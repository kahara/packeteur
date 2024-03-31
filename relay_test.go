package main

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func TestRelay(t *testing.T) {
	x := generateTestPacket()

	log.Debug().Any("x", x).Msg("")
}
