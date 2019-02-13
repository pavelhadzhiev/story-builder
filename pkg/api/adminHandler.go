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

// BanHandler is an http handler for the story builder's gameplay API
func (server *SBServer) BanHandler(w http.ResponseWriter, r *http.Request) {
	urlSuffix := strings.TrimPrefix(r.URL.Path, "/admin/ban/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")

	if len(urlSuffixSplit) == 1 || len(urlSuffixSplit) > 3 || (len(urlSuffixSplit) == 3 && urlSuffixSplit[2] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Request URL is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]
	playerToBan := urlSuffixSplit[1]

	switch r.Method {
	case http.MethodPost:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist."))
			return
		}

		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}

		if userExists, err := server.Database.UserExists(playerToBan); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Database lookup failed."))
			return
		} else if !userExists {
			w.WriteHeader(404)
			w.Write([]byte("User \"" + playerToBan + "\" doesn't exist."))
			return
		}

		if room.IsBanned(playerToBan) {
			w.WriteHeader(409)
			w.Write([]byte("User is already banned from \"" + roomName + "\". No action will be taken."))
			return
		}

		if err := room.BanPlayer(playerToBan, issuer); err != nil {
			w.WriteHeader(403)
			w.Write([]byte("You don't have admin access for room \"" + roomName + "\"."))
			return
		}

		w.Write([]byte("Player \"" + playerToBan + "\" has been banned from \"" + roomName + "\"."))
		return
	default:
		w.WriteHeader(405)
		return
	}
}
