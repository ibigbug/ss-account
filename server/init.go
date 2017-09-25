package server

import "net/http"
import "github.com/prometheus/client_golang/prometheus"

func init() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/deregister", deregisterHandler)
	http.HandleFunc("/usage", allManaged)

	prometheus.MustRegister(connConnectCounter, connClosedCounter, bytesDownloadVec, bytesUploadVec)
	http.Handle("/metrics", prometheus.Handler())
}
