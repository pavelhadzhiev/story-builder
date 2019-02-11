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
	"errors"

	"github.com/pavelhadzhiev/story-builder/pkg/api/game"
)

// GetGame returns the currently played game in the room with the provided name.
// If there isn't one, it returns the last finished game.
// Returns an error if no game has ever been played or such a room doesn't exist.
func (sbServer *SBServer) GetGame(roomName string) (*game.Game, error) {
	for index, room := range sbServer.Rooms {
		if room.Name == roomName {
			if game := sbServer.Rooms[index].GetGame(); game != nil {
				return game, nil
			}
			return nil, errors.New("there hasn't been a started game in room \"" + roomName + "\" yet")
		}
	}

	return nil, errors.New("room \"" + roomName + "\" doesn't exist")
}
