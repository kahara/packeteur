package main

//import (
//	"github.com/google/gopacket"
//	"github.com/packetcap/go-pcap"
//	"github.com/rs/zerolog/log"
//)
//
//func readPackets(path string, packets chan<- gopacket.Packet) {
//	if handle, err := pcap.OpenOffline(path); err != nil {
//		log.Err(err).Msg("Could not open pcap file")
//	} else {
//		source := gopacket.NewPacketSource(handle, handle.LinkType())
//		for packet := range source.Packets() {
//			packets <- pcap.Packet{
//				B:     packet.Data(),
//				Info:  gopacket.CaptureInfo{},
//				Error: nil,
//			}
//		}
//	}
//}
