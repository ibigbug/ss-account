package server

import (
	"encoding/json"
	"net/http"
)

func Jsonify(w http.ResponseWriter, res interface{}, cors bool) {
	if cors {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Content-Type", "applicaton/json")
	json.NewEncoder(w).Encode(res)
}
