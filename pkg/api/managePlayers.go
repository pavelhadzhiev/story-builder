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
	"fmt"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
)

// JoinRoom puts the player in the room with the provided name.
// Returns error if a room with this name doesn't exist or the user doesn't have permission to join that room.
func (sbServer *SBServer) JoinRoom(roomName, player string) error {
	var index int
	var room rooms.Room
	roomExists := false
	for index, room = range sbServer.Rooms {
		if room.Name == roomName {
			for _, banned := range room.Banned {
				if banned == player {
					return errors.New("player doesn't have permissions to join room \"" + roomName + "\"")
				}
			}
			roomExists = true
			fmt.Println("JOINED A ROOM")
			break
		}
	}
	if roomExists {
		sbServer.Rooms[index].Players = append(sbServer.Rooms[index].Players, player)
		return nil
	}

	fmt.Println("ROOM NOT FOUND")
	return errors.New("room with name \"" + roomName + "\" doesn't exist")
}
