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

// Room represents a story builder room, in which a select group of players can play the game
type Room struct {
	Name    string
	Creator string
	Rules   RoomRules

	Turn    string
	Players []string
	Story   []Entry
}

// Entry represents a single player's turn in the story builder game
type Entry struct {
	Text   string
	Player string
}

// RoomRules keeps some configurations for the gameplay in the story builder room.
type RoomRules struct {
	Timeout int
}
