package client

import "fmt"

type ResponseError struct {
	URL        string
	StatusCode int
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("URL: %s, Status: %v;", e.URL, e.StatusCode)
}
