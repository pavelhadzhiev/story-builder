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

func (client *SBClient) call(method string, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	fullURL := client.config.URL + path

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header = *client.headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewTestSBClient(config *config.SBConfiguration, httpClient *http.Client) *SBClient {
	client := &SBClient{config: config, httpClient: httpClient}
	client.headers = &http.Header{}
	client.headers.Add("Content-Type", "application/json")
	if len(client.config.Authorization) > 0 {
		client.headers.Add("Authorization", client.config.Authorization)
	}

	return client
}
