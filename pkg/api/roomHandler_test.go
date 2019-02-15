package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
	"github.com/pavelhadzhiev/story-builder/pkg/client"
	"github.com/pavelhadzhiev/story-builder/pkg/config"
)

var _ = Describe("Story Builder Admin Handlers test", func() {
	var sbClient *client.SBClient
	var clientConfig *config.SBConfiguration
	var sbServer *SBServer
	var room *rooms.Room
	var ts *httptest.Server

	username := "username"
	password := "password"
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	roomName := "Test Room"

	BeforeEach(func() {
		// Create a room
		room = rooms.NewRoom(roomName, username)

		// Join the room as the creator
		room.Online = append(room.Online, username)

		// Create the server and add the configured room to it
		sbServer = &SBServer{
			Rooms:  make([]rooms.Room, 0),
			Online: make([]string, 0),
		}
		sbServer.Rooms = append(sbServer.Rooms, *room)
	})

	Describe("Handle room API requests", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.RoomHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Describe("Specifically get all rooms request", func() {
			Context("When request is valid", func() {
				It("should return all rooms and not return error", func() {
					responseRooms, err := sbClient.GetAllRooms()

					Expect(err).ShouldNot(HaveOccurred())
					Expect(responseRooms[0].String()).To(Equal(sbServer.Rooms[0].String()))
				})
			})
		})

		Describe("Specifically create room request", func() {
			Context("When request is valid", func() {
				It("should create a room successfully and not return error", func() {
					otherRoom := rooms.NewRoom("other-room", username)
					err := sbClient.CreateNewRoom(otherRoom)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(len(sbServer.Rooms)).To(Equal(2))
					Expect(sbServer.Rooms[1].String()).To(Equal(otherRoom.String()))
				})
			})

			Context("When there already exists a room with this name", func() {
				It("should not create a new room and return error", func() {
					err := sbClient.CreateNewRoom(room)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + room.Name + "\" already exists"))
					Expect(len(sbServer.Rooms)).To(Equal(1))
					Expect(sbServer.Rooms[0].String()).To(Equal(room.String()))
				})
			})
		})

		Describe("Specifically get a single room request", func() {
			Context("When request is valid", func() {
				It("should return the room and not return error", func() {
					responseRoom, err := sbClient.GetRoom(roomName)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(responseRoom.String()).To(Equal(sbServer.Rooms[0].String()))
				})
			})
			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					responseRoom, err := sbClient.GetRoom(roomName)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist"))
					Expect(responseRoom).To(BeNil())
				})
			})

			Context("When invalid URL is requested", func() {
				It("should return error", func() {
					responseRoom, err := sbClient.GetRoom("invalid/roomname")

					Expect(err).Should(HaveOccurred())
					Expect(responseRoom).To(BeNil())
				})
			})
		})

		Describe("Specifically delete a room request", func() {
			Context("When request is valid", func() {
				It("should delete the room and not return error", func() {
					err := sbClient.DeleteRoom(roomName)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(len(sbServer.Rooms)).To(Equal(0))
				})
			})

			Context("When room does not exist", func() {
				It("should return error", func() {
					sbServer.Rooms = make([]rooms.Room, 0)

					err := sbClient.DeleteRoom(roomName)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist"))
				})
			})

			Context("When user is not the room creator", func() {
				It("should return error", func() {
					sbServer.Rooms[0] = *rooms.NewRoom(roomName, "other-creator")

					err := sbClient.DeleteRoom(roomName)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("user doesn't have permissions to delete this room"))
				})
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Post(ts.URL+"/rooms/"+roomName+"/", "", nil)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})

	Describe("Handle join room request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.JoinRoomHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())

			sbServer.Rooms[0].Online = make([]string, 0)
		})
		Context("When request is valid", func() {
			It("should join the room and not return error", func() {
				err := sbClient.JoinRoom(roomName)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(sbServer.Rooms[0].Online)).To(Equal(1))
				Expect(sbServer.Rooms[0].Online[0]).To(Equal(username))
			})
		})

		Context("When room does not exist", func() {
			It("should return error", func() {
				sbServer.Rooms = make([]rooms.Room, 0)

				err := sbClient.JoinRoom(roomName)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist"))
			})
		})

		Context("When user is banned from the room", func() {
			It("should return error", func() {
				sbServer.Rooms[0].Banned = append(sbServer.Rooms[0].Banned, username)

				err := sbClient.JoinRoom(roomName)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user doesn't have permissions to join this room"))
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.JoinRoom(roomName)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/join-room/" + roomName + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return error", func() {
				err := sbClient.JoinRoom("invalid/roomname")

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Handle leave room request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.LeaveRoomHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader, Room: roomName}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Context("When request is valid", func() {
			It("should leave the room and not return error", func() {
				err := sbClient.LeaveRoom(roomName)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(sbServer.Rooms[0].Online)).To(Equal(0))
			})
		})

		Context("When room does not exist", func() {
			It("should return error", func() {
				sbServer.Rooms = make([]rooms.Room, 0)

				err := sbClient.LeaveRoom(roomName)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("room \"" + roomName + "\" doesn't exist"))
			})
		})

		Context("When user is not in the room", func() {
			It("should return error", func() {
				sbServer.Rooms[0].Online = make([]string, 0)

				err := sbClient.LeaveRoom(roomName)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user is not in room \"" + roomName + "\"."))
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid", Room: roomName}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.LeaveRoom(roomName)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/leave-room/" + roomName + "/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return error", func() {
				err := sbClient.LeaveRoom("invalid/roomname")

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
