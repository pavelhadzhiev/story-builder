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
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/spf13/cobra"
)

// StartGameCmd is a wrapper for the story-builder start-game command
type StartGameCmd struct {
	*cmd.Context

	timeLimit    int
	maxLength    int
	entriesCount int
}

// Command builds and returns a cobra command that will be added to the root command
func (sgc *StartGameCmd) Command() *cobra.Command {
	result := sgc.buildCommand()

	return result
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (sgc *StartGameCmd) RequiresConnection() *cmd.Context {
	return sgc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (sgc *StartGameCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (sgc *StartGameCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (sgc *StartGameCmd) Run() error {
	if err := sgc.Client.StartGame(sgc.timeLimit, sgc.maxLength, sgc.entriesCount); err != nil {
		return err
	}

	fmt.Println("You've successfully started a game.")
	fmt.Println("You can use the get-game and add-entry commands to play.")
	return nil
}

func (sgc *StartGameCmd) buildCommand() *cobra.Command {
	var startGameCmd = &cobra.Command{
		Use:     "start-game",
		Aliases: []string{"sg"},
		Short:   "Starts a game in the joined room.",
		Long:    `Starts a game in the joined room. Requires admin access. If a game is already started, returns error. Supports configurations of time limit per turn and max length of entries. Default values are 60 seconds and 100 symbols. If you don't want to use any of these features, pass 0 with the according flag.`,
		PreRunE: cmd.PreRunE(sgc),
		RunE:    cmd.RunE(sgc),
	}

	startGameCmd.Flags().IntVarP(&sgc.timeLimit, "time", "t", 60, "the time limit to complete a turn in seconds")
	startGameCmd.Flags().IntVarP(&sgc.maxLength, "length", "l", 100, "the max length for an entry in symbols")
	startGameCmd.Flags().IntVarP(&sgc.entriesCount, "entires", "e", 0, "the amount of entries that will be played out before the game ends")

	return startGameCmd
}
