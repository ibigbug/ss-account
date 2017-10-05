package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var tags = []string{
	"username",
	"port",
	"backend",
}

var ConnConnectCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "connection_created_total",
		Help: "A counter for connection created",
	},
	tags,
)

var ConnClosedCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "connection_closed_total",
		Help: "A counter for connection closed",
	},
	tags,
)

var BytesUploadVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "uploaded_bytes",
		Help:    "Bytes uploaded",
		Buckets: []float64{.005, .01, .025, .05},
	},
	tags,
)

var BytesDownloadVec = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "downloaded_bytes",
		Help:    "Bytes downloaded",
		Buckets: []float64{.005, .01, .025, .05},
	},
	tags,
)
