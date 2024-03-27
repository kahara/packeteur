package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/packetcap/go-pcap"
	"github.com/seancfoley/ipaddress-go/ipaddr"
)

func relay(packets <-chan pcap.Packet, endpoint string) {
	var addressFamily = "undefined"

	for packet := range packets {
		p := gopacket.NewPacket(packet.B, layers.LayerTypeEthernet, gopacket.Default)
		if net := p.NetworkLayer(); net != nil {
			_, dst := net.NetworkFlow().Endpoints()
			if addr, err := ipaddr.NewIPAddressFromBytes(dst.Raw()); err == nil {
				if addr.IsIPv6() {
					addressFamily = "IPv6"
				} else if addr.IsIPv4() {
					addressFamily = "IPv4"
				}
			}
			// Record our beloved metrics
			captured_total_metric.WithLabelValues(addressFamily).Inc()
			captured_bytes_metric.WithLabelValues(addressFamily).Observe(float64(len(packet.B)))

			//
		}
	}
}
