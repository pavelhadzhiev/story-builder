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
	"fmt"
	"net/http"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"
)

// TODO: check docs !!!!!!!!!!!!!!!!!!

// GetAllRooms retrieves all rooms from the server's database and returns them.
// Returns error in case of a database error.
func (client *SBClient) GetAllRooms() ([]rooms.Room, error) {
	response, err := client.call(http.MethodGet, "/rooms/", nil)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode {
	// handle status codes
	}

	defer response.Body.Close()
	var rooms []rooms.Room
	if err := json.NewDecoder(response.Body).Decode(rooms); err != nil {
		return nil, err
	}
	fmt.Println("GET ALL ROOMS")
	return rooms, nil
}

// CreateNewRoom creates a new room in the server's database by the provided model.
// Returns error in case of a database error.
func (client *SBClient) CreateNewRoom(room *rooms.Room) error {
	requestBody, err := json.Marshal(room)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.call(http.MethodPost, "/rooms/", buffer)
	if err != nil {
		return err
	}
	switch response.StatusCode {
	// handle status codes
	}

	fmt.Println("CREATE NEW ROOM")
	return nil
}

// GetRoom retrieves the room with the provided name from the server's database.
// Returns error if room is not found or in case of a database error.
func (client *SBClient) GetRoom(roomName string) (*rooms.Room, error) {
	response, err := client.call(http.MethodGet, "/rooms/"+roomName, nil)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode {
	// handle status codes
	}

	defer response.Body.Close()
	var room *rooms.Room
	if err := json.NewDecoder(response.Body).Decode(room); err != nil {
		return nil, err
	}

	fmt.Println("GET A ROOM")
	return room, nil
}

// UpdateRoom updates the room with the provided name from the server's database with the provided room model.
// Returns error if room is not found or in case of a database error.
func (client *SBClient) UpdateRoom(roomName string, room *rooms.Room) (*rooms.Room, error) {
	requestBody, err := json.Marshal(room)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.call(http.MethodPut, "/rooms/"+roomName, buffer)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode {
	// handle status codes
	}

	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(room); err != nil {
		return nil, err
	}

	fmt.Println("UPDATE A ROOM")
	return room, nil
}

// DeleteRoom deletes the room with the provided name from the server's database.
// Returns error if room is not found or in case of a database error.
func (client *SBClient) DeleteRoom(roomName string) error {
	response, err := client.call(http.MethodDelete, "/rooms/"+roomName, nil)
	if err != nil {
		return err
	}
	switch response.StatusCode {
	// handle status codes
	}

	fmt.Println("DELETE A ROOM")
	return nil
}
