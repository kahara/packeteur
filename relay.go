package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/packetcap/go-pcap"
	"github.com/rs/zerolog/log"
	"github.com/seancfoley/ipaddress-go/ipaddr"
	"net"
)

func relay(packets <-chan pcap.Packet, endpoint string) {
	var (
		err          error
		endpointAddr *net.UDPAddr
		conn         *net.UDPConn
		dstAddr      *ipaddr.IPAddress
	)

	log.Info().Str("endpoint", endpoint).Msg("Packeteur is relaying")

	if endpointAddr, err = net.ResolveUDPAddr("udp", endpoint); err != nil {
		log.Err(err).Msg("")
	}
	if conn, err = net.DialUDP("udp", nil, endpointAddr); err != nil {
		log.Err(err).Msg("")
	}

	for packet := range packets {
		p := gopacket.NewPacket(packet.B, layers.LayerTypeEthernet, gopacket.Default)
		log.Debug().Str("packet", p.String()).Msg("Received packet")

		// Skip own traffic
		addressFamily := "undefined"
		if netLayer := p.NetworkLayer(); netLayer != nil {
			_, dst := netLayer.NetworkFlow().Endpoints()

			// Take note of address family while on it
			if dstAddr, err = ipaddr.NewIPAddressFromBytes(dst.Raw()); err == nil {
				if dstAddr.IsIPv6() {
					addressFamily = "IPv6"
				} else if dstAddr.IsIPv4() {
					addressFamily = "IPv4"
				}
			}

			// This should compare []bytes and also ports
			if endpointAddr.IP.String() == dst.String() {
				continue
			}
		}

		// Send the packet
		if _, _, err = conn.WriteMsgUDP(packet.B, nil, nil); err != nil {
			log.Err(err).Msg("Something went wrong while sending packet to collector")
			continue
		}

		// Record our beloved metrics
		relayed_total_metric.WithLabelValues(addressFamily).Inc()
		relayed_bytes_metric.WithLabelValues(addressFamily).Observe(float64(len(packet.B)))
	}
}
