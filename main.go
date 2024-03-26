package main

import (
	"github.com/rs/zerolog"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	captureDevice := "enp5s0"
	captureFilter := "udp port 53"

	capture(captureDevice, captureFilter)
}
