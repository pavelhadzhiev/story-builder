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

// GetGameCmd is a wrapper for the story-builder story command
type GetGameCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (ggc *GetGameCmd) Command() *cobra.Command {
	result := ggc.buildCommand()

	return result
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (ggc *GetGameCmd) RequiresConnection() *cmd.Context {
	return ggc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (ggc *GetGameCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (ggc *GetGameCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (ggc *GetGameCmd) Run() error {
	game, err := ggc.Client.GetGame()
	if err != nil {
		return err
	}
	fmt.Println(game)
	return nil
}

func (ggc *GetGameCmd) buildCommand() *cobra.Command {
	var getGameCmd = &cobra.Command{
		Use:     "get-game",
		Aliases: []string{"game"},
		Short:   "Prints the status of the current or last played game.",
		Long:    `Prints the status of the current or last played game. If no game has ever been started in the room, returns error.`,
		PreRunE: cmd.PreRunE(ggc),
		RunE:    cmd.RunE(ggc),
	}
	return getGameCmd
}
