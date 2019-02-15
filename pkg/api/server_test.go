package api

import (
	"testing"

	"github.com/pavelhadzhiev/story-builder/pkg/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStoryBuilderServer(t *testing.T) {
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

func TestServerCreation(t *testing.T) {
	database := &db.SBDatabase{}
	sbServer := NewSBServer(database, 8080)
	expectedAddress := ":8080"

	if sbServer.srv.Addr != expectedAddress {
		t.Errorf("got '%s' want '%s'", sbServer.srv.Addr, expectedAddress)
	}

	if sbServer.Database != database {
		t.Errorf("got '%v' want '%v'", sbServer.Database, database)
	}
}
