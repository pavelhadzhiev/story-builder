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
	"smctl/pkg/errors"
	"smctl/pkg/httputil"

	"github.com/pavelhadzhiev/story-builder/pkg/config"
)

// StoryBuilderClient is an HTTP client wrapper that executes story builder API requests.
type StoryBuilderClient struct {
	config     *config.SBConfiguration
	httpClient *http.Client
	headers    *http.Header
}

// NewStoryBuilderClient creates a new StoryBuilderClient with a given SBConfiguration.
// It attaches an application/json content-type header and authorization, if set in the configuration.
func NewStoryBuilderClient(config *config.SBConfiguration) *StoryBuilderClient {
	client := &StoryBuilderClient{config: config, httpClient: &http.Client{}}
	client.headers = &http.Header{}
	client.headers.Add("Content-Type", "application/json")
	if len(client.config.Authorization) > 0 {
		client.headers.Add("Authorization", client.config.Authorization)
	}

	return client
}

// Register makes a request to the story builder server to register the user in the configuration
func (client *StoryBuilderClient) Register() error {
	if _, err := client.call(http.MethodPost, "/register/", nil); err != nil {
		return err
	}
	return nil
}

// Login makes a request to the story builder server to check whether the user in the configuration is registered in the server DB.
func (client *StoryBuilderClient) Login() error {
	if _, err := client.call(http.MethodPost, "/login/", nil); err != nil {
		return err
	}
	return nil
}

func (client *StoryBuilderClient) call(method string, path string, body io.Reader) (*http.Response, error) {
	URL := httputil.NormalizeURL(client.config.URL)
	fullURL := URL + path

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
		respErr := errors.ResponseError{
			URL:        fullURL,
			StatusCode: resp.StatusCode,
		}

		return nil, respErr
	}

	return resp, nil
}
