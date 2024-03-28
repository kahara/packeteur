package main

import (
	"github.com/oleiade/reflections"
	"github.com/rs/zerolog/log"
	"github.com/stoewer/go-strcase"
	"os"
)

type Config struct {
	Mode            string // One of ("capture", "collect"); in the future there could be an additional "re-relay" mode, https://github.com/kahara/packeteur/issues/2
	Device          string // Something like "eth0", "enp5s0"; for Mode "capture"
	Filter          string // For example, "udp port 53"; for Mode "capture"
	RelayEndpoint   string // Where to send packets to when Mode is "capture"
	CollectEndpoint string // Where to listen for packets when Mode is "collect"
	MetricsAddrport string // Exposed for Prometheus
}

func NewConfig() *Config {
	var config Config

	for setting, defaultValue := range map[string]string{
		"MODE":             "capture",
		"DEVICE":           "lo",
		"FILTER":           "",
		"RELAY_ENDPOINT":   "localhost:7386",
		"COLLECT_ENDPOINT": "localhost:7386",
		"METRICS_ADDRPORT": ":9108"} {
		value := os.Getenv(setting)
		if value == "" {
			value = defaultValue
		}
		reflections.SetField(&config, strcase.UpperCamelCase(setting), value)
	}

	if config.Mode == "capture" {
		config.CollectEndpoint = ""
	} else if config.Mode == "collect" {
		config.Device = ""
		config.RelayEndpoint = ""
	}

	log.Info().Any("config", config).Msg("")

	return &config
}
