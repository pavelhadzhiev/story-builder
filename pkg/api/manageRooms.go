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
)

// GetAllRooms retrieves all rooms from the server and returns them.
func (sbServer *SBServer) GetAllRooms() []rooms.Room {
	return sbServer.Rooms
}

// CreateNewRoom creates a new room in the server, using the provided model.
// Returns error if a room with this name already exists.
func (sbServer *SBServer) CreateNewRoom(room *rooms.Room) error {
	if _, err := sbServer.GetRoom(room.Name); err == nil {
		return errors.New("a room with this name already exists")
	}
	sbServer.Rooms = append(sbServer.Rooms, *room)
	return nil
}

// GetRoom retrieves the room with the provided name from the server.
// Returns error if a room with this name doesn't exist.
func (sbServer *SBServer) GetRoom(roomName string) (*rooms.Room, error) {
	for _, room := range sbServer.Rooms {
		if room.Name == roomName {
			return &room, nil
		}
	}
	return nil, errors.New("room with name \"" + roomName + "\" doesn't exist")
}

// DeleteRoom deletes the room with the provided name from the server.
// Returns error if a room with this name doesn't exist or the issuer doesn't have the permissions to delete it.
func (sbServer *SBServer) DeleteRoom(roomName, issuer string) error {
	var index int
	var room rooms.Room
	roomExists := false
	for index, room = range sbServer.Rooms {
		if room.Name == roomName {
			roomExists = true
			break
		}
	}
	if !roomExists {
		return errors.New("room with name \"" + roomName + "\" doesn't exist")
	}
	if room.Creator != issuer {
		return errors.New("user doesn't have permission to delete this room")
	}

	sbServer.Rooms = append(sbServer.Rooms[:index], sbServer.Rooms[index+1:]...)
	return nil
}
