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
	"errors"
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
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

// Run is used to build the RunE function for the cobra command
func (bc *BanCmd) Run() error {
	cfg, err := bc.Configurator.Load()
	if err != nil {
		return err
	}
	if err := cfg.ValidateConnection(); err != nil {
		return fmt.Errorf("there is no valid connection with a server: %v", err)
	}
	if cfg.Authorization == "" {
		return errors.New("users is not logged in")
	}
	if cfg.Room == "" {
		return errors.New("user has not joined in a room")
	}

	if err := bc.Client.BanPlayer(cfg.Room, bc.player); err != nil {
		return err
	}

	fmt.Printf("You've banned \"%s\" from room \"%s\".\n", bc.player, cfg.Room)
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
