package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/pavelhadzhiev/story-builder/pkg/api/game"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
	"github.com/pavelhadzhiev/story-builder/pkg/client"
	"github.com/pavelhadzhiev/story-builder/pkg/config"
	"github.com/pavelhadzhiev/story-builder/pkg/db/dbfakes"
)

var _ = Describe("Story Builder Admin Handlers test", func() {
	var sbClient *client.SBClient
	var clientConfig *config.SBConfiguration
	var sbServer *SBServer
	var room *rooms.Room
	var database *dbfakes.FakeUserDatabase
	var ts *httptest.Server

	username := "username"
	password := "password"
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	roomName := "Test Room"

	player := "test-player"

	timeLimit := 60
	maxLength := 100
	entriesCount := 0
	entry := "Test story entry."

	BeforeEach(func() {
		// Create a room
		room = rooms.NewRoom(roomName, username)

		// Join the room as the creator
		room.Online = append(room.Online, username)

		// Put a non-admin player in the room
		room.Online = append(room.Online, player)

		// Configure a game
		room.StartGame(username, timeLimit, maxLength, entriesCount)

		// Fake a user database and make it return true when the argument matches the non-admin player
		database = &dbfakes.FakeUserDatabase{}
		database.UserExistsStub = func(username string) (bool, error) {
			if username == player {
				return true, nil
			}
			return false, nil
		}

		// Create the server and add the configured room to it
		sbServer = &SBServer{
			Database: database,
			Rooms:    make([]rooms.Room, 0),
			Online:   make([]string, 0),
		}
		sbServer.Rooms = append(sbServer.Rooms, *room)
	})

	Describe("Handle gameplay requests", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.GameplayHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Describe("Specifically get game request", func() {
			Context("When request is valid", func() {
				Context("And game is not finished", func() {
					It("should return the game and not return error", func() {
						responseGame, err := sbClient.GetGame()

						Expect(err).ShouldNot(HaveOccurred())
						Expect(responseGame.String()).To(Equal(sbServer.Rooms[0].GetGame().String()))
					})
				})

				Context("And game is finished", func() {
					Context("And no new games are started", func() {
						It("should return the game and not return error", func() {
							sbServer.Rooms[0].GetGame().Finished = true
							responseGame, err := sbClient.GetGame()

							Expect(err).ShouldNot(HaveOccurred())
							Expect(responseGame.String()).To(Equal(sbServer.Rooms[0].GetGame().String()))
						})
					})

					Context("And a new game is started", func() {
						It("should return the game and not return error", func() {
							sbServer.Rooms[0].GetGame().Finished = true
							sbServer.Rooms[0].StartGame(username, timeLimit, maxLength, entriesCount)
							responseGame, err := sbClient.GetGame()

							Expect(err).ShouldNot(HaveOccurred())
							Expect(responseGame.String()).To(Equal(sbServer.Rooms[0].GetGame().String()))
							Expect(responseGame.Finished).ToNot(BeTrue())
						})
					})
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					responseGame, err := sbClient.GetGame()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist or no games have been started"))
					Expect(responseGame).To(BeNil())
				})
			})

			Context("When no games have ever been started", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 1)
					sbServer.Rooms[0] = *rooms.NewRoom(roomName, username)

					responseGame, err := sbClient.GetGame()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist or no games have been started"))
					Expect(responseGame).To(BeNil())
				})
			})

			Context("When invalid URL is requested", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: "invalid/roomname"}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					responseGame, err := sbClient.GetGame()

					Expect(err).Should(HaveOccurred())
					Expect(responseGame).To(BeNil())
				})
			})
		})

		Describe("Specifically add entry request", func() {
			Context("When request is valid", func() {
				Context("And it is indeed the user's turn", func() {
					It("should successfully add the entry to the game and not return error", func() {
						entry := "added by admin"
						err := sbClient.AddEntry(entry)

						Expect(err).ShouldNot(HaveOccurred())

						// New entry should be inserted in the game
						Expect(sbServer.Rooms[0].GetGame().String()).To(ContainSubstring(entry))

						// Turn should have incremented
						time.Sleep(100 * time.Millisecond)
						Expect(sbServer.Rooms[0].GetGame().Turn).To(Equal(player))
					})
				})

				Context("And it is not the the user's turn", func() {
					It("should fail to add the entry and return error", func() {
						sbServer.Rooms[0].GetGame().Turn = player
						entry := "added by admin"
						err := sbClient.AddEntry(entry)

						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("illegal entry: There was an error while adding your entry: invalid entry - not this player's turn"))

						// New entry should be inserted in the game
						Expect(sbServer.Rooms[0].GetGame().String()).ToNot(ContainSubstring(entry))

						// Turn should not have incremented
						time.Sleep(100 * time.Millisecond)
						Expect(sbServer.Rooms[0].GetGame().Turn).To(Equal(player))
					})
				})
			})

			Context("When there isn't a started game", func() {
				It("should return error", func() {
					sbServer.Rooms[0].GetGame().Finished = true

					err := sbClient.AddEntry(entry)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist or no games have been started"))
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					err := sbClient.AddEntry(entry)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist or no games have been started"))
				})
			})

			Context("When no entry is provided", func() {
				It("should return error", func() {
					err := sbClient.AddEntry("")

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("missing Entry-Text header from request"))
				})
			})

			Context("When an invalid authorization header is provided", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.AddEntry(entry)

					Expect(err).Should(HaveOccurred())
				})
			})

			Context("When invalid URL is requested", func() {
				It("should return error status code HTTP 400", func() {
					resp, err := http.Get(ts.URL + "/invalid/roomname")

					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				})
			})
		})
	})

	Describe("Handle game management requests", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.ManageGamesHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Describe("Specifically start game request", func() {
			BeforeEach(func() {
				sbServer.Rooms[0].GetGame().Finished = true
			})
			Context("When request is valid", func() {
				Context("And previous game is not finished", func() {
					It("should return error", func() {
						sbServer.Rooms[0].GetGame().Finished = false
						err := sbClient.StartGame(timeLimit, maxLength, entriesCount)

						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("a game is already running in \"" + roomName + "\""))
					})
				})

				Context("And previous game is finished", func() {
					It("should start game succesffuly and not return error", func() {
						err := sbClient.StartGame(timeLimit, maxLength, entriesCount)

						Expect(err).ShouldNot(HaveOccurred())
						Expect(sbServer.Rooms[0].GetGame().Finished).To(BeFalse())
					})
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					err := sbClient.StartGame(timeLimit, maxLength, entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist"))
				})
			})

			Context("When user doesn't have permissions", func() {
				It("should return error", func() {
					sbServer.Rooms[0].Admins = make([]string, 0)

					err := sbClient.StartGame(timeLimit, maxLength, entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot start game: requires admin access"))
				})
			})

			Context("When an invalid authorization header is provided", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.StartGame(timeLimit, maxLength, entriesCount)

					Expect(err).Should(HaveOccurred())
				})
			})

			Context("When invalid URL is requested", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: "invalid/roomname"}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.StartGame(timeLimit, maxLength, entriesCount)

					Expect(err).Should(HaveOccurred())
				})
			})
		})

		Describe("Specifically end game request", func() {
			BeforeEach(func() {
				sbServer.Rooms[0].GetGame().Finished = false
				entriesCount = 1
			})
			Context("When request is valid", func() {
				Context("And game is not finished", func() {
					It("should return error", func() {
						err := sbClient.EndGame(entriesCount)

						Expect(err).ShouldNot(HaveOccurred())
						Expect(sbServer.Rooms[0].GetGame().EntriesLeft).To(Equal(entriesCount))
					})
				})

				Context("And game is finished", func() {
					It("should return the game and not return error", func() {
						sbServer.Rooms[0].GetGame().Finished = true
						err := sbClient.EndGame(entriesCount)

						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("there is no running game in \"" + roomName + "\""))
					})
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					err := sbClient.EndGame(entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist"))
				})
			})

			Context("When user doesn't have permissions", func() {
				It("should return error", func() {
					sbServer.Rooms[0].Admins = make([]string, 0)

					err := sbClient.EndGame(entriesCount)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot end game: requires admin access"))
				})
			})

			Context("When an invalid authorization header is provided", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.EndGame(entriesCount)

					Expect(err).Should(HaveOccurred())
				})
			})

			Context("When invalid URL is requested", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: "invalid/roomname"}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.EndGame(entriesCount)

					Expect(err).Should(HaveOccurred())
				})
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/manage-games/" + room.Name + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})

	Describe("Handle vote requests", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.VoteHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Describe("Specifically trigger vote request", func() {
			BeforeEach(func() {
				sbServer.Rooms[0].GetGame().VoteKick = nil
			})
			Context("When request is valid", func() {
				Context("And game is not finished", func() {
					It("should trigger voting successfully and not return error", func() {
						err := sbClient.TriggerVoteKick(player)

						Expect(err).ShouldNot(HaveOccurred())
						Expect(sbServer.Rooms[0].GetGame().VoteKick).ToNot(BeNil())
						Expect(sbServer.Rooms[0].GetGame().VoteKick.Player).To(Equal(player))
						Expect(sbServer.Rooms[0].GetGame().VoteKick.Issuer).To(Equal(username))
					})
				})

				Context("And game is finished", func() {
					It("should not trigger voting and return error", func() {
						sbServer.Rooms[0].GetGame().Finished = true
						err := sbClient.TriggerVoteKick(player)

						Expect(err).Should(HaveOccurred())
						Expect(sbServer.Rooms[0].GetGame().VoteKick).To(BeNil())
					})
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					err := sbClient.TriggerVoteKick(player)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not trigger vote: %s", "Room \""+roomName+"\" doesn't exist or no games have been started.")))
				})
			})

			Context("When no games have been started", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 1)
					sbServer.Rooms[0] = *rooms.NewRoom("new room", username)

					err := sbClient.TriggerVoteKick(player)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not trigger vote: %s", "Room \""+roomName+"\" doesn't exist or no games have been started.")))
				})
			})

			Context("When there is already an ongoing vote", func() {
				It("should return error", func() {
					sbClient.TriggerVoteKick(player)
					err := sbClient.TriggerVoteKick(player)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("there is already an ongoing vote"))
				})
			})

			Context("When an invalid authorization header is provided", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.TriggerVoteKick(player)

					Expect(err).Should(HaveOccurred())
				})
			})

			Context("When invalid URL is requested", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: "invalid/roomname"}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.TriggerVoteKick(player)

					Expect(err).Should(HaveOccurred())
				})
			})
		})

		Describe("Specifically submit vote request", func() {
			BeforeEach(func() {
				sbServer.Rooms[0].GetGame().TriggerVoteKick(username, player, 1, 60)
			})
			Context("When request is valid", func() {
				Context("And game is not finished", func() {
					Context("And there is ongoing vote", func() {
						It("should submit vote successfully and not return error", func() {
							err := sbClient.SubmitVote()

							// No errors should occur
							Expect(err).ShouldNot(HaveOccurred())

							// Vote should be counted
							Expect(sbServer.Rooms[0].GetGame().VoteKick.Count).To(Equal(1))
						})
					})
					Context("And there is no ongoing vote", func() {
						It("should not submit vote and return error", func() {
							sbServer.Rooms[0].GetGame().VoteKick = nil

							err := sbClient.SubmitVote()

							Expect(err).Should(HaveOccurred())
							Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("cannot vote: %s", "There are no ongoing votes.")))
						})
					})
				})

				Context("And game is finished", func() {
					It("should not submit vote and return error", func() {
						sbServer.Rooms[0].GetGame().Finished = true
						err := sbClient.SubmitVote()

						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("cannot vote: %s", "There is no running game.")))
					})
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					err := sbClient.SubmitVote()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("cannot vote: %s", "Room \""+roomName+"\" doesn't exist or no games have been started.")))
				})
			})

			Context("When no games have been started", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 1)
					sbServer.Rooms[0] = *rooms.NewRoom("new room", username)

					err := sbClient.SubmitVote()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("cannot vote: %s", "Room \""+roomName+"\" doesn't exist or no games have been started.")))
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					// Increase treshold to test double voting
					sbServer.Rooms[0].GetGame().VoteKick = game.NewVoteKick(username, player, 2, 60)

					sbClient.SubmitVote()
					err := sbClient.SubmitVote()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot vote: user has already voted once"))
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms[0].GetGame().Players = make([]string, 0)
					err := sbClient.SubmitVote()

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot vote: user is not part of the game"))
				})
			})

			Context("When an invalid authorization header is provided", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.SubmitVote()

					Expect(err).Should(HaveOccurred())
				})
			})

			Context("When invalid URL is requested", func() {
				It("should return error", func() {
					clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: "invalid/roomname"}
					sbClient = client.NewTestSBClient(clientConfig, ts.Client())

					err := sbClient.SubmitVote()

					Expect(err).Should(HaveOccurred())
				})
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/vote/" + room.Name + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})
})
