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
	"io"
	"net/http"

	"github.com/pavelhadzhiev/story-builder/pkg/config"
)

// SBClient is an HTTP client wrapper that executes story builder API requests.
type SBClient struct {
	config     *config.SBConfiguration
	httpClient *http.Client
	headers    *http.Header
}

// NewSBClient creates a new StoryBuilderClient with a given SBConfiguration.
// It attaches an application/json content-type header and authorization, if set in the configuration.
func NewSBClient(config *config.SBConfiguration) *SBClient {
	client := &SBClient{config: config, httpClient: &http.Client{}}
	client.headers = &http.Header{}
	client.headers.Add("Content-Type", "application/json")
	if len(client.config.Authorization) > 0 {
		client.headers.Add("Authorization", client.config.Authorization)
	}

	return client
}

// Register makes a request to the story builder server to register the user in the configuration
func (client *SBClient) Register() (*http.Response, error) {
	resp, err := client.call(http.MethodPost, "/register/", nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Login makes a request to the story builder server to check whether the user in the configuration is registered in the server DB.
func (client *SBClient) Login() (*http.Response, error) {
	resp, err := client.call(http.MethodPost, "/login/", nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *SBClient) call(method string, path string, body io.Reader) (*http.Response, error) {
	fullURL := client.config.URL + path

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header = *client.headers

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respErr := ResponseError{
			URL:        fullURL,
			StatusCode: resp.StatusCode,
		}

		return nil, respErr
	}

	return resp, nil
}
