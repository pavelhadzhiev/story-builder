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
	"net/http"
	"strings"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// RoomHandler is an http handler for the story builder's room API
func (server *SBServer) RoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/rooms/" {
		switch r.Method {
		case http.MethodGet:
			rooms := server.GetAllRooms()
			responseBody, err := json.Marshal(rooms)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error during serialization of retrieved rooms."))
				return
			}
			w.Write(responseBody)
			return
		case http.MethodPost:
			var room = &rooms.Room{}
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(room); err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error during serialization of retrieved rooms."))
				return
			}

			if err := server.CreateNewRoom(room); err != nil {
				w.WriteHeader(409)
				w.Write([]byte("Cannot create more room. A room with this name already exists"))
				return
			}

			w.WriteHeader(201)
			return
		default:
			w.WriteHeader(405)
			return
		}
	}

	urlSuffix := strings.TrimPrefix(r.URL.Path, "/rooms/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Room name is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]
	switch r.Method {
	case http.MethodGet:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" not found."))
			return
		}
		responseBody, err := json.Marshal(room)
		if err == nil {
			w.Write(responseBody)
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("Error during serialization of retrieved room."))
		return
	case http.MethodDelete:
		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
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
		if err := server.DeleteRoom(roomName, issuer); err != nil {
			w.WriteHeader(403)
			w.Write([]byte("You are not authorized to delete this room."))
			return
		}
		w.WriteHeader(204)
		return
	default:
		w.WriteHeader(405)
		return
	}
}

// JoinRoomHandler is an http handler for the story builder's join room API
func (server *SBServer) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
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
		player, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
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

// LeaveRoomHandler is an http handler for the story builder's leave room API
func (server *SBServer) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	urlSuffix := strings.TrimPrefix(r.URL.Path, "/leave-room/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Room name is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]

	switch r.Method {
	case http.MethodPost:
		player, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
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
		if err := server.LeaveRoom(roomName, player); err != nil {
			w.WriteHeader(403)
			w.Write([]byte("User \"" + player + "\" is not in room \"" + roomName + "\"."))
			return
		}
		w.Write([]byte("Left room \"" + roomName + "\" successfully."))
		return
	default:
		w.WriteHeader(405)
		return
	}
}
