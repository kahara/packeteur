# packeteur

Capture packets at remote spots and relay them to a single collector destination.
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

* [packetcap/go-pcap](https://github.com/packetcap/go-pcap)
* [pebbe/zmq4](https://github.com/pebbe/zmq4)

## Usage

```console
pcktr relay -device eth0 -filter 'udp port 53' -destination 10.0.1.2:7386
```

```console
pcktr collect -address :7386  # running on 10.0.1.2; write  to stdout
```
