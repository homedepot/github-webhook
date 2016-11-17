package main

import (
	"expvar"
	"github.com/homedepot/github-webhook/configuration"
	"github.com/homedepot/github-webhook/executor"
	"github.com/homedepot/github-webhook/matcher"
	"github.com/homedepot/github-webhook/postreceive"
	"github.com/homedepot/github-webhook/status"
	"github.com/homedepot/github-webhook/web"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	app          = kingpin.New("github-webhook", "Configurable GitHub Webhook").Author("THD Engineering & Operations")
	debug        = app.Flag("debug", "Enables debug mode").Default("false").OverrideDefaultFromEnvar("DEBUG").Bool()
	dryRun       = app.Flag("dry-run", "Disables actual execution of webhooks").Default("false").OverrideDefaultFromEnvar("DRYRUN").Bool()
	validateOnly = app.Flag("validate-only", "Validates configuration, then exits").Default("false").OverrideDefaultFromEnvar("VALIDATE_ONLY").Bool()
	port         = app.Flag("port", "Port to listen on").Default("8080").OverrideDefaultFromEnvar("PORT").Int()
	configFile   = app.Flag("config-file", "Configuration File").OverrideDefaultFromEnvar("CONFIG_FILE").Default("").String()

	logStdout = app.Flag("log-stdout", "Log to stdout").Default("true").Bool()
	logPrefix = app.Flag("log-prefix", "Log prefix").OverrideDefaultFromEnvar("LOG_PREFIX").Default("").String()
	logFormat = app.Flag("log-format", "Log format").OverrideDefaultFromEnvar("LOG_FORMAT").Default(strconv.Itoa(log.LstdFlags)).Int()

	githubSecret = app.Flag("github-secret", "Optional Github Secret").OverrideDefaultFromEnvar("GITHUB_SECRET").Default("").String()

	statusInterval = app.Flag("status-interval", "Status log message interval").Default("300s").OverrideDefaultFromEnvar("STATUS_INTERVAL").Duration()
)

var (
	Version     = "0.0.1"
	GitCommit   = "HEAD"
	BuildStamp  = "UNKNOWN"
	FullVersion = Version + "+" + GitCommit + "-" + BuildStamp
	Splash      = []byte("github-webhook " + FullVersion)
)

var metrics = expvar.NewMap("metrics")

func init() {
	app.Version(FullVersion)
}

func main() {

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *logStdout {
		log.SetOutput(os.Stdout)
	}

	if len(*logPrefix) > 0 {
		log.SetPrefix(*logPrefix + " ")
	}

	log.SetFlags(*logFormat)

	config, err := configuration.LoadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if *validateOnly {
		log.Println("Configuration seems to OK.")
		os.Exit(0)
	}

	postreceive.SetSecret(*githubSecret)
	postreceive.SetMetrics(metrics)

	matcher.SetDebug(*debug)
	executor.SetDebug(*debug)
	executor.SetDryRun(*dryRun)
	executor.SetConfig(config)
	executor.SetMetrics(metrics)

	http.HandleFunc("/postreceive", postreceive.HttpHandler)
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(Splash)
	})

	web.Listen(*port)

	status.SetMetrics(metrics)
	statusDone := status.Start(*statusInterval)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	statusDone <- true

}
