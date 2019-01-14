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

package config

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// SBConfigurator is an objet that can save and load Story Builder Configuration objects.
type SBConfigurator interface {
	Save(*SBConfiguration) error
	Load() (*SBConfiguration, error)
}

// SBConfiguration contains the configuration of the Story Builder CLI.
type SBConfiguration struct {
	URL           string `json:"url"`
	Authorization string `json:"authorization,omitempty"`
	Room          string `json:"room,omitempty"`
}

// ValidateConnection ensures the configuration has a valid URL for an online server.
func (sbConfig SBConfiguration) ValidateConnection() error {
	url := sbConfig.URL
	if err := validateURL(url); err != nil {
		return err
	}
	if err := healthCheckServer(url + "/healthcheck/"); err != nil {
		return err
	}
	return nil
}

// ValidateUser ensures the configuration has Authorization
func (sbConfig SBConfiguration) ValidateUser() error {
	if sbConfig.Authorization == "" {
		return errors.New("authorization property must not be empty")
	}
	return nil
}

// ValidateRoom ensures the configuration has valid URL and Authorization
func (sbConfig SBConfiguration) ValidateRoom() error {
	if sbConfig.Room == "" {
		return errors.New("room property must not be empty")
	}
	return nil
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

func healthCheckServer(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("server is not online")
	}
	return nil
}
