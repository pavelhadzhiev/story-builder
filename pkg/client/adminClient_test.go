package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/pavelhadzhiev/story-builder/pkg/config"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Story Builder Game Client test", func() {
	var client *SBClient
	var responseStatusCode int
	var responseBody []byte
	var sbServer *httptest.Server

	username := "user"
	password := "password"
	authHeader := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	room := rooms.NewRoom("roomName", username)
	room.Online = append(room.Online, username)

	player := "test-player"
	room.Online = append(room.Online, player)

	testHandler := func() http.HandlerFunc {
		return func(response http.ResponseWriter, req *http.Request) {
			response.WriteHeader(responseStatusCode)
			response.Write([]byte(responseBody))
		}
	}

	BeforeEach(func() {
		sbServer = httptest.NewServer(testHandler())
		clientConfig := &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader, Room: room.Name}
		client = NewSBClient(clientConfig)
	})

	setupFaultyServer := func() {
		sbServer = httptest.NewUnstartedServer(testHandler())
		clientConfig := &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader, Room: room.Name}
		client = NewSBClient(clientConfig)
	}

	Describe("Ban player", func() {
		Context("When request is valid", func() {
			Context("And the room exists", func() {
				It("should ban the player and not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.BanPlayer(player)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound
					errorMessage := "some error"
					responseBody = []byte(errorMessage)

					err := client.BanPlayer(player)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not ban player: %s", errorMessage)))
				})
			})
		})

		Context("When issuer is not an admin", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user does not have permissions to ban in room \"" + room.Name + "\""))
			})
		})

		Context("When player is already banned", func() {
			It("should return error", func() {
				room.Banned = append(room.Banned, player)
				responseStatusCode = http.StatusConflict

				err := client.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("player is already banned"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Kick player", func() {
		timeLimit := 60
		maxLength := 100
		entriesCount := 0
		room.StartGame(username, timeLimit, maxLength, entriesCount)
		entry := "Test story entry."
		room.AddEntry(entry, username)

		Context("When request is valid", func() {
			Context("And the room exists and a game is started", func() {
				It("should kick the player and not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.KickPlayer(player)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist or game is not started", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound
					errorMessage := "some error"
					responseBody = []byte(errorMessage)

					err := client.KickPlayer(player)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not kick player: %s", errorMessage)))
				})
			})
		})

		Context("When issuer is not an admin", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user does not have permissions to kick in room \"" + room.Name + "\""))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Promote player to admin", func() {
		Context("When request is valid", func() {
			Context("And the room exists", func() {
				It("should promote the player to admin and not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.PromoteAdmin(player)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound
					errorMessage := "some error"
					responseBody = []byte(errorMessage)

					err := client.PromoteAdmin(player)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not promote user: %s", errorMessage)))
				})
			})
		})

		Context("When issuer is not an admin", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user does not have permissions to promote in room \"" + room.Name + "\""))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
