package client

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStoryBuilderClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

type BadRoom struct {
	Name    bool   `json:"name"`
	Creator string `json:"creator,omitempty"`
}

type BadGame struct {
	Turn int `json:"turn,omitempty"`
}

func TestingHandler(responseBody *[]byte, responseStatusCode *int) http.HandlerFunc {
	return func(response http.ResponseWriter, req *http.Request) {
		response.WriteHeader(*responseStatusCode)
		response.Write([]byte(*responseBody))
	}
}
