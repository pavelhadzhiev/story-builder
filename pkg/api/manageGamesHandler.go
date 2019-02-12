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
	"strconv"
	"strings"

	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// ManageGamesHandler is an http handler for the story builder's game management API
func (server *SBServer) ManageGamesHandler(w http.ResponseWriter, r *http.Request) {
	urlSuffix := strings.TrimPrefix(r.URL.Path, "/manage-games/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Room name is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]

	switch r.Method {
	case http.MethodPost:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist."))
			return
		}
		if game, err := server.GetGame(roomName); err == nil && !game.Finished {
			w.WriteHeader(409)
			w.Write([]byte("There is already a running game."))
		}

		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}

		var timeLimit int
		timeLimitString := r.Header.Get("Time-Limit")
		if timeLimitString != "" {
			timeLimit, err = strconv.Atoi(timeLimitString)
			if err != nil || timeLimit < 0 {
				w.WriteHeader(400)
				w.Write([]byte("Illegal Time-Limit header value."))
				return
			}
		} else {
			timeLimit = 60 // Set default time limit, if one is not provided
		}

		var maxLength int
		maxLengthString := r.Header.Get("Max-Length")
		if maxLengthString != "" {
			maxLength, err = strconv.Atoi(maxLengthString)
			if err != nil || maxLength < 0 {
				w.WriteHeader(400)
				w.Write([]byte("Illegal Max-Length header value."))
				return
			}
		} else {
			maxLength = 100 // Set default max length, if one is not provided
		}

		if err := room.StartGame(issuer, timeLimit, maxLength); err != nil {
			w.WriteHeader(403)
			w.Write([]byte("Game cannot be started. Requires user to be joined and have admin access."))
			return
		}

		w.Write([]byte("Game successfully started in room \"" + roomName + "\"."))
	case http.MethodDelete:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist."))
			return
		}
		if game, err := server.GetGame(roomName); err != nil || game.Finished {
			w.WriteHeader(409)
			w.Write([]byte("There is no running game."))
		}

		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}

		if err := room.EndGame(issuer); err != nil {
			w.WriteHeader(403)
			w.Write([]byte("Game cannot be ended. Requires user to be joined and have admin access."))
			return
		}

		w.WriteHeader(202)
		w.Write([]byte("Game end successfully triggered in room \"" + roomName + "\". Next move will be the last."))
	default:
		w.WriteHeader(405)
		return
	}
}
