package main

import (
	"bytes"
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
	)

	log.Info().Str("endpoint", endpoint).Msg("Packeteur is relaying")

	if addr, err = net.ResolveUDPAddr("udp", endpoint); err != nil {
		log.Err(err)
	}

	if conn, err = net.DialUDP("udp", nil, addr); err != nil {
		log.Err(err)
	}
	//defer conn.Close()

	for packet := range packets {
		p := gopacket.NewPacket(packet.B, layers.LayerTypeEthernet, gopacket.Default)
		if netLayer := p.NetworkLayer(); netLayer != nil {
			log.Debug().Any("net", netLayer).Msg("Extract net layer")
			_, dst := netLayer.NetworkFlow().Endpoints()

			if a, err := ipaddr.NewIPAddressFromBytes(dst.Raw()); err == nil {
				if bytes.Compare(a.Bytes(), dst.Raw()) == 0 {
					continue
				}
				if a.IsIPv6() {
					addressFamily = "IPv6"
				} else if a.IsIPv4() {
					addressFamily = "IPv4"
				}
			}

			// Send the packet
			if _, _, err = conn.WriteMsgUDP(packet.B, nil, nil); err != nil {
				log.Err(err)
			}
			//if _, err = conn.Write(packet.B); err != nil {
			//	log.Err(err)
			//}

			// Record our beloved metrics
			captured_total_metric.WithLabelValues(addressFamily).Inc()
			captured_bytes_metric.WithLabelValues(addressFamily).Observe(float64(len(packet.B)))
		}
	}
}
