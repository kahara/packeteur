package main

import (
	"github.com/rs/zerolog/log"
	"net"
)

func collect(endpoint string) {
	var (
		err    error
		addr   *net.UDPAddr
		conn   *net.UDPConn
		buf    []byte
		count  int
		source *net.UDPAddr
	)

	log.Info().Str("endpoint", endpoint).Msg("Packeteur is collecting")

	if addr, err = net.ResolveUDPAddr("udp", endpoint); err != nil {
		log.Err(err)
	}
	if conn, err = net.ListenUDP("udp", addr); err != nil {
		log.Err(err)
	}

	log.Info().Str("addr", addr.String()).Msg("Listening")

	for {
		if count, _, _, source, err = conn.ReadMsgUDP(buf, nil); err != nil {
			log.Err(err)
		}
		log.Debug().Int("count", count).Any("source", source).Bytes("packet", buf).Msg("Packet received")
	}
}
