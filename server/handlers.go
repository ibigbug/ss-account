package server

import (
	"net/http"

	"github.com/ibigbug/ss-account/user"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		AllowCors(w)
		return
	}
	username := r.FormValue("username")
	backend := r.FormValue("backend")
	port := r.FormValue("port")

	if _, err := user.AddOneUser(backend, username, port); err != nil {
		WriteError(w, err, http.StatusPaymentRequired, true)
	} else {
		u := user.GetManagerByUsername(username)
		Jsonify(w, u, true)
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
	if r.Method == "OPTIONS" {
		AllowCors(w)
		return
	}
	if usage, err := user.GetAllUserUsage(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		Jsonify(w, usage, true)
	}
}
