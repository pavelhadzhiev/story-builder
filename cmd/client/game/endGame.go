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
	"strconv"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// EndGameCmd is a wrapper for the story-builder end-game command
type EndGameCmd struct {
	*cmd.Context

	entriesCount int
}

// Command builds and returns a cobra command that will be added to the root command
func (egc *EndGameCmd) Command() *cobra.Command {
	result := egc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (egc *EndGameCmd) Validate(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("requires a single arg or no args")
	}

	if len(args) == 1 {
		count, err := strconv.Atoi(args[0])
		if err != nil || count <= 0 {
			return err
		}
		egc.entriesCount = count
	} else {
		egc.entriesCount = 1 // Set default end game countdown in case one is not provided
	}
	return nil
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (egc *EndGameCmd) RequiresConnection() *cmd.Context {
	return egc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (egc *EndGameCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (egc *EndGameCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (egc *EndGameCmd) Run() error {
	action := fmt.Sprintf("end game in room after %d moves", egc.entriesCount)
	if !util.ConfirmationPrompt(action) {
		fmt.Println("Operation cancelled. No action taken.")
		return nil
	}
	if err := egc.Client.EndGame(egc.entriesCount); err != nil {
		return err
	}

	fmt.Println("You've successfully triggered a game end.")
	return nil
}

func (egc *EndGameCmd) buildCommand() *cobra.Command {
	var endGameCmd = &cobra.Command{
		Use:     "end-game [entries-left]",
		Aliases: []string{"eg", "end"},
		Short:   "Ends the game in the joined room. Executing this means that game will have the provided number of entries left until it is finished. If not provided it will allow only one turn.",
		Long:    `Ends the game in the joined room. Executing this means that game will have the provided number of entries left until it is finished. If not provided it will allow only one turn. Requires admin access. If there is no running game, returns error.`,
		PreRunE: cmd.PreRunE(egc),
		RunE:    cmd.RunE(egc),
	}
	return endGameCmd
}
