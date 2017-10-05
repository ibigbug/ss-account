package server

import (
	"log"
	"net/http"
)

// Start runs the web console panel
func Start(bind string) error {

	log.Println("listening:", bind)
	return http.ListenAndServe(bind, nil)
}
