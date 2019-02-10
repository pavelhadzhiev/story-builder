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

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// RoomHandler is an http handler for the story builder's room API
func (server *SBServer) RoomHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ROOM REQUEST:", r)
	if !strings.HasPrefix(r.URL.Path, "/rooms/") {
		w.WriteHeader(404)
		return
	}

	if r.URL.Path == "/rooms/" {
		switch r.Method {
		case http.MethodGet:
			rooms, err := server.GetAllRooms()
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error while retrieving rooms from database."))
				fmt.Println("Error while retrieving room from database.")
				return
			}
			if rooms != nil {
				responseBody, err := json.Marshal(rooms)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte("Error during serialization of retrieved room."))
					fmt.Println("Error during serialization of retrieved room.")
					return
				}
				fmt.Println("Rooms found and returned as body")
				w.Write(responseBody)
				return
			}
			w.WriteHeader(404) // REMOVE
			return
		case http.MethodPost:
			server.CreateNewRoom(&rooms.Room{})
			return
		default:
			w.WriteHeader(405)
			return
		}
	}

	var urlSuffix = strings.TrimPrefix(r.URL.Path, "/rooms/")
	var urlSuffixSplit = strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Room name is illegal."))
		return
	}
	var roomName = urlSuffixSplit[0]
	fmt.Println("ROOM NAME:", roomName)
	switch r.Method {
	case http.MethodGet:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error while retrieving room from database."))
			fmt.Println("Error while retrieving room from database.")
			return
		}
		if room != nil {
			responseBody, err := json.Marshal(room)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error during serialization of retrieved room."))
				fmt.Println("Error during serialization of retrieved room.")
				return
			}
			fmt.Println("Room found and returned as body")
			w.Write(responseBody)
			return
		}
		w.WriteHeader(404)
		fmt.Println("404 ROOM NOT FOUND")
		return
	case http.MethodPut:
		room, err := server.UpdateRoom(roomName, &rooms.Room{})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error while updating room in database."))
			return
		}
		if room != nil {
			responseBody, err := json.Marshal(room)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Error during serialization of updated room."))
				return
			}
			w.Write(responseBody)
			return
		}
		w.WriteHeader(404)
		return
	case http.MethodDelete:
		issuer, err := util.DecodeBasicAuthorization(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}
		server.DeleteRoom(roomName, issuer)
		return
	default:
		w.WriteHeader(405)
		return
	}
}
