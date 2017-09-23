package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	backend := r.FormValue("backend")
	portStr := r.FormValue("port")
	var port int
	if p, err := strconv.Atoi(portStr); err == nil {
		port = p
	}

	if port, err := AddOneUser(backend, username, port); err != nil {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(fmt.Sprintf("%d", port)))
	}
}

func deregisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	m := GetManagerByUsername(username)
	if m == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("no such user"))
		return
	}

	m.Stop()
	DefaultManaged.Remove(m)
	w.WriteHeader(http.StatusNoContent)
}

func allManaged(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(DefaultManaged.String()))
}
