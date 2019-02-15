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
	"math"
	"strings"
	"time"

	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// Game represents a story builder game. It holds a
type Game struct {
	Turn        string    `json:"turn,omitempty"`
	Story       []Entry   `json:"story,omitempty"`
	Players     []string  `json:"players,omitempty"`
	Finished    bool      `json:"finished,omitempty"`
	TimeLeft    int       `json:"timeLeft,omitempty"`
	MaxLength   int       `json:"maxLength,omitempty"`
	MaxEntries  int       `json:"maxEntries,omitempty"`
	EntriesLeft int       `json:"entriesLeft,omitempty"`
	VoteKick    *VoteKick `json:"votekick,omiempty"`

	playerTurn int
	timeLimit  int
}

func (game *Game) String() string {
	gameString := "\n"
	if game.VoteKick != nil {
		gameString += "ATTENTION: There is a kick vote going on!\n"
		gameString += game.VoteKick.String()
		gameString += "\nYou can cast your vote using the vote command.\n\n"
	}

	gameString += "Players in the game: "
	for _, player := range game.Players {
		gameString += player + ", "
	}
	gameString = strings.TrimSuffix(gameString, ", ")

	gameString += "\n--------------------------------\n"
	for _, entry := range game.Story {
		gameString += entry.String() + "\n"
	}
	gameString += "--------------------------------\n"

	if game.Finished {
		gameString += "The game has finished. You can now start the next one!\n"
	} else {
		gameString += fmt.Sprintf("Next turn: Player \"%s\"\n", game.Turn)
		if game.MaxLength != 0 {
			gameString += fmt.Sprintf("Max length: %d symbols\n", game.MaxLength)
		}
		if game.TimeLeft != 0 {
			gameString += fmt.Sprintf("Time left: %d seconds\n", game.TimeLeft)
		}
		if game.MaxEntries != 0 {
			if game.EntriesLeft == 1 {
				gameString += "\nNext entry will be the story ending. Make it a good one!\n"
			} else {
				gameString += fmt.Sprintf("Entires left: %d \n", game.EntriesLeft)
			}
		}
	}

	return gameString
}

// StartGame creates a game, initializing all required structures and arrays, with the provided players and initiator.
// Supports configuration of time limit for turns (in seconds) and max length of entries (in symbols). Pass 0 if you don't want any of these features.
func StartGame(initiator string, players []string, timeLimit, maxLength, entriesCount int) *Game {
	playersCopy := make([]string, len(players))
	copy(playersCopy, players)
	for index, player := range playersCopy {
		if player == initiator { // arrange players array so that the initiator is first
			playersCopy = append([]string{initiator}, append(playersCopy[:index], playersCopy[index+1:]...)...)
			break
		}
	}

	game := &Game{
		Turn:        initiator,
		Story:       make([]Entry, 0),
		Players:     playersCopy,
		Finished:    false,
		TimeLeft:    timeLimit,
		MaxLength:   maxLength,
		MaxEntries:  entriesCount,
		EntriesLeft: entriesCount,
		VoteKick:    nil,

		playerTurn: 1,
		timeLimit:  timeLimit,
	}

	if timeLimit > 0 {
		go game.monitorTime()
	}

	return game
}

// AddEntry sets the game to end after the next turn.
func (game *Game) AddEntry(entry string, issuer string) error {
	if issuer != game.Turn {
		return errors.New("invalid entry - not this player's turn")
	}
	if game.MaxLength > 0 && len(entry) > game.MaxLength {
		return fmt.Errorf("invalid entry - entry is above max length (%v)", game.MaxLength)
	}

	game.Story = append(game.Story, Entry{Text: entry, Player: issuer})
	game.setNextTurn()

	if game.MaxEntries != 0 {
		game.EntriesLeft--
		if game.EntriesLeft <= 0 {
			game.Finished = true
		}
	}
	return nil
}

// EndGame sets the left entries count to one, meaning the next move will finish the story.
func (game *Game) EndGame(entries int) {
	game.MaxEntries = entries
	game.EntriesLeft = entries
}

// TriggerVoteKick starts a vote to kick a player. It requires an issuer on whose behalf the vote is triggered.
// Parameters of the campaign are the acceptance ratio (a number between 0 and 1, indicating what part of the players must submit a vote in order to kick the player)
// and a time limit (how many seconds before the campaign is considered unsuccessful)
// Return error if there is already a running vote or if the player to be kicked is not in the game.
func (game *Game) TriggerVoteKick(issuer, playerToKick string, acceptanceRatio float64, timeLimit int) error {
	if game.Finished {
		return errors.New("there is no running game")
	}
	if game.VoteKick != nil {
		return fmt.Errorf("there is an ongoing vote to kick player \"%s\"", game.VoteKick.Player)
	}
	for _, player := range game.Players {
		if player == playerToKick {
			voteTreshold := int(math.Ceil(float64(len(game.Players)) * acceptanceRatio))
			game.VoteKick = NewVoteKick(issuer, playerToKick, voteTreshold, timeLimit)
			go game.monitorVote()
			return nil
		}
	}
	return fmt.Errorf("player \"%s\" is not in the game", playerToKick)
}

// Vote submits a vote on behalf of the provided voter to the current campaign.
// Returns errors if the issuer has already voted or he's not part of the game or if there is no ongoing vote at all.
func (game *Game) Vote(voter string) error {
	if game.Finished {
		return errors.New("there is no running game")
	}
	if game.VoteKick == nil {
		return errors.New("there is no ongoing vote")
	}
	for _, player := range game.Players {
		if player == voter {
			if !game.VoteKick.hasVoted(voter) {
				game.VoteKick.voted = append(game.VoteKick.voted, player)
				game.VoteKick.Count++
				return nil
			}
			return fmt.Errorf("player \"%s\" has already voted for this vote", voter)
		}
	}
	return fmt.Errorf("player \"%s\" cannot vote as he's not part of the game", voter)
}

// Kick kicks a player from the game immediately, iterating player turn if necessary.
// Returns error of the player to be kicked is not part of the game
func (game *Game) Kick(toRemove string) error {
	for index, player := range game.Players {
		if player == toRemove {
			game.Players = util.DeleteFromSlice(game.Players, index)
			if player == game.Turn {
				game.setNextTurn()
			}
			return nil
		}
	}
	return fmt.Errorf("player \"%s\" is not part of the game", toRemove)
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

	if len(game.Players) > 0 {
		game.Turn = game.Players[game.playerTurn-1]
		game.TimeLeft = game.timeLimit
	} else {
		game.Finished = true
		game.Turn = ""
	}
}

func (game *Game) monitorVote() {
	for game.VoteKick != nil {
		if game.VoteKick.Count >= game.VoteKick.Treshold {
			game.Kick(game.VoteKick.Player)
			game.VoteKick = nil
			return
		}
		time.Sleep(1 * time.Second)
		game.VoteKick.TimeLeft--
		if game.VoteKick.TimeLeft <= 0 {
			game.VoteKick = nil
			return
		}
	}
}
