package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ibigbug/ss-account/user"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	backend := r.FormValue("backend")
	portStr := r.FormValue("port")
	var port int
	if p, err := strconv.Atoi(portStr); err == nil {
		port = p
	}

	if port, err := user.AddOneUser(backend, username, port); err != nil {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(fmt.Sprintf("%d", port)))
	}
}

func deregisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	m := user.GetManagerByUsername(username)
	if m == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("no such user"))
		return
	}

	m.Stop()
	user.DefaultManaged.Remove(m)
	w.WriteHeader(http.StatusNoContent)
}

func allManaged(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(user.DefaultManaged.String()))
}
