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
		room.AddEntry(entry, username)

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

	Describe("Handle ban request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.BanHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Context("When request is valid", func() {
			It("player should be banned and request not return error", func() {
				sbServer.Rooms[0].GetGame().Turn = player
				err := sbClient.BanPlayer(player)

				// Response should be 200
				Expect(err).Should(BeNil())

				// Banned player should be kicked from the game
				Expect(len(sbServer.Rooms[0].GetGame().Players)).To(Equal(1))
				Expect(sbServer.Rooms[0].GetGame().Players[0]).To(Equal(username))

				// Banned player should be removed from turn
				Expect(sbServer.Rooms[0].GetGame().Turn).To(Equal(username))

				// Banned player should be in the banned list of the room
				Expect(len(sbServer.Rooms[0].Banned)).To(Equal(1))
				Expect(sbServer.Rooms[0].Banned[0]).To(Equal(player))
			})
		})

		Context("When issuer is not an admin", func() {
			It("should return error", func() {
				sbServer.Rooms[0].Admins[0] = ""

				err := sbClient.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user does not have permissions to ban in room \"" + room.Name + "\""))
			})
		})

		Context("When player is already banned", func() {
			It("should return error", func() {
				sbServer.Rooms[0].Banned = append(sbServer.Rooms[0].Banned, player)

				err := sbClient.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("player is already banned"))
			})
		})

		Context("When room does not exist", func() {
			It("should return error", func() {
				sbServer.Rooms = make([]rooms.Room, 0)

				err := sbClient.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not ban player: %s", "Room \""+room.Name+"\" doesn't exist.")))
			})
		})

		Context("When user does not exist", func() {
			It("should return error", func() {
				database.UserExistsStub = func(username string) (bool, error) {
					return false, nil
				}

				err := sbClient.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not ban player: %s", "User \""+player+"\" doesn't exist.")))
			})
		})

		Context("When user lookup fails", func() {
			It("should return error", func() {
				database.UserExistsStub = func(username string) (bool, error) {
					return false, errors.New("")
				}

				err := sbClient.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.BanPlayer(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/admin/ban/" + room.Name + "/" + player + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return error if invalid URL", func() {
				err := sbClient.BanPlayer("not/a/valid/playername")

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Handle kick request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.KickHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Context("When request is valid", func() {
			It("player should be kicked from the game and request not return error", func() {
				sbServer.Rooms[0].GetGame().Turn = player
				err := sbClient.KickPlayer(player)

				// Response should be 200
				Expect(err).Should(BeNil())

				// Kicked player should be kicked from the game
				Expect(len(sbServer.Rooms[0].GetGame().Players)).To(Equal(1))
				Expect(sbServer.Rooms[0].GetGame().Players[0]).To(Equal(username))

				// Kicked player should be removed from turn
				Expect(sbServer.Rooms[0].GetGame().Turn).To(Equal(username))
			})
		})

		Context("When issuer is not an admin", func() {
			It("should return error", func() {
				sbServer.Rooms[0].Admins[0] = ""

				err := sbClient.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user does not have permissions to kick in room \"" + room.Name + "\""))
			})
		})

		Context("When player is not in the game", func() {
			It("should return error", func() {
				sbServer.Rooms[0].GetGame().Players = make([]string, 0)

				err := sbClient.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not kick player: %s", "The user to kick is not in the game.")))
			})
		})

		Context("When room does not exist", func() {
			It("should return error", func() {
				sbServer.Rooms = make([]rooms.Room, 0)

				err := sbClient.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not kick player: %s", "Room \""+roomName+"\" doesn't exist.")))
			})
		})

		Context("When game does not exist or is finished", func() {
			It("should return error", func() {
				sbServer.Rooms[0].GetGame().Finished = true

				err := sbClient.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not kick player: %s", "There is no running game in room \""+roomName+"\".")))
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.KickPlayer(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/admin/kick/" + room.Name + "/" + player + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return error if invalid URL", func() {
				err := sbClient.KickPlayer("not/a/valid/playername")

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Handle promote admin request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.PromoteAdminHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Context("When request is valid", func() {
			It("player should be promoted to admin for the room and request should not return error", func() {
				err := sbClient.PromoteAdmin(player)

				// Response should be 200
				Expect(err).Should(BeNil())

				// Promoted player should be an admin in the room
				Expect(len(sbServer.Rooms[0].Admins)).To(Equal(2))
				Expect(sbServer.Rooms[0].Admins[1]).To(Equal(player))
			})
		})

		Context("When issuer is not an admin", func() {
			It("should return error", func() {
				sbServer.Rooms[0].Admins[0] = ""

				err := sbClient.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user does not have permissions to promote in room \"" + room.Name + "\""))
			})
		})

		Context("When room does not exist", func() {
			It("should return error", func() {
				sbServer.Rooms = make([]rooms.Room, 0)

				err := sbClient.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not promote user: %s", "Room \""+roomName+"\" doesn't exist.")))
			})
		})

		Context("When user does not exist", func() {
			It("should return error", func() {
				database.UserExistsStub = func(username string) (bool, error) {
					return false, nil
				}

				err := sbClient.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("could not promote user: %s", "User \""+player+"\" doesn't exist.")))
			})
		})

		Context("When user lookup fails", func() {
			It("should return error", func() {
				database.UserExistsStub = func(username string) (bool, error) {
					return false, errors.New("")
				}

				err := sbClient.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.PromoteAdmin(player)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/admin/" + room.Name + "/" + player + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return error if invalid URL", func() {
				err := sbClient.PromoteAdmin("not/a/valid/playername")

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
