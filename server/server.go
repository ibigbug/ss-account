package server

import (
	"net/http"

	"github.com/ibigbug/ss-account/user"
)

// Start runs the web console panel
func Start(bind string) error {
	logger.Println("initializing user")
	user.Init()

	logger.Println("listening:", bind)
	return http.ListenAndServe(bind, nil)
}
