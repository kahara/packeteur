package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strings"
)

func toss(x ...func() []gopacket.SerializableLayer) []gopacket.SerializableLayer {
	return x[rand.Intn(len(x))]()
}

// https://pkg.go.dev/github.com/google/gopacket#hdr-Creating_Packet_Data

func generateRandomPacket() gopacket.SerializeBuffer {
	var (
		buf  = gopacket.NewSerializeBuffer()
		opts = gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		l = generateRandomEthernetLayer()
	)

	log.Debug().Any("layers", l).Msg("Layers generated")

	if err := gopacket.SerializeLayers(buf, opts, l...); err != nil {
		log.Err(err).Msg("SerializeLayers")
	}

	return buf
}

func generateRandomPayload() gopacket.Payload {
	var payload = make([]byte, 23+rand.Intn(450))
	crand.Read(payload)
	return gopacket.Payload(payload)
}

func generateRandomHardwareAddress() []byte {
	var address = make([]byte, 6)
	crand.Read(address)
	return address
}

func generateRandomIPAddress(kind string) []byte {
	var address []byte

	switch strings.ToLower(kind) {
	case "ipv4":
		address = make([]byte, 4)
	case "ipv6":
		address = make([]byte, 16)
	default:
		panic(fmt.Sprintf("I don't know how to come up with an IP address of kind %s", kind))
	}

	crand.Read(address)
	return address
}

func generateRandomPort() uint16 {
	return uint16(rand.Intn(65536))
}

func generateRandomFlowLabel() uint32 {
	var flowLabel = make([]byte, 4)
	crand.Read(flowLabel)
	return 0xFFFFF & binary.BigEndian.Uint32(flowLabel)
}

func generateRandomEthernetLayer() []gopacket.SerializableLayer {
	var (
		l            []gopacket.SerializableLayer
		encapsulated = toss(generateRandomIPv4Layer, generateRandomIPv6Layer)
	)

	l = append(l, &layers.Ethernet{
		BaseLayer: layers.BaseLayer{},
		SrcMAC:    generateRandomHardwareAddress(),
		DstMAC:    generateRandomHardwareAddress(),
		EthernetType: func(x gopacket.SerializableLayer) layers.EthernetType {
			switch x.LayerType() {
			case layers.LayerTypeIPv4:
				return layers.EthernetTypeIPv4
			case layers.LayerTypeIPv6:
				return layers.EthernetTypeIPv6
			default:
				log.Debug().Any("type", x.LayerType().String()).Msg("unknown layer type")
				return layers.EthernetTypeLLC
			}
		}(encapsulated[0]),
	})

	l = append(l, encapsulated...)

	// Marry IP payloads
	for i, x := range l {
		//log.Debug().Int("index", i).Any("type", x.LayerType().String()).Any("layer", x).Msg("dumping layers")
		switch x.LayerType() {
		case layers.LayerTypeICMPv6:
			x.(*layers.ICMPv6).SetNetworkLayerForChecksum(l[i-1].(gopacket.NetworkLayer))
		case layers.LayerTypeUDP:
			x.(*layers.UDP).SetNetworkLayerForChecksum(l[i-1].(gopacket.NetworkLayer))
		case layers.LayerTypeTCP:
			x.(*layers.TCP).SetNetworkLayerForChecksum(l[i-1].(gopacket.NetworkLayer))
		default:
			continue
		}
	}
	return l
}

func generateRandomIPv4Layer() []gopacket.SerializableLayer {
	var (
		l            []gopacket.SerializableLayer
		encapsulated = toss(generateRandomICMPv4Layer, generateRandomUDPLayer, generateRandomTCPLayer)
	)

	l = append(l, &layers.IPv4{
		BaseLayer:  layers.BaseLayer{},
		Version:    4,
		IHL:        0,
		TOS:        0,
		Id:         0,
		Flags:      0,
		FragOffset: 0,
		TTL:        uint8(rand.Intn(64)),
		Protocol: func(x gopacket.SerializableLayer) layers.IPProtocol {
			switch x.LayerType() {
			case layers.LayerTypeICMPv4:
				return layers.IPProtocolICMPv4
			case layers.LayerTypeUDP:
				return layers.IPProtocolUDP
			case layers.LayerTypeTCP:
				return layers.IPProtocolTCP
			default:
				return layers.IPProtocol(0)
			}
		}(encapsulated[0]),
		Checksum: 0,
		SrcIP:    generateRandomIPAddress("ipv4"),
		DstIP:    generateRandomIPAddress("ipv4"),
		Options:  nil,
		Padding:  nil,
	})

	l = append(l, encapsulated...)

	return l
}

func generateRandomIPv6Layer() []gopacket.SerializableLayer {
	var (
		l            []gopacket.SerializableLayer
		encapsulated = toss(generateRandomICMPv6Layer, generateRandomUDPLayer, generateRandomTCPLayer)
	)

	l = append(l, &layers.IPv6{
		BaseLayer:    layers.BaseLayer{},
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    generateRandomFlowLabel(),
		NextHeader:   layers.IPProtocolNoNextHeader,
		HopLimit:     uint8(rand.Intn(64)),
		SrcIP:        generateRandomIPAddress("ipv6"),
		DstIP:        generateRandomIPAddress("ipv6"),
		//HopByHop: nil,
	})

	l = append(l, encapsulated...)

	return l
}

func generateRandomICMPv4Layer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.ICMPv4{
		BaseLayer: layers.BaseLayer{},
		TypeCode:  layers.ICMPv4TypeCode(uint16(rand.Intn(65536))),
		Id:        uint16(rand.Intn(65536)),
		Seq:       uint16(rand.Intn(65536)),
	})

	return l
}

func generateRandomICMPv6Layer() []gopacket.SerializableLayer {
	var l []gopacket.SerializableLayer

	l = append(l, &layers.ICMPv6{
		BaseLayer: layers.BaseLayer{},
		TypeCode:  layers.ICMPv6TypeCode(uint16(rand.Intn(65536))),
	})

	return l
}

func generateRandomUDPLayer() []gopacket.SerializableLayer {
	var (
		l            []gopacket.SerializableLayer
		encapsulated = generateRandomPayload()
	)

	l = append(l, &layers.UDP{
		BaseLayer: layers.BaseLayer{},
		SrcPort:   layers.UDPPort(generateRandomPort()),
		DstPort:   layers.UDPPort(generateRandomPort()),
		Checksum:  0,
	})

	l = append(l, encapsulated)

	return l
}

func generateRandomTCPLayer() []gopacket.SerializableLayer {
	var (
		l            []gopacket.SerializableLayer
		encapsulated = generateRandomPayload()
	)

	l = append(l, &layers.TCP{
		BaseLayer:  layers.BaseLayer{},
		SrcPort:    layers.TCPPort(generateRandomPort()),
		DstPort:    layers.TCPPort(generateRandomPort()),
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

	l = append(l, encapsulated)

	return l
}
