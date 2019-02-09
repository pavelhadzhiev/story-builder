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
)

// RoomHandler is an http handler for the story builder's room API
func (server *SBServer) RoomHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/rooms/") {
		w.WriteHeader(404)
		return
	}

	if r.URL.Path == "/rooms/" {
		switch r.Method {
		case http.MethodGet:
			//GetAllRooms()
			fmt.Println("GET ALL ROOMS")
		case http.MethodPost:
			//CreateNewRoom(room)
			fmt.Println("CREATE NEW ROOM")
		default:
			w.WriteHeader(405)
			return
		}
	}

	var urlSuffix = strings.TrimPrefix(r.URL.Path, "/rooms/")
	var urlSuffixSplit = strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) > 2 || (len(urlSuffixSplit) == 2 && urlSuffixSplit[1] != "") {
		fmt.Println("INVALID ROOM NAME")
		w.WriteHeader(400)
		w.Write([]byte("Room name is illegal."))
		return
	}
	var roomName = urlSuffixSplit[0]
	fmt.Println("ROOM NAME:", roomName)
	switch r.Method {
	case http.MethodGet:
		//GetRoom(roomName)
		fmt.Println("GET A ROOM")
	case http.MethodPut:
		//UpdateRoom(roomName, room)
		fmt.Println("UPDATE A ROOM")
	case http.MethodDelete:
		//DeleteRoom(roomName)
		fmt.Println("DELETE A ROOM")
	default:
		w.WriteHeader(405)
		return
	}

	w.Write([]byte("Called the rooms API!"))
}
