package api

import (
	"net/http"
)

// RoomHandler is an http handler for the story builder's rooms API
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

	w.Write([]byte("Called the rooms API!"))
}
