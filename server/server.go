package server

import (
	"log"
	"net/http"

	"github.com/ibigbug/ss-account/config"
)

// Start runs the web console panel
func Start(bind string) error {

	log.Println("listening:", bind)
	return http.ListenAndServe(bind, nil)
}

// AddOneUser register a user with specific backend
// a random port will be used upon the port range if
// no port given
func AddOneUser(backend, username string, port int) (p int, err error) {
	if port == 0 {
		p = config.GetRandomPort()
	} else {
		p = port
	}

	m := Manager{
		Username: username,
		Backend:  backend,
		Port:     port,
	}

	if err = m.Bind(); err != nil {
		return
	}

	if err = m.Start(); err != nil {
		return
	}

	DefaultManaged.Add(&m)

	return
}
