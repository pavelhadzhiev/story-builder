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

package util

import (
	"encoding/base64"
	"errors"
	"strings"
)

// ExtractUsernameFromAuthorizationHeader takes a Basic Authorization header and decodes it to return the username
func ExtractUsernameFromAuthorizationHeader(authorizationHeader string) (string, error) {
	authorizationHeaderValue := strings.TrimPrefix(authorizationHeader, "Basic ")
	if authorizationHeader == authorizationHeaderValue {
		return "", errors.New("unsupported authorization header type")
	}

	credentials, err := base64.StdEncoding.DecodeString(authorizationHeaderValue)
	if err != nil {
		return "", errors.New("invalid authorization header")
	}

	split := strings.Split(string(credentials), ":")
	username := split[0]

	return username, nil
}

// ExtractCredentialsFromAuthorizationHeader takes a Basic Authorization header and decodes it to return the username and passowrd
func ExtractCredentialsFromAuthorizationHeader(authorizationHeader string) (string, string, error) {
	authorizationHeaderValue := strings.TrimPrefix(authorizationHeader, "Basic ")
	if authorizationHeader == authorizationHeaderValue {
		return "", "", errors.New("unsupported authorization header type")
	}

	credentials, err := base64.StdEncoding.DecodeString(authorizationHeaderValue)
	if err != nil {
		return "", "", errors.New("invalid authorization header")
	}

	split := strings.Split(string(credentials), ":")
	username, password := split[0], split[1]
	return username, password, nil
}
