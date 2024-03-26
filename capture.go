package main

import (
	pcap "github.com/packetcap/go-pcap"
	"github.com/rs/zerolog/log"
)

func capture(captureDevice string, captureFilter string) {

	var (
		err    error
		handle *pcap.Handle
	)

	if handle, err = pcap.OpenLive(captureDevice, 1600, true, 0, false); err != nil {
		log.Err(err)
	}
	if err = handle.SetBPFFilter(captureFilter); err != nil {
		log.Err(err)
	}
	for packet := range handle.Listen() {
		log.Debug().Any("info", packet.Info).Bytes("packet", packet.B).Msg("")
	}
}
