// Copyright © 2019 Pavel Hadzhiev <p.hadzhiev96@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pavelhadzhiev/story-builder/pkg/config"
)

// HealthCheck validates the configured URL and the health of the connection.
// Also validates authentication and room if provided.
func (client *SBClient) HealthCheck(configurator config.SBConfigurator) error {
	defer configurator.Save(client.config)
	if err := validateURL(client.config.URL); err != nil {
		client.wipeConnection()
		return fmt.Errorf("%v. Configuration was wiped clean. Use the connect command to connect to a server", err)
	}
	response, err := client.call(http.MethodPost, "/healthcheck/"+client.config.Room, nil, nil)
	if err != nil {
		client.wipeConnection()
		return errors.New("server is not valid or unhealthy: Configuration was wiped clean. Use the connect command to connect to a server")
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 401:
		defer client.wipeAuthorization()
		return errors.New("authentication failed. User and room from configuration were wiped clean")
	case 403:
		defer client.wipeRoom()
		return errors.New("player is not in room \"" + client.config.Room + "\". Room from configuration was wiped clean")
	case 404:
		defer client.wipeRoom()
		return errors.New("room \"" + client.config.Room + "\" doesn't exist. Room from configuration was wiped clean")
	default:
		return errors.New("something went really wrong :(")
	}
}

func validateURL(URL string) error {
	if URL == "" {
		return errors.New("missing URL")
	}

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("unparsable URL: %s", err)
	}

	if !parsedURL.IsAbs() || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return fmt.Errorf("non-HTTP URL: %s", URL)
	}

	return nil
}

func (client *SBClient) wipeConnection() {
	client.config.URL = ""
	client.wipeAuthorization()
}

func (client *SBClient) wipeAuthorization() {
	client.config.Authorization = ""
	client.wipeRoom()
}

func (client *SBClient) wipeRoom() {
	client.config.Room = ""
}
