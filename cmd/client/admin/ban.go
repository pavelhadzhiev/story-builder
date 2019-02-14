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

package admin

import (
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// BanCmd is a wrapper for the story-builder ban command
type BanCmd struct {
	*cmd.Context

	player string
}

// Command builds and returns a cobra command that will be added to the root command
func (bc *BanCmd) Command() *cobra.Command {
	result := bc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (bc *BanCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	bc.player = args[0]
	return nil
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (bc *BanCmd) RequiresConnection() *cmd.Context {
	return bc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (bc *BanCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (bc *BanCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (bc *BanCmd) Run() error {
	action := fmt.Sprintf("ban player \"%s\"", bc.player)
	if !util.ConfirmationPrompt(action) {
		fmt.Println("Operation cancelled. No action taken.")
		return nil
	}
	if err := bc.Client.BanPlayer(bc.player); err != nil {
		return err
	}

	fmt.Printf("You've banned \"%s\".\n", bc.player)
	return nil
}

func (bc *BanCmd) buildCommand() *cobra.Command {
	var banCmd = &cobra.Command{
		Use:     "ban [player]",
		Short:   "An admin command that bans the player provided as argument.",
		Long:    `An admin command that bans the player provided as argument from the current room. Returns error if you don't have admin access for the room.`,
		PreRunE: cmd.PreRunE(bc),
		RunE:    cmd.RunE(bc),
	}

	return banCmd
}
