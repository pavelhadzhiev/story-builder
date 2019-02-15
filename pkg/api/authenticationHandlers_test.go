package api

import (
	"encoding/base64"
	"errors"
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
	var database *dbfakes.FakeUserDatabase
	var ts *httptest.Server

	username := "username"
	password := "password"
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	BeforeEach(func() {
		// Fake a user database to be able to mock user lookups in database
		database = &dbfakes.FakeUserDatabase{}
		database.UserExistsReturns(false, nil)

		// Create the server and add the configured room to it
		sbServer = &SBServer{
			Database: database,
			Rooms:    make([]rooms.Room, 0),
			Online:   make([]string, 0),
		}
	})

	Describe("Handle register request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.RegistrationHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Context("When request is valid", func() {
			It("player should be registered and request not return error", func() {
				err := sbClient.Register()

				// Response should be 200
				Expect(err).Should(BeNil())

				// User should be online in the server
				Expect(len(sbServer.Online)).To(Equal(1))
				Expect(sbServer.Online[0]).To(Equal(username))
			})
		})

		Context("When username is taken", func() {
			It("should return error", func() {
				database.UserExistsReturns(true, nil)

				err := sbClient.Register()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("username already exists"))
			})
		})

		Context("When user lookup transaction fails", func() {
			It("should return error", func() {
				database.UserExistsReturns(false, errors.New(""))

				err := sbClient.Register()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When user registration transaction fails", func() {
			It("should return error", func() {
				database.RegisterUserReturns(errors.New(""))

				err := sbClient.Register()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid"}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.Register()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("credentials have illegal characters"))
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/register/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return HTTP 404 status cose", func() {
				resp, err := http.Get(ts.URL + "/register/zxc")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("Handle login request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.LoginHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Context("When request is valid", func() {
			It("player should be logged in and request not return error", func() {
				err := sbClient.Login()

				// Response should be 200
				Expect(err).Should(BeNil())

				// User should be online in the server
				Expect(len(sbServer.Online)).To(Equal(1))
				Expect(sbServer.Online[0]).To(Equal(username))
			})
		})

		Context("When user is already logged in", func() {
			It("should return error", func() {
				sbServer.Online = append(sbServer.Online, username)

				err := sbClient.Login()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user is already logged in"))
			})
		})

		Context("When user authentication fails", func() {
			It("should return error", func() {
				database.LoginUserReturns(errors.New(""))

				err := sbClient.Login()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user doesn't exist or password is wrong"))
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid"}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.Login()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("credentials have illegal characters"))
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/login/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return HTTP 404 status cose", func() {
				resp, err := http.Get(ts.URL + "/login/zxc")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("Handle logout request", func() {
		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(sbServer.LogoutHandler))

			clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: authHeader}
			sbClient = client.NewTestSBClient(clientConfig, ts.Client())
		})
		Context("When request is valid", func() {
			It("player should be logged out and request not return error", func() {
				// Log user in
				sbServer.Online = append(sbServer.Online, username)

				err := sbClient.Logout()

				// Response should be 200
				Expect(err).Should(BeNil())

				// User should not be online in the server
				Expect(len(sbServer.Online)).To(Equal(0))
			})
		})

		Context("When user authentication fails", func() {
			It("should return error", func() {
				database.LoginUserReturns(errors.New(""))

				err := sbClient.Logout()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user doesn't exist or password is wrong"))
			})
		})

		Context("When an invalid authorization header is provided", func() {
			It("should return error", func() {
				clientConfig = &config.SBConfiguration{URL: ts.URL, Authorization: "invalid"}
				sbClient = client.NewTestSBClient(clientConfig, ts.Client())

				err := sbClient.Logout()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("credentials have illegal characters"))
			})
		})

		Context("When request is a wrong method type", func() {
			It("should return HTTP 405 status code", func() {
				resp, err := http.Get(ts.URL + "/logout/")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})

		Context("When invalid URL is requested", func() {
			It("should return HTTP 404 status cose", func() {
				resp, err := http.Get(ts.URL + "/logout/zxc")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
