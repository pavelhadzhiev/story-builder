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
	"fmt"
	"net/http"
	"strings"

	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// JoinRoomHandler is an http handler for the story builder's join room API
func (server *SBServer) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("JOIN ROOM REQUEST:", r)
	if !strings.HasPrefix(r.URL.Path, "/join-room/") {
		w.WriteHeader(404)
		return
	}

	urlSuffix := strings.TrimPrefix(r.URL.Path, "/join-room/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Room name is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]

	switch r.Method {
	case http.MethodPost:
		player, err := util.DecodeBasicAuthorization(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}
		if _, err := server.GetRoom(roomName); err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" not found."))
			return
		}
		if err := server.JoinRoom(roomName, player); err != nil {
			w.WriteHeader(403)
			w.Write([]byte("The user doesn't have permissions to join that room."))
			return
		}
		w.Write([]byte("Joined room \"" + roomName + "\" successfully."))
		return
	default:
		w.WriteHeader(405)
		return
	}
}
