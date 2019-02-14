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
	"net/http"
	"strings"

	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// HealthcheckHandler is an http handler for the story builder's healthcheck endpoint
func (server *SBServer) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {

			// Validate authentication
			user, pass, err := util.ExtractCredentialsFromAuthorizationHeader(authHeader)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error during decoding of authorization header."))
				return
			}
			if err := server.Database.LoginUser(user, pass); err != nil {
				w.WriteHeader(401)
				w.Write([]byte("Authentication for user \"" + user + "\" failed."))
				return
			}

			// Validate room
			urlSuffix := strings.TrimPrefix(r.URL.Path, "/healthcheck/")
			if urlSuffix != "" {
				urlSuffixSplit := strings.Split(urlSuffix, "/")
				if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
					w.WriteHeader(400)
					w.Write([]byte("Room name is illegal."))
					return
				}
				roomName := urlSuffixSplit[0]
				room, err := server.GetRoom(roomName)
				if err != nil {
					w.WriteHeader(404)
					w.Write([]byte("Room \"" + roomName + "\" not found."))
					return
				}
				if !room.IsOnline(user) {
					w.WriteHeader(403)
					w.Write([]byte("Room \"" + roomName + "\" exists but player \"" + user + "\" is not int it."))
					return
				}
			}
		}

		// Returns status code 200 to show server is online and healthy and all configurations are valid
		w.WriteHeader(200)
	default:
		w.WriteHeader(405)
	}
}
