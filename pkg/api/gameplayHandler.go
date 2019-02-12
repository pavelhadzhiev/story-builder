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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// GameplayHandler is an http handler for the story builder's gameplay API
func (server *SBServer) GameplayHandler(w http.ResponseWriter, r *http.Request) {
	urlSuffix := strings.TrimPrefix(r.URL.Path, "/gameplay/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Room name is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]

	switch r.Method {
	case http.MethodGet:
		game, err := server.GetGame(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist or no games have been started."))
			return
		}

		if responseBody, err := json.Marshal(game); err == nil {
			w.Write(responseBody)
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("Error during serialization of retrieved game."))
		return
	case http.MethodPost:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist."))
			return
		}

		if game, err := server.GetGame(roomName); err != nil || game.Finished {
			w.WriteHeader(404)
			w.Write([]byte("There is no running game."))
			return
		}

		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}

		entry := r.Header.Get("Entry-Text")
		if entry == "" {
			w.WriteHeader(400)
			w.Write([]byte("Missing Entry-Text header."))
			return
		}

		if err := room.AddEntry(entry, issuer); err != nil {
			w.WriteHeader(403)
			w.Write([]byte(fmt.Sprintf("There was an error while adding your entry: %v", err)))
			return
		}

		w.Write([]byte("Entry successfully submitted."))
	default:
		w.WriteHeader(405)
		return
	}
}
