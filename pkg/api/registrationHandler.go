package api

import (
	"net/http"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register/" {
		w.WriteHeader(404)
		return
	}
	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("Let's say you've registered.\n"))
	default:
		w.WriteHeader(405)
	}
}