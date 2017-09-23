package server

import "net/http"

func init() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/deregister", deregisterHandler)
	http.HandleFunc("/usage", allManaged)
}
