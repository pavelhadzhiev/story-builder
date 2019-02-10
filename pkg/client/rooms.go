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

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
)

// GetAllRooms retrieves all rooms from the server and returns them.
func (client *SBClient) GetAllRooms() ([]rooms.Room, error) {
	response, err := client.call(http.MethodGet, "/rooms/", nil)
	if err != nil {
		return nil, fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		defer response.Body.Close()
		var rooms = make([]rooms.Room, 0, 100)
		if err := json.NewDecoder(response.Body).Decode(&rooms); err != nil {
			return nil, fmt.Errorf("failed to deserialize response from server: %e", err)
		}
		return rooms, nil
	default:
		return nil, errors.New("something went really wrong :(")
	}
}

// CreateNewRoom creates a new room in the server, using the provided model.
// Returns error if a room with this name already exists.
func (client *SBClient) CreateNewRoom(room *rooms.Room) error {
	requestBody, err := json.Marshal(room)
	if err != nil {
		return fmt.Errorf("failed to serialize room: %e", err)
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.call(http.MethodPost, "/rooms/", buffer)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 204:
		return nil
	case 403:
		return errors.New("room \"" + room.Name + "\" already exists")
	default:
		return errors.New("something went really wrong :(")
	}
}

// GetRoom retrieves the room with the provided name from the server.
// Returns error if a room with this name doesn't exist.
func (client *SBClient) GetRoom(roomName string) (*rooms.Room, error) {
	response, err := client.call(http.MethodGet, "/rooms/"+roomName, nil)
	if err != nil {
		return nil, fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		defer response.Body.Close()
		var room = &rooms.Room{}
		if err := json.NewDecoder(response.Body).Decode(room); err != nil {
			return nil, err
		}
		return room, nil
	case 404:
		return nil, errors.New("room \"" + roomName + "\" doesn't exist")
	default:
		return nil, errors.New("something went really wrong :(")
	}
}

// DeleteRoom deletes the room with the provided name from the server.
// Returns error if a room with this name doesn't exist or the issuer doesn't have the permissions to delete it.
func (client *SBClient) DeleteRoom(roomName string) error {
	response, err := client.call(http.MethodDelete, "/rooms/"+roomName, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 204:
		return nil
	case 403:
		return errors.New("user doesn't have permissions to delete this room")
	case 404:
		return errors.New("room \"" + roomName + "\" doesn't exist")
	default:
		return errors.New("something went really wrong :(")
	}
}

// JoinRoom puts the player in the room with the provided name.
// Returns error if a room with this name doesn't exist or the user doesn't have permission to join that room.
func (client *SBClient) JoinRoom(roomName string) error {
	response, err := client.call(http.MethodPost, "/join-room/"+roomName, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		return errors.New("user doesn't have permissions to join this room")
	case 404:
		return errors.New("room \"" + roomName + "\" doesn't exist")
	default:
		return errors.New("something went really wrong :(")
	}
}

// LeaveRoom removes the player from the room with the provided name.
func (client *SBClient) LeaveRoom(roomName string) error {
	response, err := client.call(http.MethodPost, "/leave-room/"+roomName, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		return errors.New("user is not in room \"" + roomName + "\".")
	case 404:
		return errors.New("room \"" + roomName + "\" doesn't exist")
	default:
		return errors.New("something went really wrong :(")
	}
}
