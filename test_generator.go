package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"math/rand"
)

func toss(x ...func() []gopacket.SerializableLayer) []gopacket.SerializableLayer {
	return x[rand.Intn(len(x))]()
}

// https://pkg.go.dev/github.com/google/gopacket#hdr-Creating_Packet_Data

func generateRandomPacket() gopacket.SerializeBuffer {
	var (
		buf  = gopacket.NewSerializeBuffer()
		opts = gopacket.SerializeOptions{}
	)

	_ = gopacket.SerializeLayers(buf, opts, generateRandomEthernetLayer()...)

	return buf
}

func generateRandomPayload() gopacket.Payload {
	var payload = gopacket.Payload([]byte{73, 86})

	return payload
}

func generateRandomEthernetLayer() []gopacket.SerializableLayer {
	var (
		l            []gopacket.SerializableLayer
		encapsulated = toss(generateRandomIPv4Layer, generateRandomIPv6Layer)
	)

	l = append(l, &layers.Ethernet{
		BaseLayer: layers.BaseLayer{},
		SrcMAC:    nil,
		DstMAC:    layers.EthernetBroadcast,
		EthernetType: func(x gopacket.SerializableLayer) layers.EthernetType {
			switch x.LayerType() {
			case layers.LayerTypeIPv4:
				return layers.EthernetTypeIPv4
			case layers.LayerTypeIPv6:
				return layers.EthernetTypeIPv6
			default:
				return layers.EthernetTypeLLC
			}
		}(encapsulated[0]),
		Length: 0,
	})

	l = append(l, encapsulated...)

	return l
}

func generateRandomIPv4Layer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.IPv4{
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
	})

	l = append(l, toss(generateRandomICMPv4Layer, generateRandomUDPLayer, generateRandomTCPLayer)...)

	return l
}

func generateRandomIPv6Layer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.IPv6{
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
	})

	l = append(l, toss(generateRandomICMPv6Layer, generateRandomUDPLayer, generateRandomTCPLayer)...)

	return l
}

func generateRandomICMPv4Layer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.ICMPv4{
		BaseLayer: layers.BaseLayer{},
		TypeCode:  0,
		Checksum:  0,
		Id:        0,
		Seq:       0,
	})

	return l
}

func generateRandomICMPv6Layer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.ICMPv6{
		BaseLayer: layers.BaseLayer{},
		TypeCode:  0,
		Checksum:  0,
		TypeBytes: nil,
	})

	return l
}

func generateRandomUDPLayer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.UDP{
		BaseLayer: layers.BaseLayer{},
		SrcPort:   0,
		DstPort:   0,
		Length:    0,
		Checksum:  0,
	})

	return l
}

func generateRandomTCPLayer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.TCP{
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
	})

	return l
}
