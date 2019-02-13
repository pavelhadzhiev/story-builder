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

// BanHandler is an http handler for the story builder's admin API
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
	case http.MethodDelete:
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

		room.GetGame().Kick(playerToBan)
		w.Write([]byte("Player \"" + playerToBan + "\" has been banned from room \"" + roomName + "\"."))
		return
	default:
		w.WriteHeader(405)
		return
	}
}

// KickHandler is an http handler for the story builder's admin API
func (server *SBServer) KickHandler(w http.ResponseWriter, r *http.Request) {
	urlSuffix := strings.TrimPrefix(r.URL.Path, "/admin/kick/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")

	if len(urlSuffixSplit) == 1 || len(urlSuffixSplit) > 3 || (len(urlSuffixSplit) == 3 && urlSuffixSplit[2] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Request URL is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]
	playerToKick := urlSuffixSplit[1]

	switch r.Method {
	case http.MethodDelete:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist."))
			return
		}

		game, err := server.GetGame(roomName)
		if err != nil || game.Finished {
			w.WriteHeader(404)
			w.Write([]byte("There is no runnig game in room \"" + roomName + "\"."))
			return
		}

		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}

		inGame := false
		for _, player := range game.Players {
			if player == playerToKick {
				inGame = true
				break
			}
		}
		if !inGame {
			w.WriteHeader(404)
			w.Write([]byte("The user to kick is not in the game."))
			return
		}

		isAdmin := false
		for _, admin := range room.Admins {
			if admin == issuer {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			w.WriteHeader(403)
			w.Write([]byte("You don't have admin access for room \"" + roomName + "\"."))
			return
		}

		game.Kick(playerToKick)
		w.Write([]byte("Player \"" + playerToKick + "\" has been kicked from the game in room \"" + roomName + "\"."))
		return
	default:
		w.WriteHeader(405)
		return
	}
}

// PromoteAdminHandler is an http handler for the story builder's admin API
func (server *SBServer) PromoteAdminHandler(w http.ResponseWriter, r *http.Request) {
	urlSuffix := strings.TrimPrefix(r.URL.Path, "/admin/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")

	if len(urlSuffixSplit) == 1 || len(urlSuffixSplit) > 3 || (len(urlSuffixSplit) == 3 && urlSuffixSplit[2] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Request URL is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]
	userToPromote := urlSuffixSplit[1]

	switch r.Method {
	case http.MethodPost:
		room, err := server.GetRoom(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist."))
			return
		}

		if userExists, err := server.Database.UserExists(userToPromote); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Database lookup failed."))
			return
		} else if !userExists {
			w.WriteHeader(404)
			w.Write([]byte("User \"" + userToPromote + "\" doesn't exist."))
			return
		}

		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}
		if err := room.PromoteAdmin(userToPromote, issuer); err != nil {
			w.WriteHeader(403)
			w.Write([]byte("You don't have admin access for room \"" + roomName + "\"."))
			return
		}

		w.Write([]byte("User \"" + userToPromote + "\" has been promoted to admin in room \"" + roomName + "\"."))
	default:
		w.WriteHeader(405)
		return
	}
}
