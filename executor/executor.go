package executor

import (
	"github.com/homedepot/github-webhook/configuration"
	"github.com/homedepot/github-webhook/events"
	"github.com/homedepot/github-webhook/matcher"
	"expvar"
	"os/exec"
	"log"
	"os"
)

var config configuration.Configuration
var debug bool
var dryRun bool
var metrics *expvar.Map

func ProcessEvent(event events.Event) error {

	var errors []error
	for _, trigger := range config.Triggers {
		if matcher.Matches(trigger, event) {
			errors = append(errors, Execute(trigger, event))
		}
	}

	return nil
}

func Execute(trigger configuration.Trigger, event events.Event) error {

	log.Println("Executing",trigger.Name)

	if dryRun {
		log.Println("DryRun enabled. Execution skipped.")
		return nil
	}

	cmd := exec.Command(
		trigger.Run.Path,
		trigger.Run.Args...
	)

	if debug {
		log.Println("Redirecting stdout/stderr for",trigger.Name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	log.Println(trigger.Name, "was successful!")

	return nil
}

func SetConfig(c configuration.Configuration) {
	config = c
}

func SetDebug(d bool) {
	debug = d
}

func SetDryRun(d bool) {
	dryRun = d
}

func SetMetrics(m *expvar.Map) {
	metrics = m
}
