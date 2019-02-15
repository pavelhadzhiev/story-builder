package client

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	"github.com/pavelhadzhiev/story-builder/pkg/config"
	"github.com/pavelhadzhiev/story-builder/pkg/config/configfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Story Builder Room Client test", func() {
	var client *SBClient
	var responseStatusCode int
	var responseBody []byte
	var sbServer *httptest.Server
	var clientConfig *config.SBConfiguration

	configurator := &configfakes.FakeSBConfigurator{}
	configurator.LoadReturns(clientConfig, nil)
	configurator.SaveStub = func(config *config.SBConfiguration) error {
		client.config = config
		return nil
	}

	username := "user"
	password := "password"
	authHeader := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	roomName := "someRoom"

	testHandler := func() http.HandlerFunc {
		return func(response http.ResponseWriter, req *http.Request) {
			response.WriteHeader(responseStatusCode)
			response.Write([]byte(responseBody))
		}
	}

	BeforeEach(func() {
		sbServer = httptest.NewServer(testHandler())
		clientConfig = &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader, Room: roomName}
		client = NewSBClient(clientConfig)
	})

	setupFaultyServer := func() {
		sbServer = httptest.NewUnstartedServer(testHandler())
		clientConfig := &config.SBConfiguration{URL: sbServer.URL, Authorization: authHeader, Room: roomName}
		client = NewSBClient(clientConfig)
	}

	Describe("Call healthcheck", func() {
		Context("With a valid request", func() {
			Context("And an online server", func() {
				It("should not return error", func() {
					responseStatusCode = http.StatusOK

					err := client.HealthCheck(configurator)

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			Context("And a dead server", func() {
				It("should return error", func() {
					setupFaultyServer()

					sbServer.URL = "http://valid-definitely-unexisting-url.com"
					client.config.URL = sbServer.URL
					responseStatusCode = http.StatusOK

					err := client.HealthCheck(configurator)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Configuration was wiped clean"))
					Expect(client.config.URL).To(Equal(""))
					Expect(client.config.Authorization).To(Equal(""))
					Expect(client.config.Room).To(Equal(""))
				})
			})
		})

		Context("With an invalid request", func() {
			Context("With an invalid URL", func() {
				It("should return error", func() {
					client.config.URL = "invalid-url"
					err := client.HealthCheck(configurator)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Configuration was wiped clean"))
					Expect(client.config.URL).To(Equal(""))
					Expect(client.config.Authorization).To(Equal(""))
					Expect(client.config.Room).To(Equal(""))
				})
			})

			Context("With a missing URL", func() {
				It("should return error", func() {
					client.config.URL = ""
					err := client.HealthCheck(configurator)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Configuration was wiped clean"))
					Expect(client.config.URL).To(Equal(""))
					Expect(client.config.Authorization).To(Equal(""))
					Expect(client.config.Room).To(Equal(""))
				})
			})

			Context("With an invalid authentication", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusUnauthorized

					err := client.HealthCheck(configurator)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("User and room from configuration were wiped clean"))
					Expect(client.config.URL).To(Equal(clientConfig.URL))
					Expect(client.config.Authorization).To(Equal(""))
					Expect(client.config.Room).To(Equal(""))
				})
			})

			Context("With a non-existing room", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusNotFound

					err := client.HealthCheck(configurator)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Room from configuration was wiped clean"))
					Expect(client.config.URL).To(Equal(clientConfig.URL))
					Expect(client.config.Authorization).To(Equal(clientConfig.Authorization))
					Expect(client.config.Room).To(Equal(""))
				})
			})

			Context("With a room that the player hasn't really joined", func() {
				It("should return error", func() {
					responseStatusCode = http.StatusForbidden

					err := client.HealthCheck(configurator)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Room from configuration was wiped clean"))
					Expect(client.config.URL).To(Equal(clientConfig.URL))
					Expect(client.config.Authorization).To(Equal(clientConfig.Authorization))
					Expect(client.config.Room).To(Equal(""))
				})
			})
		})

		Context("When invalid status code is returned by SB server", func() {
			It("should return error", func() {
				responseStatusCode = http.StatusCreated

				err := client.HealthCheck(configurator)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
