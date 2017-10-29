package server

import (
	"encoding/json"
	"net/http"
)

func AllowCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(200)
}

func Jsonify(w http.ResponseWriter, res interface{}, cors bool) {
	if cors {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Content-Type", "applicaton/json")
	json.NewEncoder(w).Encode(res)
}

func WriteError(w http.ResponseWriter, err error, status int, cors bool) {
	if cors {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
