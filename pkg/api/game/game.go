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

package game

import (
	"errors"
	"fmt"
)

// Game represents a story builder game. It holds a
type Game struct {
	Turn     string   `json:"turn,omitempty"`
	Story    []Entry  `json:"story,omitempty"`
	Players  []string `json:"players,omitempty"`
	Finished bool     `json:"finished,omitempty"`

	endGame     bool
	turnCounter int
}

// Entry represents a single player's turn in the story builder game
type Entry struct {
	Text   string `json:"text"`
	Player string `json:"player"`
}

func (entry Entry) String() string {
	return fmt.Sprintf("(%s) %s", entry.Text, entry.Player)
}

// NewGame creates a game, initializing all required structures and arrays, with the provided players and initiator.
func NewGame(initiator string, players []string) *Game {
	playersCopy := make([]string, len(players))
	copy(playersCopy, players)
	for index, player := range playersCopy {
		if player == initiator { // arrange players array so that the initiator is first
			playersCopy = append([]string{initiator}, append(playersCopy[:index], playersCopy[index+1:]...)...)
			break
		}
	}
	return &Game{
		Turn:     initiator,
		Story:    make([]Entry, 0),
		Players:  playersCopy,
		Finished: false,

		endGame:     false,
		turnCounter: 1,
	}
}

// AddEntry sets the game to end after the next turn.
func (game *Game) AddEntry(entry string, issuer string) error {
	if issuer != game.Turn {
		return errors.New("invalid entry - not this player's turn")
	}
	game.Story = append(game.Story, Entry{Text: entry, Player: issuer})
	game.turnCounter++
	if game.turnCounter > len(game.Players) {
		game.turnCounter = 1
	}

	game.Turn = game.Players[game.turnCounter-1]

	if game.endGame {
		game.Finished = true
	}
	return nil
}

// EndGame sets the game to end after the next turn.
func (game *Game) EndGame() {
	game.endGame = true
}
