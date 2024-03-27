package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	Namespace = "packeteur"
)

var (
	captured_total_metric  *prometheus.CounterVec
	captured_bytes_metric  *prometheus.HistogramVec
	relayed_total_metric   *prometheus.CounterVec
	collected_total_metric *prometheus.CounterVec
)

func setupMetrics(mode string) {
	switch mode {
	case "capture":
		captured_total_metric = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "capture",
			Name:      "total",
		}, []string{"address_family"})
		captured_bytes_metric = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: Namespace,
			Subsystem: "capture",
			Name:      "bytes",
			Buckets:   []float64{64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536},
		}, []string{"address_family"})
		relayed_total_metric = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "relay",
			Name:      "total",
		}, []string{})
	case "collect":
		collected_total_metric = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Subsystem: "collect",
			Name:      "total",
		}, []string{})
	default:
		panic("Not sure what went wrong, but we're done here.")
	}
}

func metrics(addrPort string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addrPort, nil); err != nil {
		log.Fatal().Err(err).Str("addrport", addrPort).Msg("Could not expose Prometheus metrics")
	}
}
