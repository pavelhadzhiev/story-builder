package api

import (
	"fmt"
	"net/http"
)

func RoomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Serve the resource.
	case http.MethodPost:
		// Create a new record.
	case http.MethodDelete:
		// Remove the record.
	default:
		// Give an error message.
	}

	fmt.Println("I received a request: ", r)
	w.Write([]byte("Hello from your new room! Here is URL: " + r.URL.Path))
}
