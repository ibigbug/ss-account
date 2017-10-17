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

var BytesUploadVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "uploaded_bytes",
		Help: "Bytes uploaded",
	},
	tags,
)

var BytesDownloadVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "downloaded_bytes",
		Help: "Bytes downloaded",
	},
	tags,
)
