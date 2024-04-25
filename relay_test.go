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
	RelayTestFilename    = "/tmp/packeteur-test-relay.pcap"
)

func TestRelay(t *testing.T) {
	var (
		outgoing = make(chan pcap.Packet, 1000)
		seen     = make(map[uint64]bool)
		incoming = make(chan []byte, 1000)
	)

	go relay(outgoing, RelayTestEndpoint)
	go collect(RelayTestEndpoint, incoming)

	for i := 0; i < RelayTestPacketCount; i++ {
		outgoingPacket := generateRandomPacket()
		seen[xxhash.Sum64(outgoingPacket.Bytes())] = true
		log.Debug().Bytes("packet", outgoingPacket.Bytes()).Msg("Sending packet to relayer")
		outgoing <- pcap.Packet{
			B:     outgoingPacket.Bytes(),
			Info:  gopacket.CaptureInfo{},
			Error: nil,
		}
	}

	count := RelayTestPacketCount
	for incomingPacket := range incoming {
		log.Debug().Bytes("packet", incomingPacket).Msg("Received packet from collector")

		// Verify collector isn't sending anything extraneous
		count--
		if count < 0 {
			t.Errorf("Too many packets seen, expected %d count", RelayTestPacketCount)
		}

		// Verify all the packets that went out made the round-trip intact
		h := xxhash.Sum64(incomingPacket)
		if seen[h] {
			delete(seen, h)
		}
		if len(seen) == 0 {
			break
		}
	}
}
