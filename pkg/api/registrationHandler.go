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

// RegistrationHandler is an http handler for the story builder's registration endpoint
func (server *SBServer) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register/" {
		w.WriteHeader(404)
		return
	}
	authorizationHeader := r.Header.Get("Authorization")
	authorizationHeaderValue := strings.TrimPrefix(authorizationHeader, "Basic ")
	if authorizationHeader == authorizationHeaderValue {
		w.Write([]byte("Unsupported authorization header."))
		w.WriteHeader(400)
		return
	}
	switch r.Method {
	case http.MethodPost:
		credentials, err := base64.StdEncoding.DecodeString(authorizationHeaderValue)
		if err != nil {
			w.Write([]byte("Invalid authorization header."))
			w.WriteHeader(400)
			return
		}

		splitted := strings.Split(string(credentials), ":")
		username, password := splitted[0], splitted[1]
		usernameTaken, err := server.Database.UserExists(username)
		if err != nil {
			w.Write([]byte("Database lookup failed."))
			w.WriteHeader(500)
			return
		}
		if usernameTaken {
			w.Write([]byte("Username is already taken."))
			w.WriteHeader(409)
			return
		}

		server.Database.RegisterUser(username, password)

		fmt.Println("Registered user with name \"", username, "\".")
		w.Write([]byte("Successfully registered! Welcome, " + username + "."))
	default:
		w.WriteHeader(405)
	}
}
