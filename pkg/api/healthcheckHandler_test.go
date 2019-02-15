package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
	"github.com/pavelhadzhiev/story-builder/pkg/client"
	"github.com/pavelhadzhiev/story-builder/pkg/config"
	"github.com/pavelhadzhiev/story-builder/pkg/config/configfakes"
	"github.com/pavelhadzhiev/story-builder/pkg/db/dbfakes"
)

var _ = Describe("Story Builder Admin Handlers test", func() {
	var sbClient *client.SBClient
	var clientConfig *config.SBConfiguration
	var sbServer *SBServer
	var room *rooms.Room
	var database *dbfakes.FakeUserDatabase
	var ts *httptest.Server

	configurator := &configfakes.FakeSBConfigurator{}
	configurator.LoadReturns(clientConfig, nil)
	configurator.SaveStub = func(config *config.SBConfiguration) error {
		clientConfig = config
		sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		return nil
	}
	username := "username"
	password := "password"
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	roomName := "Test Room"

	BeforeEach(func() {
		// Create a room
		room = rooms.NewRoom(roomName, username)

		// Join the room as the creator
		room.Online = append(room.Online, username)

		// Fake a user database and make it return true when the argument matches the non-admin player
		database = &dbfakes.FakeUserDatabase{}
		database.UserExistsReturns(true, nil)

		// Create the server and add the configured room to it
		sbServer = &SBServer{
			Database: database,
			Rooms:    make([]rooms.Room, 0),
			Online:   make([]string, 0),
		}
		sbServer.Rooms = append(sbServer.Rooms, *room)

		ts = httptest.NewServer(http.HandlerFunc(sbServer.HealthcheckHandler))

		clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
		sbClient = client.NewTestSBClient(clientConfig, ts.Client())

	})
	Describe("Handle healthcheck request", func() {
		Context("When request is valid", func() {
			It("should not return error", func() {
				err := sbClient.HealthCheck(configurator)

				// Response should be 200
				Expect(err).Should(BeNil())
			})
		})

		Context("When the room does not exist", func() {
			It("should return error", func() {
				sbServer.Rooms = make([]rooms.Room, 0)

				err := sbClient.HealthCheck(configurator)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist"))
			})
		})

		Context("When player is not in the room", func() {
			It("should return error", func() {
				sbServer.Rooms[0].Online = make([]string, 0)

				err := sbClient.HealthCheck(configurator)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("player is not in room \"" + roomName + "\""))
			})
		})

		Context("When authentication fails", func() {
			It("should return error", func() {
				database.LoginUserReturns(errors.New("error"))

				err := sbClient.HealthCheck(configurator)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("authentication failed")))
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.HealthCheck(configurator)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/healthcheck/" + roomName + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return error if invalid URL", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: "invalid/roomname"}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.HealthCheck(configurator)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
