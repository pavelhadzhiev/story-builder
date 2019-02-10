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

package rooms

import "fmt"

// Room represents a story builder room, in which a select group of players can play the game
type Room struct {
	Name    string    `json:"name"`
	Creator string    `json:"creator,omitempty"`
	Rules   RoomRules `json:"rules,omitempty"`
	Admins  []string  `json:"admins,omitempty"`
	Banned  []string  `json:"banned,omitempty"`

	Turn    string   `json:"turn,omitempty"`
	Players []string `json:"players,omitempty"`
	Story   []Entry  `json:"story,omitempty"`
}

func (room Room) String() string {
	return fmt.Sprintf("Name: %s\nCreator: %s\nPlayers:%v\n", room.Name, room.Creator, room.Players)
}

// NewRoom creates a room with the provided name and creator, initializing all required structures and arrays and using the default timeout (180 seconds)
func NewRoom(name, creator string) *Room {
	admins := make([]string, 1)
	admins[0] = creator
	return &Room{
		Name:    name,
		Creator: creator,
		Rules:   RoomRules{Timeout: 180},
		Admins:  admins,
		Banned:  make([]string, 0),

		Turn:    creator,
		Players: make([]string, 0),
		Story:   make([]Entry, 0),
	}
}

// Entry represents a single player's turn in the story builder game
type Entry struct {
	Text   string `json:"text"`
	Player string `json:"player"`
}

// RoomRules keeps some configurations for the gameplay in the story builder room.
type RoomRules struct {
	Timeout int `json:"timeout,omitempty"`
}
