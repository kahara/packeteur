# packeteur

**Capture** packets at remote spots and relay them to a single **collector** destination.
Then, at the destination, combine streams of relayed packets from multiple sources
into a single output for e.g. [tcpdump(1)'s](https://www.tcpdump.org/) consumption.

Captured packets are relayed over [ZeroMQ](https://zeromq.org/), which allows us
to keep things a bit more flexible compared to, idk, compressing the packets and
relaying them encapsulated in UDP, for example.

*I am aware of* [rpcapd(8)](https://www.tcpdump.org/manpages/rpcapd.8.html), but am
atm pretty confused about the whole thing. For starters, where's the source and how
should rpcapd be installed? It could be a straightforward exercise to get it
running on Windows, but this doesn't seem to be the case for Linux.

The following libraries kindly provide the core functionality on which Packeteur
is able to build upon:

* [rs/zerolog](https://github.com/rs/zerolog)
* [packetcap/go-pcap](https://github.com/packetcap/go-pcap)
* [google/gopacket](https://github.com/google/gopacket)
* [seancfoley/ipaddress-go](github.com/seancfoley/ipaddress-go)
* [prometheus/client_golang](https://github.com/prometheus/client_golang)
* [pebbe/zmq4](https://github.com/pebbe/zmq4)

## Configuration

All settings go through environment variables:

**Setting**|      **Default**       |**Notes**
:-----:|:----------------------:|:-----:
MODE|       `capture`        |Either `capture` or `collect`
DEVICE|          `lo`          |Something like `eth0` or `enp5s0`; for MODE `capture`
FILTER|                        |For example `udp port 53`; for MODE `capture`
RELAY\_ENDPOINT| `tcp://localhost:7386` |Where to send packets to when MODE is `capture`
COLLECT\_ENDPOINT| `tcp://localhost:7386` |Where to listen for packets when MODE is `collect`
METRICS\_ADDRPORT|        `:9108`         |Exposed for Prometheus; see below

## Running

To capture:

```console
MODE=capture\
    DEVICE=enp5s0 \
    FILTER="udp port 53" \
    RELAY_ENDPOINT=tcp://localhost:7386 \
    pcktr
```

To collect:

```console
MODE=collect \
    COLLECT_ENDPOINT=tcp://localhost:7386 \
    pcktr | tcpdump -r - -tttt -vvvv
```

## Metrics

In addition to the [builtins](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#hdr-Metrics),
the following Packeteur-specific ones are exposed:

```
packeteur_capture_bytes_bucket{address_family="IPv4",le="..."}  # Buckets are
packeteur_capture_bytes_bucket{address_family="IPv6",le="..."}  # 2^6..2^16
packeteur_capture_bytes_count{address_family="IPv4"}
packeteur_capture_bytes_count{address_family="IPv6"}
packeteur_capture_bytes_sum{address_family="IPv4"}
packeteur_capture_bytes_sum{address_family="IPv6"}
packeteur_capture_total{address_family="IPv4"}
packeteur_capture_total{address_family="IPv6"}
```

## End-to-end testing

