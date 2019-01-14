package api

import (
	"fmt"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

// StartStoryBuilderServer starts a story builder server at localhost:<port>
func StartStoryBuilderServer(port int) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/rooms/", RoomHandler)
	http.HandleFunc("/register/", RegistrationHandler)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	return srv
}
