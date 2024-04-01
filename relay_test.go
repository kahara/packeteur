package main

import (
	"github.com/cespare/xxhash/v2"
	"github.com/google/gopacket"
	"github.com/packetcap/go-pcap"
	"github.com/rs/zerolog/log"
	"testing"
)

const (
	RelayTestEndpoint    = "localhost:17386"
	RelayTestPacketCount = 10
)

func TestRelay(t *testing.T) {
	var (
		outgoing = make(chan pcap.Packet)
		seen     = make(map[uint64]bool)
		incoming = make(chan []byte, 1024)
	)

	go relay(outgoing, RelayTestEndpoint)
	go collect(RelayTestEndpoint, incoming)

	for i := 0; i < RelayTestPacketCount; i++ {
		outgoingPacket := generateTestPacket()
		seen[xxhash.Sum64(outgoingPacket.Bytes())] = true
		log.Debug().Bytes("packet", outgoingPacket.Bytes()).Msg("Sending packet to relayer")
		outgoing <- pcap.Packet{
			B:     outgoingPacket.Bytes(),
			Info:  gopacket.CaptureInfo{},
			Error: nil,
		}
	}

	for incomingPacket := range incoming {
		log.Debug().Bytes("packet", incomingPacket).Msg("Received packet from collector")
	}
}
