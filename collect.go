package main

import (
	"github.com/ef-ds/deque/v2"
	"github.com/rs/zerolog/log"
	"net"
	"sync"
	"time"
)

func collect(endpoint string, packets chan []byte) {
	var (
		err   error
		conn  net.PacketConn //*net.UDPConn
		addr  net.Addr
		count int
		stage deque.Deque[[]byte]
		m     sync.Mutex
	)

	log.Info().Str("endpoint", endpoint).Msg("Packeteur is collecting")

	if packets == nil {
		packets = make(chan []byte, 1024)
		go func() {
			log.Debug().Msg("Processor is processing")
			for {
				buf := make([]byte, 65536)
				select {
				case buf = <-packets:
					log.Debug().Int("length", len(buf)).Msg("Processing packet")
				}
			}
		}()
	}

	// Deliver absorbed packets
	go func() {
		for {
			time.Sleep(time.Microsecond)
			m.Lock()
			if e, _ := stage.PopFront(); e != nil {
				select {
				case packets <- e:
				default:
					stage.PushFront(e)
					break
				}
			}
			m.Unlock()
		}
	}()

	if conn, err = net.ListenPacket("udp", endpoint); err != nil {
		log.Err(err).Any("addr", addr).Msg("Something went wrong while attempting to listen")
	}

	for {
		buf := make([]byte, 65536) // Consider if possible to recycle the buffer
		if count, addr, err = conn.ReadFrom(buf); err != nil {
			log.Err(err).Msg("")
			continue
		}
		log.Debug().Int("length", count).Str("source", addr.String()).Msg("Relayed packet received")

		// Absorb collected packets
		m.Lock()
		stage.PushBack(buf[:count])
		m.Unlock()
	}
}
