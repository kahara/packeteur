package main

import (
	"github.com/packetcap/go-pcap"
	"github.com/rs/zerolog/log"
)

func capture(captureDevice string, captureFilter string) chan pcap.Packet {

	var (
		err    error
		handle *pcap.Handle
	)

	log.Info().Str("device", captureDevice).Str("filter", captureFilter).Msg("Packeteur is capturing")

	if handle, err = pcap.OpenLive(captureDevice, 1600, true, 0, false); err != nil {
		log.Err(err).Msg("")
	}
	if err = handle.SetBPFFilter(captureFilter); err != nil {
		log.Err(err).Msg("")
	}

	return handle.Listen()
}
