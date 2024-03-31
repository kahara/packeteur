package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
	"math/rand"
)

func toss() bool {
	if rand.Float64() < 0.5 {
		return false
	}
	return true
}

// https://pkg.go.dev/github.com/google/gopacket#hdr-Creating_Packet_Data

func generateTestPacket() []byte {
	var (
		buf  = gopacket.NewSerializeBuffer()
		opts = gopacket.SerializeOptions{}
	)

	_ = gopacket.SerializeLayers(buf, opts,
		&layers.Ethernet{
			BaseLayer:    layers.BaseLayer{},
			SrcMAC:       nil,
			DstMAC:       layers.EthernetBroadcast,
			EthernetType: 0,
			Length:       0,
		},
		&layers.IPv4{
			BaseLayer:  layers.BaseLayer{},
			Version:    0,
			IHL:        0,
			TOS:        0,
			Length:     0,
			Id:         0,
			Flags:      0,
			FragOffset: 0,
			TTL:        0,
			Protocol:   17,
			Checksum:   0,
			SrcIP:      nil,
			DstIP:      nil,
			Options:    nil,
			Padding:    nil,
		},
		&layers.TCP{
			BaseLayer:  layers.BaseLayer{},
			SrcPort:    0,
			DstPort:    123,
			Seq:        0,
			Ack:        0,
			DataOffset: 0,
			FIN:        false,
			SYN:        false,
			RST:        false,
			PSH:        false,
			ACK:        false,
			URG:        false,
			ECE:        false,
			CWR:        false,
			NS:         false,
			Window:     0,
			Checksum:   0,
			Urgent:     0,
			Options:    nil,
			Padding:    nil,
		},
		gopacket.Payload([]byte{73, 86}))

	log.Debug().Any("buf", buf).Msg("")

	return buf.Bytes()
}
