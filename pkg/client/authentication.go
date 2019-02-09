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
	"net/http"
)

// Register makes a request to the story builder server to register the user in the configuration
func (client *SBClient) Register() (*http.Response, error) {
	response, err := client.call(http.MethodPost, "/register/", nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return response, nil
	case 400:
		return response, errors.New("credentials have illegal characters")
	case 409:
		return response, errors.New("username already exists")
	default:
		return response, errors.New("something went really wrong :(")
	}
}

// Login makes a request to the story builder server to check whether the user in the configuration is registered in the server DB.
func (client *SBClient) Login() (*http.Response, error) {
	response, err := client.call(http.MethodPost, "/login/", nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return response, nil
	case 400:
		return response, errors.New("credentials have illegal characters")
	case 401:
		return response, errors.New("user doesn't exist or password is wrong")
	default:
		return response, errors.New("something went really wrong :(")
	}
}
