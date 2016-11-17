package status

import (
	"fmt"
	"time"
	"expvar"
	"log"
)

const status_format = "STATUS" +
" requests=\"%v\"" +
" pings=\"%v\"" +
" pushes=\"%v\"" +
" pull_requests=\"%v\"" +
" errors=%v"

func Start(statusInterval time.Duration) chan bool {

	done := make(chan bool, 1)

	//start status routine
	go func(done chan bool, statusInterval time.Duration) {
		statusTicker := time.NewTicker(statusInterval)
		defer statusTicker.Stop()

		StatusLoop:
		for {
			select {
			case <-statusTicker.C:
				log.Println(fmt.Sprintf(
					status_format,
					zeroOrValue(metrics.Get("request")),
					zeroOrValue(metrics.Get("ping")),
					zeroOrValue(metrics.Get("push")),
					zeroOrValue(metrics.Get("pull_request")),
					zeroOrValue(metrics.Get("error")),
				))
			case <-done:
				log.Println("Exiting status thread due to signal")
				break StatusLoop
			}
		}
		close(done)
	}(done, statusInterval)

	return done
}

var metrics *expvar.Map

func SetMetrics(m *expvar.Map) {
	metrics = m
}

func zeroOrValue(value interface{}) interface{} {
	if value == nil {
		return 0
	}
	return value
}


