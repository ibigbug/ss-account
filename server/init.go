package server

import (
	"log"
	"net/http"
	"os"

	"github.com/ibigbug/ss-account/utils"
	"github.com/prometheus/client_golang/prometheus"
)

var logger *log.Logger

func init() {
	if utils.IsaTty(os.Stdout.Fd()) && utils.IsaTty(os.Stderr.Fd()) {
		logger = log.New(os.Stderr, "\x1b[36m[Server]: \x1b[0m", log.LstdFlags)
	} else {
		logger = log.New(os.Stderr, "[Server]: ", log.LstdFlags)
	}

	logger.Println("mounting routers")

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/deregister", deregisterHandler)
	http.HandleFunc("/usage", allManaged)

	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/0.0", renderPaymentForm)
	http.HandleFunc("/payment-handler", paymentHandler)

	http.Handle("/metrics", prometheus.Handler())
}
