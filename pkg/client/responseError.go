package client

import "fmt"

// ResponseError represents an error that will be returned by the client, holding context of an error response from the server.
type ResponseError struct {
	URL        string
	StatusCode int
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("URL: %s, Status: %v;", e.URL, e.StatusCode)
}
