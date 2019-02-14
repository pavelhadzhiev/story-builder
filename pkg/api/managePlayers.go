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

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// JoinRoom puts the player in the room with the provided name.
// Returns error if a room with this name doesn't exist or the user doesn't have permission to join that room.
func (sbServer *SBServer) JoinRoom(roomName, player string) error {
	var index int
	var room rooms.Room
	roomExists := false
	for index, room = range sbServer.Rooms {
		if room.Name == roomName {
			if room.IsBanned(player) {
				return errors.New("player is not allowed to join room \"" + roomName + "\"")
			}
			roomExists = true
			break
		}
	}
	if roomExists {
		sbServer.Rooms[index].Online = append(sbServer.Rooms[index].Online, player)
		return nil
	}

	return errors.New("room with name \"" + roomName + "\" doesn't exist")
}

// LeaveRoom removes the player from the room with the provided name.
// Returns error if a room with this name doesn't exist or the user was not in it to begin with.
func (sbServer *SBServer) LeaveRoom(roomName, player string) error {
	for roomIndex, room := range sbServer.Rooms {
		if room.Name == roomName {
			for playerIndex, playerName := range sbServer.Rooms[roomIndex].Online {
				if playerName == player {
					sbServer.Rooms[roomIndex].Online = util.DeleteFromSlice(sbServer.Rooms[roomIndex].Online, playerIndex)
					return nil
				}
			}
			return errors.New("player \"" + player + "\" is not in room " + roomName + "\".")
		}
	}

	return errors.New("room with name \"" + roomName + "\" doesn't exist")
}
