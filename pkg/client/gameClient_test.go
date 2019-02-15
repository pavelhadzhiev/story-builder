package client

import (
	"encoding/base64"
	"encoding/json"
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
	testHandler := TestingHandler(&responseBody, &responseStatusCode)

	username := "user"
	password := "password"
	authHeader := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	room := rooms.NewRoom("roomName", username)
	room.Online = append(room.Online, username)
	timeLimit := 60
	maxLength := 100
	entriesCount := 0
	room.StartGame(username, timeLimit, maxLength, entriesCount)
	entry := "Test story entry."
	room.AddEntry(entry, username)
	gameObj := room.GetGame()

	BeforeEach(func() {
		sbServer = httptest.NewServer(testHandler)
		clientConfig := &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader, Room: room.Name}
		client = NewSBClient(clientConfig)
	})

	setupFaultyServer := func() {
		sbServer = httptest.NewUnstartedServer(testHandler)
		clientConfig := &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader, Room: room.Name}
		client = NewSBClient(clientConfig)
	}

	Describe("Get game", func() {
		Context("When request is valid", func() {
			Context("And room exists", func() {
				It("should return the game successfully", func() {
					responseStatusCode = http.StatusOK
					responseBody, _ = json.Marshal(gameObj)

					responseGame, err := client.GetGame()

					Expect(err).ShouldNot(HaveOccurred())
					Expect(responseGame.String()).To(Equal(gameObj.String()))
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound
					responseBody, _ = json.Marshal(gameObj)

					nonExistentRoom := "non-existent-room"
					client.config.Room = nonExistentRoom
					responseGame, err := client.GetGame()

					Expect(err).Should(HaveOccurred())
					Expect(responseGame).To(BeNil())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistentRoom)))
				})
			})
		})

		Context("When invalid game is returned by SB server", func() {
			It("should return error", func() {
				badGame := &BadGame{Turn: 2}

				responseBody, _ = json.Marshal(badGame)
				responseStatusCode = http.StatusOK

				responseGame, err := client.GetGame()

				Expect(err).Should(HaveOccurred())
				Expect(responseGame).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseBody, _ = json.Marshal(gameObj)
				responseStatusCode = http.StatusCreated

				responseGame, err := client.GetGame()

				Expect(err).Should(HaveOccurred())
				Expect(responseGame).To(BeNil())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseBody, _ = json.Marshal(gameObj)
				responseStatusCode = http.StatusOK

				responseGame, err := client.GetGame()

				Expect(err).Should(HaveOccurred())
				Expect(responseGame).To(BeNil())
			})
		})
	})

	Describe("Add entry to game", func() {
		Context("When request is valid", func() {
			Context("And room exists", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.AddEntry(entry)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound

					nonExistentRoom := "non-existent-room"
					client.config.Room = nonExistentRoom
					err := client.AddEntry(entry)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistentRoom)))
				})
			})
		})

		Context("When entry header is missing", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusBadRequest

				err := client.AddEntry(entry)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("missing Entry-Text header from request"))
			})
		})

		Context("When entry was illegal", func() {
			It("should return error", func() {
				errorMessage := "some error according to game rules"
				responseBody = []byte(errorMessage)
				responseStatusCode = http.StatusForbidden

				err := client.AddEntry(entry)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("illegal entry: %s", errorMessage)))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.AddEntry(entry)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.AddEntry(entry)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Start a new game", func() {
		Context("When request is valid", func() {
			Context("And room exists", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.StartGame(timeLimit, maxLength, entriesCount)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound

					nonExistentRoom := "non-existent-room"
					client.config.Room = nonExistentRoom
					err := client.StartGame(timeLimit, maxLength, entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistentRoom)))
				})
			})
		})

		Context("When a game is already running", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusConflict

				err := client.StartGame(timeLimit, maxLength, entriesCount)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("a game is already running in \"" + room.Name + "\""))
			})
		})

		Context("When issuer does not have admin access", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.StartGame(timeLimit, maxLength, entriesCount)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot start game: requires admin access"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.StartGame(timeLimit, maxLength, entriesCount)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.StartGame(timeLimit, maxLength, entriesCount)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("With illegal input", func() {
			Context("Without being in a room", func() {
				It("should return error", func() {
					client.wipeRoom()
					err := client.StartGame(timeLimit, maxLength, entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot start game: requires user to be joined in the room"))
				})
			})

			Context("With illegal time left setting", func() {
				It("should return error", func() {
					err := client.StartGame(-1, maxLength, entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot start game: negative time limit value"))
				})
			})

			Context("With illegal max length setting", func() {
				It("should return error", func() {
					err := client.StartGame(timeLimit, -1, entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot start game: negative max length value"))
				})
			})

			Context("With illegal entries count setting", func() {
				It("should return error", func() {
					err := client.StartGame(timeLimit, maxLength, -1)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot start game: negative entries value"))
				})
			})
		})
	})

	Describe("End a running game", func() {
		entriesCount = 1
		Context("When request is valid", func() {
			Context("And room exists", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusAccepted

					err := client.EndGame(entriesCount)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound

					nonExistentRoom := "non-existent-room"
					client.config.Room = nonExistentRoom
					err := client.EndGame(entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistentRoom)))
				})
			})
		})

		Context("When a game is already running", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusConflict

				err := client.EndGame(entriesCount)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("there is no running game in \"" + room.Name + "\""))
			})
		})

		Context("When issuer does not have admin access", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.EndGame(entriesCount)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot end game: requires admin access"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.EndGame(entriesCount)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusAccepted

				err := client.EndGame(entriesCount)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("With illegal input", func() {
			Context("Without being in a room", func() {
				It("should return error", func() {
					client.wipeRoom()
					err := client.EndGame(entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot end game: requires user to be joined in the room"))
				})
			})

			Context("With illegal entries count setting", func() {
				It("should return error", func() {
					err := client.EndGame(-1)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot end game: negative entries value"))
				})
			})
		})
	})

	Describe("Trigger a vote to kick a player", func() {
		playerToKick := "some-other-player"
		room := rooms.NewRoom("roomName", username)
		room.Online = append(room.Online, username)
		room.Online = append(room.Online, playerToKick)
		room.StartGame(username, timeLimit, maxLength, entriesCount)

		Context("When request is valid", func() {
			Context("And room exists", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusAccepted

					err := client.TriggerVoteKick(playerToKick)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					errorMessage := "some server error"
					responseBody = []byte(errorMessage)
					responseStatusCode = http.StatusNotFound

					nonExistentRoom := "non-existent-room"
					client.config.Room = nonExistentRoom
					err := client.TriggerVoteKick(playerToKick)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not trigger vote: %s", errorMessage)))
				})
			})
		})

		Context("When a vote is already running", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusConflict

				err := client.TriggerVoteKick(playerToKick)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("there is already an ongoing vote"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.TriggerVoteKick(playerToKick)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusAccepted

				err := client.TriggerVoteKick(playerToKick)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Submit a vote to kick a player", func() {
		playerToKick := "some-other-player"
		room := rooms.NewRoom("roomName", username)
		room.Online = append(room.Online, username)
		room.Online = append(room.Online, playerToKick)
		room.StartGame(username, timeLimit, maxLength, entriesCount)
		room.GetGame().TriggerVoteKick(username, playerToKick, 0.5, timeLimit)

		Context("When request is valid", func() {
			Context("And room exists", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.SubmitVote()

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And room does not exist", func() {
				It("should return error", func() {
					errorMessage := "some server error"
					responseBody = []byte(errorMessage)
					responseStatusCode = http.StatusNotFound

					nonExistentRoom := "non-existent-room"
					client.config.Room = nonExistentRoom
					err := client.SubmitVote()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("cannot vote: %s", errorMessage)))
				})
			})
		})

		Context("When the user is not part of the game", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.SubmitVote()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot vote: user is not part of the game"))
			})
		})

		Context("When the user has already voted for this campaign", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusConflict

				err := client.SubmitVote()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot vote: user has already voted once"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.SubmitVote()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.SubmitVote()

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
