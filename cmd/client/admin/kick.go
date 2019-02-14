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

// KickCmd is a wrapper for the story-builder kick command
type KickCmd struct {
	*cmd.Context

	player string
}

// Command builds and returns a cobra command that will be added to the root command
func (kc *KickCmd) Command() *cobra.Command {
	result := kc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (kc *KickCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	kc.player = args[0]
	return nil
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (kc *KickCmd) RequiresConnection() *cmd.Context {
	return kc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (kc *KickCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (kc *KickCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (kc *KickCmd) Run() error {
	action := fmt.Sprintf("kick player \"%s\" from the game", kc.player)
	if !util.ConfirmationPrompt(action) {
		fmt.Println("Operation cancelled. No action taken.")
		return nil
	}
	if err := kc.Client.KickPlayer(kc.player); err != nil {
		return err
	}

	fmt.Printf("You've kicked \"%s\" from the game.\n", kc.player)
	return nil
}

func (kc *KickCmd) buildCommand() *cobra.Command {
	var KickCmd = &cobra.Command{
		Use:     "kick [player]",
		Short:   "An admin command that kicks the player provided as argument.",
		Long:    `An admin command that kicks the player provided as argument from the game in the current room. Returns error if you don't have admin access for the room.`,
		PreRunE: cmd.PreRunE(kc),
		RunE:    cmd.RunE(kc),
	}

	return KickCmd
}
