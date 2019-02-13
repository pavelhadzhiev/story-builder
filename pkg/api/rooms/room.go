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
	Name    string   `json:"name"`
	Creator string   `json:"creator,omitempty"`
	Admins  []string `json:"admins,omitempty"`
	Banned  []string `json:"banned,omitempty"`
	Online  []string `json:"online,omitempty"`

	game         *game.Game
	previousGame *game.Game
}

// NewRoom creates a room with the provided name and creator, initializing all required structures and arrays and using the default timeout (180 seconds)
func NewRoom(name, creator string) *Room {
	admins := make([]string, 1)
	admins[0] = creator
	return &Room{
		Name:    name,
		Creator: creator,
		Admins:  admins,
		Banned:  make([]string, 0),
		Online:  make([]string, 0),

		game:         nil,
		previousGame: nil,
	}
}

func (room Room) String() string {
	return fmt.Sprintf("Name: %s\nCreator: %s\nOnline: %v\nAdmins: %v\n", room.Name, room.Creator, room.Online, room.Admins)
}

// StartGame starts a new game, including all online players and giving the provided initiator the first turn.
// Returns error if a game is already started and still ongoing or if user doesn't have admin access or is not in the room.
func (room *Room) StartGame(initiator string, timeLimit, maxLength, entriesCount int) error {
	if err := room.checkUserPermissions(initiator); err != nil {
		return err
	}

	if room.game != nil {
		if room.game.Finished {
			room.previousGame = room.game
			room.game = nil
		} else {
			return errors.New("there is an unfinished game")
		}
	}
	room.game = game.StartGame(initiator, room.Online, timeLimit, maxLength, entriesCount)
	return nil
}

// AddEntry add the provided entry text to the story on the issuers behalf.
// Returns error if there isn't a started game or it's not the issuer's turn.
func (room *Room) AddEntry(entry, issuer string) error {
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
func (room *Room) GetGame() *game.Game {
	if room.game != nil {
		return room.game
	}
	if room.previousGame != nil {
		return room.previousGame
	}
	return nil
}

// EndGame sets the currently played game to finish after the next move.
// Returns error if there isn't a started game to end or if user doesn't have admin access or is not in the room.
func (room *Room) EndGame(issuer string, entries int) error {
	if err := room.checkUserPermissions(issuer); err != nil {
		return err
	}

	if room.game == nil {
		return errors.New("there isn't a started game")
	}

	room.game.EndGame(entries)
	return nil
}

// BanPlayer bans the provided player, on behalf of the provider issuer. The banned player is instantly removed from the room and prevented from joining again.
// Returns error if the issuer doesn't have admin access.
func (room *Room) BanPlayer(playerToBan, issuer string) error {
	if err := room.checkUserPermissions(issuer); err != nil {
		return err
	}
	for index, online := range room.Online {
		if online == playerToBan {
			room.Online = append(room.Online[:index], room.Online[index+1:]...)
		}
	}
	room.Banned = append(room.Banned, playerToBan)
	return nil
}

// IsBanned returns true of the provided player has been banned from the room and false otherwise.
func (room *Room) IsBanned(player string) bool {
	for _, banned := range room.Banned {
		if banned == player {
			return true
		}
	}
	return false
}

// checkUserPermissions returns error if the user is not an admin or joined in the room.
func (room *Room) checkUserPermissions(user string) error {
	isAdmin, isOnline := false, false
	for _, admin := range room.Admins {
		if admin == user {
			isAdmin = true
		}
	}
	for _, online := range room.Online {
		if online == user {
			isOnline = true
		}
	}
	if !isAdmin {
		return errors.New("user is not an admin")
	}
	if !isOnline {
		return errors.New("user is not in the room")
	}
	return nil
}
