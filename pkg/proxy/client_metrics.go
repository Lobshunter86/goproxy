package proxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type LocalServerMetricsVec struct {
	RequestCount             *prometheus.CounterVec
	HandshakeErrCount        *prometheus.CounterVec
	CurrentConnGauge         *prometheus.GaugeVec
	RequestHandshakeDuration *prometheus.HistogramVec
	RequestHandlingDuration  *prometheus.HistogramVec // not includes handshake time
}

var metricVecs = LocalServerMetricsVec{
	RequestCount: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "counter of requests",
		},
		[]string{"protocol"},
	),

	HandshakeErrCount: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hanshake_err_count",
			Help: "counter of errors in quic handshake, not includes errors after that",
		},
		[]string{"protocol"},
	),

	CurrentConnGauge: promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "current_connections",
			Help: "current connection count",
		},
		[]string{"protocol"},
	),

	RequestHandshakeDuration: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_handshake_duration",
			Help:    "bucketed histogram of quic handshaking time (ms) of requests",
			Buckets: prometheus.LinearBuckets(200, 100, 18),
		},
		[]string{"protocol"},
	),

	RequestHandlingDuration: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_handling_duration",
			Help:    "bucketed histogram of quic processing time (s) of requests, not includes quic handshake time",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10),
		},
		[]string{"protocol"},
	),
}
