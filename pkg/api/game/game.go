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
	"strings"
	"time"
)

// Game represents a story builder game. It holds a
type Game struct {
	Turn     string   `json:"turn,omitempty"`
	Story    []Entry  `json:"story,omitempty"`
	Players  []string `json:"players,omitempty"`
	Finished bool     `json:"finished,omitempty"`
	TimeLeft int      `json:"timeleft,omitempty"`

	endGame    bool
	playerTurn int
	timeLimit  int
}

// Entry represents a single player's turn in the story builder game
type Entry struct {
	Text   string `json:"text"`
	Player string `json:"player"`
}

func (game Game) String() string {
	gameString := "Players in the current game: "
	for _, player := range game.Players {
		gameString += player + ","
	}
	gameString = strings.TrimSuffix(gameString, ",")

	gameString += "\n--------------------------------\n"
	for _, entry := range game.Story {
		gameString += entry.String() + "\n"
	}
	gameString += "--------------------------------\n"

	if game.Finished {
		gameString += "The game has finished. You can now start the next one!\n"
	} else {
		gameString += fmt.Sprintf("Next turn: Player \"%s\"\n", game.Turn)
		gameString += fmt.Sprintf("Time left: %d\n", game.TimeLeft)
	}

	return gameString
}

func (entry Entry) String() string {
	return fmt.Sprintf("%s (By \"%s\")", entry.Text, entry.Player)
}

// StartGame creates a game, initializing all required structures and arrays, with the provided players and initiator.
func StartGame(initiator string, players []string, timeLimit int) *Game {
	playersCopy := make([]string, len(players))
	copy(playersCopy, players)
	for index, player := range playersCopy {
		if player == initiator { // arrange players array so that the initiator is first
			playersCopy = append([]string{initiator}, append(playersCopy[:index], playersCopy[index+1:]...)...)
			break
		}
	}

	game := &Game{
		Turn:     initiator,
		Story:    make([]Entry, 0),
		Players:  playersCopy,
		Finished: false,
		TimeLeft: timeLimit,

		endGame:    false,
		playerTurn: 1,
		timeLimit:  timeLimit,
	}

	go game.monitorTime()

	return game
}

// AddEntry sets the game to end after the next turn.
func (game *Game) AddEntry(entry string, issuer string) error {
	if issuer != game.Turn {
		return errors.New("invalid entry - not this player's turn")
	}

	game.Story = append(game.Story, Entry{Text: entry, Player: issuer})
	game.setNextTurn()

	if game.endGame {
		game.Finished = true
	}
	return nil
}

// EndGame sets the game to end after the next turn.
func (game *Game) EndGame() {
	game.endGame = true
}

func (game *Game) monitorTime() {
	for !game.Finished {
		game.TimeLeft--
		if game.TimeLeft <= 0 {
			game.setNextTurn()
		}
		time.Sleep(1 * time.Second)
	}
}

func (game *Game) setNextTurn() {
	game.playerTurn++
	if game.playerTurn > len(game.Players) {
		game.playerTurn = 1
	}

	game.Turn = game.Players[game.playerTurn-1]
	game.TimeLeft = game.timeLimit
}
