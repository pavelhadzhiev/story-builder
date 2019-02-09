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

package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// LoginHandler is an http handler for the story builder's login endpoint
func (server *SBServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/" {
		w.WriteHeader(404)
		return
	}
	switch r.Method {
	case http.MethodPost:
		authorizationHeader := r.Header.Get("Authorization")
		authorizationHeaderValue := strings.TrimPrefix(authorizationHeader, "Basic ")
		if authorizationHeader == authorizationHeaderValue {
			w.WriteHeader(400)
			w.Write([]byte("Unsupported authorization header."))
			return
		}

		credentials, err := base64.StdEncoding.DecodeString(authorizationHeaderValue)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Invalid authorization header."))
			return
		}

		split := strings.Split(string(credentials), ":")
		username, password := split[0], split[1]
		if err := server.Database.LoginUser(username, password); err != nil {
			w.WriteHeader(401)
			w.Write([]byte("Could not authenticate user."))
			return
		}

		fmt.Printf("Logged in user with name \"%s\".\n", username)
		w.Write([]byte("Successfully logged in! Welcome back, " + username + "."))
	default:
		w.WriteHeader(405)
	}
}
