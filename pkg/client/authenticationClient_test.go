package client

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	"github.com/pavelhadzhiev/story-builder/pkg/config"

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

	BeforeEach(func() {
		sbServer = httptest.NewServer(testHandler)
		clientConfig := &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader}
		client = NewSBClient(clientConfig)
	})

	setupFaultyServer := func() {
		sbServer = httptest.NewUnstartedServer(testHandler)
		clientConfig := &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader}
		client = NewSBClient(clientConfig)
	}

	Describe("Register user", func() {
		Context("When request is valid", func() {
			It("should not return error", func() {
				responseStatusCode = http.StatusOK

				err := client.Register()

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("When credentials are illegal", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusBadRequest

				err := client.Register()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("credentials have illegal characters"))
			})
		})

		Context("When username is taken", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusConflict

				err := client.Register()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("username already exists"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.Register()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.Register()

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Login user", func() {
		Context("When request is valid", func() {
			It("should not return error", func() {
				responseStatusCode = http.StatusOK

				err := client.Login()

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("When credentials are illegal", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusBadRequest

				err := client.Login()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("credentials have illegal characters"))
			})
		})

		Context("When credentials are wrong", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusUnauthorized

				err := client.Login()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user doesn't exist or password is wrong"))
			})
		})

		Context("When user is already logged in", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusConflict

				err := client.Login()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user is already logged in"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.Login()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.Login()

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Logout user", func() {
		Context("When request is valid", func() {
			It("should not return error", func() {
				responseStatusCode = http.StatusOK

				err := client.Logout()

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("When credentials are illegal", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusBadRequest

				err := client.Logout()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("credentials have illegal characters"))
			})
		})

		Context("When credentials are wrong", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusUnauthorized

				err := client.Logout()

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user doesn't exist or password is wrong"))
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.Logout()

				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When there is an HTTP error", func() {
			It("should return error", func() {
				setupFaultyServer()

				responseStatusCode = http.StatusOK

				err := client.Logout()

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
