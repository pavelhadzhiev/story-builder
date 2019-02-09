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

package db

import (
	"errors"
	"fmt"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
)

// GetAllRooms retrieves all rooms from the server's database and returns them.
// Returns error in case of a database error.
func (sbdb *SBDatabase) GetAllRooms() ([]rooms.Room, error) {
	fmt.Println("GET ALL ROOMS")
	var roomArray = make([]rooms.Room, 2)
	roomArray[0] = rooms.Room{Name: "some room"}
	roomArray[1] = rooms.Room{Name: "some other room"}
	return roomArray, nil
}

// CreateNewRoom creates a new room in the server's database by the provided model.
// Returns error in case of a database error.
func (sbdb *SBDatabase) CreateNewRoom(room *rooms.Room) error {
	fmt.Println("CREATE NEW ROOM")
	return nil
}

// GetRoom retrieves the room with the provided name from the server's database.
// Returns error if room is not found or in case of a database error.
func (sbdb *SBDatabase) GetRoom(roomName string) (*rooms.Room, error) {
	fmt.Println("GET A ROOM")
	return &rooms.Room{Name: roomName}, nil
}

// UpdateRoom updates the room with the provided name from the server's database with the provided room model.
// Returns error if room is not found or in case of a database error.
func (sbdb *SBDatabase) UpdateRoom(roomName string, room *rooms.Room) (*rooms.Room, error) {
	fmt.Println("UPDATE A ROOM")
	return &rooms.Room{Name: roomName}, nil
}

// DeleteRoom deletes the room with the provided name from the server's database.
// Returns error if room is not found or in case of a database error.
func (sbdb *SBDatabase) DeleteRoom(roomName, issuer string) error {
	room, err := sbdb.GetRoom(roomName)
	if err != nil {
		return err
	}
	if room.Creator != issuer {
		return errors.New("user doesn't have permission to delete this room")
	}
	fmt.Println("DELETE A ROOM")
	return nil
}
