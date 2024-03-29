# packeteur

**Capture** packets at remote spots, compress and encapsulate them in UDP, and relay
to a single **collector** destination. Then, at the destination, combine streams of
relayed packets from multiple sources into a single output for e.g.
[tcpdump(1)'s](https://www.tcpdump.org/) consumption. No half-baked attempt at keeping
the transport confidential and tamper-proof this time &ndash; use
[Wireguard](https://www.wireguard.com/) for that.

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

## Configuration

All settings go through environment variables:

**Setting**|     **Default**     |**Notes**
:-----:|:-------------------:|:-----:
MODE|      `capture`      |Either `capture` or `collect`
DEVICE|        `lo`         |Something like `eth0` or `enp5s0`; for MODE `capture`
FILTER| `not udp port 7386` |For MODE `capture`; [BPF](https://en.wikipedia.org/wiki/Berkeley_Packet_Filter) syntax
RELAY\_ENDPOINT|  `localhost:7386`   |Where to send packets to when MODE is `capture`
COLLECT\_ENDPOINT|  `:7386`   |Where to listen for packets when MODE is `collect`
METRICS\_ADDRPORT|       `:9108`       |Exposed for Prometheus; see ["Metrics"](#metrics) below

## Running

To capture:

```console
docker run --rm -it \
    -e MODE=capture\
    -e DEVICE=enp5s0 \
    -e FILTER="not udp port 7386" \
    -e RELAY_ENDPOINT=localhost:7386 \
    ghcr.io/kahara/packeteur:latest
```

To collect:

```console
docker run --rm -it \
    -e MODE=collect \
    -e COLLECT_ENDPOINT=:7386 \
    ghcr.io/kahara/packeteur:latest
    #| tcpdump -r - -tttt -vvvv
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

