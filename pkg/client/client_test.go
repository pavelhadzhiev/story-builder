package client

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSmClient(t *testing.T) {
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
