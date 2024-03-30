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
		err           error
		addr          *net.UDPAddr
		conn          *net.UDPConn
		addressFamily = "undefined"
		count         int
	)

	log.Info().Str("endpoint", endpoint).Msg("Packeteur is relaying")

	if addr, err = net.ResolveUDPAddr("udp", endpoint); err != nil {
		log.Err(err).Msg("")
	}

	if conn, err = net.DialUDP("udp", nil, addr); err != nil {
		log.Err(err).Msg("")
	}

	for packet := range packets {
		p := gopacket.NewPacket(packet.B, layers.LayerTypeEthernet, gopacket.Default)
		if netLayer := p.NetworkLayer(); netLayer != nil {
			_, dst := netLayer.NetworkFlow().Endpoints()
			if a, err := ipaddr.NewIPAddressFromBytes(dst.Raw()); err == nil {
				// Skip own traffic
				if addr.IP.String() == dst.String() { // Would like to compare []bytes
					continue
				}
				if a.IsIPv6() {
					addressFamily = "IPv6"
				} else if a.IsIPv4() {
					addressFamily = "IPv4"
				}
			}

			// Send the packet
			if count, _, err = conn.WriteMsgUDP(packet.B, nil, nil); err != nil {
				log.Err(err).Msg("")
			} else {
				log.Debug().Str("destination", dst.String()).Int("length", count).Msg("Sent to collector")
			}

			// Record our beloved metrics
			captured_total_metric.WithLabelValues(addressFamily).Inc()
			captured_bytes_metric.WithLabelValues(addressFamily).Observe(float64(len(packet.B)))
		}
	}
}
