package main

import (
	"github.com/rs/zerolog/log"
	"net"
)

func collect(endpoint string) {
	var (
		err     error
		conn    net.PacketConn //*net.UDPConn
		addr    net.Addr
		count   int
		packets = make(chan []byte, 1024)
	)

	log.Info().Str("endpoint", endpoint).Msg("Packeteur is collecting")

	go func() {
		log.Debug().Msg("Processor is processing")
		for {
			buf := make([]byte, 65536)
			select {
			case buf = <-packets:
				log.Debug().Bytes("buf", buf).Msg("Processing packet")
			}
		}
	}()

	if conn, err = net.ListenPacket("udp", endpoint); err != nil {
		log.Err(err)
	}

	for {
		buf := make([]byte, 65536) // Consider if possible to recycle the buffer
		if count, addr, err = conn.ReadFrom(buf); err != nil {
			log.Err(err)
			continue
		}
		log.Debug().Int("count", count).Str("source", addr.String()).Bytes("packet", buf).Msg("Packet received")
		packets <- buf[:count]
	}
}
