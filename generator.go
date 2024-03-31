package main

import (
	"github.com/google/gopacket"
	"github.com/rs/zerolog/log"
	"math/rand"
)

func toss() bool {
	if rand.Float64() < 0.5 {
		return false
	}
	return true
}

func generateTestPacket() []byte {
	var packet []byte

	x := gopacket.NewPacket(nil, gopacket.DecodeUnknown, gopacket.Default)

	log.Debug().Any("x", x).Msg("")

	return packet
}
