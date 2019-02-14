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

var _ = Describe("Story Builder Room Client test", func() {
	var client *SBClient
	var responseStatusCode int
	var responseBody []byte
	var sbServer *httptest.Server

	username := "user"
	password := "password"
	authHeader := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	room := &rooms.Room{Name: "roomName", Creator: "creator"}
	roomList := []rooms.Room{*room, rooms.Room{Name: "roomName2", Creator: "creator2"}}

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

	Describe("Get all rooms", func() {
		Context("When configuration is valid", func() {
			It("should return all rooms successfully", func() {
				responseStatusCode = http.StatusOK
				responseBody, _ = json.Marshal(roomList)

				responseRooms, err := client.GetAllRooms()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responseRooms).To(Equal(roomList))
			})
		})

		Context("When invalid rooms are returned by SB server", func() {
			It("should return error", func() {
				badRoomList := []BadRoom{
					BadRoom{Name: true, Creator: "creator"},
					BadRoom{Name: false, Creator: "creator2"},
				}

				responseBody, _ = json.Marshal(badRoomList)
				responseStatusCode = http.StatusOK

				responseRooms, err := client.GetAllRooms()

				Expect(err).Should(HaveOccurred())
				Expect(responseRooms).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseBody, _ = json.Marshal(roomList)
				responseStatusCode = http.StatusCreated

				responseRooms, err := client.GetAllRooms()

				Expect(err).Should(HaveOccurred())
				Expect(responseRooms).To(BeNil())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseBody, _ = json.Marshal(roomList)
				responseStatusCode = http.StatusOK

				responseRooms, err := client.GetAllRooms()

				Expect(err).Should(HaveOccurred())
				Expect(responseRooms).To(BeNil())
			})
		})
	})

	Describe("Create a room", func() {
		Context("When configuration is valid", func() {
			It("should not return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.CreateNewRoom(room)

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("When the room to be created already exists", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusConflict

				alreadyExistingRoomName := "already-existing-room"
				alreadyExistingRoom := &rooms.Room{Name: alreadyExistingRoomName}
				err := client.CreateNewRoom(alreadyExistingRoom)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" already exists", alreadyExistingRoomName)))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseBody, _ = json.Marshal(roomList)
				responseStatusCode = http.StatusOK

				err := client.CreateNewRoom(room)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusCreated

				err := client.CreateNewRoom(room)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Get a single room", func() {
		Context("When configuration is valid", func() {
			It("should return the right room successfully", func() {
				responseStatusCode = http.StatusOK
				responseBody, _ = json.Marshal(room)

				responseRoom, err := client.GetRoom(room.Name)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(responseRoom).To(Equal(room))
			})
		})

		Context("When the searched room doesn't exist", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusNotFound

				nonExistingRoom := "non-existing room"
				responseRoom, err := client.GetRoom(nonExistingRoom)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistingRoom)))
				Expect(responseRoom).To(BeNil())
			})
		})

		Context("When an invalid room are returned by SB server", func() {
			It("should return error", func() {
				badRoom := BadRoom{Name: true, Creator: "creator"}

				responseBody, _ = json.Marshal(badRoom)
				responseStatusCode = http.StatusOK

				responseRoom, err := client.GetRoom(room.Name)

				Expect(err).Should(HaveOccurred())
				Expect(responseRoom).To(BeNil())
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseBody, _ = json.Marshal(room)
				responseStatusCode = http.StatusCreated

				responseRoom, err := client.GetRoom(room.Name)

				Expect(err).Should(HaveOccurred())
				Expect(responseRoom).To(BeNil())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseBody, _ = json.Marshal(roomList)
				responseStatusCode = http.StatusCreated

				responseRoom, err := client.GetRoom(room.Name)

				Expect(err).Should(HaveOccurred())
				Expect(responseRoom).To(BeNil())
			})
		})
	})

	Describe("Delete a room", func() {
		Context("When configuration is valid", func() {
			It("should not return error", func() {
				responseStatusCode = http.StatusNoContent

				err := client.DeleteRoom(room.Name)

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("When the user doesn't have permission to delete the room", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.DeleteRoom(room.Name)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user doesn't have permissions to delete this room"))
			})
		})

		Context("When the searched room doesn't exist", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusNotFound

				nonExistingRoom := "non-existing room"
				err := client.DeleteRoom(nonExistingRoom)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistingRoom)))
			})
		})

		Context("When an invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusOK

				err := client.DeleteRoom(room.Name)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusNoContent

				err := client.DeleteRoom(room.Name)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Join a room", func() {
		Context("When configuration is valid", func() {
			Context("And room exists", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.JoinRoom(room.Name)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound

					nonExistingRoom := "non-existing room"
					err := client.JoinRoom(nonExistingRoom)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistingRoom)))
				})
			})
		})

		Context("When the user doesn't have permission to join", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				err := client.JoinRoom(room.Name)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user doesn't have permissions to join this room"))
			})
		})

		Context("When an invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.JoinRoom(room.Name)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.JoinRoom(room.Name)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Leave a room", func() {
		Context("When configuration is valid", func() {
			Context("And room exists", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.LeaveRoom(room.Name)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			Context("And room does not exist", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound

					nonExistingRoom := "non-existing room"
					err := client.LeaveRoom(nonExistingRoom)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("room \"%s\" doesn't exist", nonExistingRoom)))
				})
			})
		})

		Context("When is not in the room", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusForbidden

				someOtherRoom := "other-room"
				err := client.LeaveRoom(someOtherRoom)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user is not in room \"" + someOtherRoom + "\"."))
			})
		})

		Context("When an invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.LeaveRoom(room.Name)

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.LeaveRoom(room.Name)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

})
