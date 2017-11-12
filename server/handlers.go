package server

import (
	"html/template"
	"net/http"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"

	"github.com/ibigbug/ss-account/config"
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

func renderPaymentForm(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/payment.html"))
	t.Execute(w, config.GetStripeKey())
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("stripeToken")
	params := &stripe.ChargeParams{
		Amount:    5000,
		Currency:  "cny",
		Desc:      "Example charge",
		Statement: "Descriptor",
	}
	params.SetSource(token)
	charge, err := charge.New(params)
	if err != nil {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(err.Error()))
	} else {
		if charge.Paid {
			if port, err := user.AddOneUser("127.0.0.1:3390", "", ""); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Please contact customer service"))
			} else {
				w.Write([]byte(port))
			}
		} else {
			w.WriteHeader(http.StatusPaymentRequired)
			w.Write([]byte(charge.FailMsg))
		}
	}
}
