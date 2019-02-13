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

import "fmt"

// VoteKick represents a vote to kick a player from a game.
// If enough votes are submitted to cover the vote treshold, the player will be removed from the game.
type VoteKick struct {
	Player   string `json:"player,omitempty"`
	TimeLeft int    `json:"timeLeft,omitempty"`
	Count    int    `json:"voteCount,omitempty"`
	Treshold int    `json:"voteTreshold,omitempty"`
	Issuer   string `json:"issuer,omitempty"`

	voted []string
}

func NewVoteKick(issuer, player string, treshold, timeleft int) *VoteKick {
	return &VoteKick{
		Player:   player,
		TimeLeft: timeleft,
		Count:    0,
		Treshold: treshold,
		Issuer:   issuer,

		voted: make([]string, 0),
	}
}

func (voteKick *VoteKick) String() (voteKickString string) {
	voteKickString = fmt.Sprintf("Triggered by: \"%s\"\n", voteKick.Issuer)
	voteKickString += fmt.Sprintf("Player to kick: \"%s\"\n", voteKick.Player)
	voteKickString += fmt.Sprintf("Required votes: %d\n", voteKick.Treshold)
	voteKickString += fmt.Sprintf("Votes so far: %d\n", voteKick.Count)
	voteKickString += fmt.Sprintf("Time left until vote end: %d seconds\n", voteKick.TimeLeft)
	return
}

func (voteKick *VoteKick) hasVoted(player string) bool {
	for _, voter := range voteKick.voted {
		if player == voter {
			return true
		}
	}
	return false
}
