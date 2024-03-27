package main

import (
	"github.com/rs/zerolog"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	//log.Info().Msg("This is Packeteur")

	var config = NewConfig()

	setupMetrics(config.Relay, config.Collect)
	go metrics(config.MetricsAddrport)

	relay(capture(config.CaptureDevice, config.CaptureFilter))
}
