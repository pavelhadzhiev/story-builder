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

import (
	"errors"
	"fmt"

	"github.com/pavelhadzhiev/story-builder/pkg/api/game"
)

// Room represents a story builder room, in which a select group of players can play the game
type Room struct {
	Name    string    `json:"name"`
	Creator string    `json:"creator,omitempty"`
	Rules   RoomRules `json:"rules,omitempty"`
	Admins  []string  `json:"admins,omitempty"`
	Banned  []string  `json:"banned,omitempty"`
	Online  []string  `json:"online,omitempty"`

	game         *game.Game
	previousGame *game.Game
}

// StartGame starts a new game, including all online players and giving the provided initiator the first turn.
// Returns error if a game is already started and still ongoing.
func (room *Room) StartGame(initiator string) error {
	isAdmin, isOnline := false, false
	for _, admin := range room.Admins {
		if admin == initiator {
			isAdmin = true
		}
	}
	for _, online := range room.Online {
		if online == initiator {
			isOnline = true
		}
	}
	if !isAdmin {
		return errors.New("game initiator is not an admin")
	}
	if !isOnline {
		return errors.New("game initiator is not in the room")
	}

	if room.game != nil {
		return errors.New("cannot start a new game until last one is finished")
	}
	room.game = game.NewGame(initiator, room.Online)
	return nil
}

// EndGame sets the currently played game to finish after the next move.
// Returns error if there isn't a started game to end.
func (room Room) EndGame() error {
	if room.game == nil {
		return errors.New("there isn't a started game")
	}
	room.game.EndGame()
	return nil
}

// AddEntry add the provided entry text to the story on the issuers behalf.
// Returns error if there isn't a started game or it's not the issuer's turn.
func (room Room) AddEntry(entry, issuer string) error {
	if room.game == nil {
		return errors.New("there isn't a started game")
	}

	if err := room.game.AddEntry(entry, issuer); err != nil {
		return err
	}

	if room.game.Finished {
		room.previousGame = room.game
		room.game = nil
	}
	return nil
}

// GetGame returns the current or last finished game.
func (room Room) GetGame() *game.Game {
	if room.game != nil {
		return room.game
	}
	if room.previousGame != nil {
		return room.previousGame
	}
	return nil
}

func (room Room) String() string {
	return fmt.Sprintf("Name: %s\nCreator: %s\nOnline:%v\n", room.Name, room.Creator, room.Online)
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
		Online:  make([]string, 0),

		game:         nil,
		previousGame: nil,
	}
}

// RoomRules keeps some configurations for the gameplay in the story builder room.
type RoomRules struct {
	Timeout int `json:"timeout,omitempty"`
}
