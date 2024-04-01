package main

import (
	"github.com/rs/zerolog"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var config = NewConfig()

	setupMetrics(config.Mode)
	go metrics(config.MetricsAddrport)

	if config.Mode == "capture" {
		relay(capture(config.Device, config.Filter), config.RelayEndpoint)
	} else if config.Mode == "collect" {
		collect(config.CollectEndpoint, nil)
	} else {
		panic("Not sure what went wrong, but we're done here.")
	}
}
