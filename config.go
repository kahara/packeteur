package main

import "flag"

type Config struct {
	Relay                    bool // These two flags represent the mode of operation;
	Collect                  bool // both would be set in https://github.com/kahara/packeteur/issues/2
	CaptureDevice            string
	CaptureFilter            string
	RelayDestinationAddrport string
	CollectAddrport          string
	MetricsAddrport          string
}

func NewConfig() *Config {
	var config Config

	flag.StringVar(&config.CaptureDevice, "device", "lo", "capture device")
	flag.StringVar(&config.CaptureFilter, "filter", "", "capture filter")
	flag.StringVar(&config.RelayDestinationAddrport, "destination", "", "relay destination addr:port")
	flag.StringVar(&config.CollectAddrport, "address", "", "collect addr:port")
	flag.StringVar(&config.MetricsAddrport, "metrics", ":9108", "metrics addr:port")
	flag.Parse()

	if config.CaptureDevice != "" {
		config.Relay = true
	}

	return &config
}
