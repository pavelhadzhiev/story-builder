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

// VoteHandler is an http handler for the story builder's voting API
func (server *SBServer) VoteHandler(w http.ResponseWriter, r *http.Request) {
	urlSuffix := strings.TrimPrefix(r.URL.Path, "/vote/")
	urlSuffixSplit := strings.Split(urlSuffix, "/")
	if len(urlSuffixSplit) == 1 || len(urlSuffixSplit) > 3 || (len(urlSuffixSplit) == 3 && urlSuffixSplit[2] != "") {
		w.WriteHeader(400)
		w.Write([]byte("Request URL is illegal."))
		return
	}
	roomName := urlSuffixSplit[0]
	playerToKick := urlSuffixSplit[1]

	switch r.Method {
	case http.MethodPost:
		game, err := server.GetGame(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist or no games have been started."))
			return
		}

		if game.VoteKick != nil {
			w.WriteHeader(409)
			w.Write([]byte("There is already an ongoing vote for player \"" + game.VoteKick.Player + "\"."))
			return
		}

		issuer, err := util.ExtractUsernameFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error during decoding of authorization header."))
			return
		}

		if err := game.TriggerVoteKick(issuer, playerToKick, 0.65, 60); err != nil {
			w.WriteHeader(404)
			w.Write([]byte("The requested player \"" + playerToKick + "\" to trigger a vote for is not in the game."))
			return
		}

		w.WriteHeader(202)
		w.Write([]byte("A vote to kick player \"" + playerToKick + "\" was successfully triggered."))
		return
	case http.MethodPut:
		game, err := server.GetGame(roomName)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Room \"" + roomName + "\" doesn't exist or no games have been started."))
			return
		}

		if game.VoteKick == nil {
			w.WriteHeader(404)
			w.Write([]byte("There are no ongoing votes."))
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
			if player == issuer {
				inGame = true
				break
			}
		}
		if !inGame {
			w.WriteHeader(403)
			w.Write([]byte("You cannot vote. You are not part of the game."))
			return
		}

		if err := game.Vote(issuer); err != nil {
			w.WriteHeader(409)
			w.Write([]byte("You have already voted. You can only vote once."))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("Your vote to kick player \"" + playerToKick + "\" was accepted."))
	default:
		w.WriteHeader(405)
		return
	}
}
