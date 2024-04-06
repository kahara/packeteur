package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"math/rand"
)

func toss() bool {
	if rand.Float64() < 0.5 {
		return false
	}
	return true
}

// https://pkg.go.dev/github.com/google/gopacket#hdr-Creating_Packet_Data

func generateRandomPacket() gopacket.SerializeBuffer {
	var (
		buf  = gopacket.NewSerializeBuffer()
		opts = gopacket.SerializeOptions{}
	)

	_ = gopacket.SerializeLayers(buf, opts,
		generateRandomEthernet(),
		generateRandomIPv4(),
		generateRandomTCP(),
		generateRandomPayload())

	return buf
}

func generateRandomPayload() gopacket.Payload {
	var payload = gopacket.Payload([]byte{73, 86})

	return payload
}

func generateRandomEthernet() *layers.Ethernet {
	var ethernet = layers.Ethernet{
		BaseLayer:    layers.BaseLayer{},
		SrcMAC:       nil,
		DstMAC:       layers.EthernetBroadcast,
		EthernetType: 0,
		Length:       0,
	}

	return &ethernet
}

func generateRandomIPv4() *layers.IPv4 {
	var ipv4 = layers.IPv4{
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
	}

	return &ipv4
}

func generateRandomIPv6() *layers.IPv6 {
	var ipv6 = layers.IPv6{
		BaseLayer:    layers.BaseLayer{},
		Version:      0,
		TrafficClass: 0,
		FlowLabel:    0,
		Length:       0,
		NextHeader:   0,
		HopLimit:     0,
		SrcIP:        nil,
		DstIP:        nil,
		HopByHop:     nil,
	}

	return &ipv6
}

func generateRandomICMPv4() *layers.ICMPv4 {
	var icmpv4 = layers.ICMPv4{
		BaseLayer: layers.BaseLayer{},
		TypeCode:  0,
		Checksum:  0,
		Id:        0,
		Seq:       0,
	}

	return &icmpv4
}

func generateRandomICMPv6() *layers.ICMPv6 {
	var icmpv6 = layers.ICMPv6{
		BaseLayer: layers.BaseLayer{},
		TypeCode:  0,
		Checksum:  0,
		TypeBytes: nil,
	}

	return &icmpv6
}

func generateRandomUDP() *layers.UDP {
	var udp = layers.UDP{
		BaseLayer: layers.BaseLayer{},
		SrcPort:   0,
		DstPort:   0,
		Length:    0,
		Checksum:  0,
	}

	return &udp
}

func generateRandomTCP() *layers.TCP {
	var tcp = layers.TCP{
		BaseLayer:  layers.BaseLayer{},
		SrcPort:    0,
		DstPort:    0,
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
	}

	return &tcp
}
