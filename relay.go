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
		count        int
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
		if netLayer := p.NetworkLayer(); netLayer != nil {
			_, dst := netLayer.NetworkFlow().Endpoints()

			// Skip own traffic
			if endpointAddr.IP.String() == dst.String() { // Would like to compare []bytes
				continue
			}

			// Send the packet
			if count, _, err = conn.WriteMsgUDP(packet.B, nil, nil); err != nil {
				log.Err(err).Msg("")
			} else {
				log.Debug().Str("destination", dst.String()).Int("length", count).Msg("Sent to collector")
			}

			// Record our beloved metrics
			addressFamily := "undefined"
			if dstAddr, err = ipaddr.NewIPAddressFromBytes(dst.Raw()); err == nil {
				if dstAddr.IsIPv6() {
					addressFamily = "IPv6"
				} else if dstAddr.IsIPv4() {
					addressFamily = "IPv4"
				}
			}
			captured_total_metric.WithLabelValues(addressFamily).Inc()
			captured_bytes_metric.WithLabelValues(addressFamily).Observe(float64(len(packet.B)))
		}
	}
}
