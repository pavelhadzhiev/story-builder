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

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// TriggerVoteCmd is a wrapper for the story-builder trigger vote command
type TriggerVoteCmd struct {
	*cmd.Context

	playerToKick string
}

// Command builds and returns a cobra command that will be added to the root command
func (tvc *TriggerVoteCmd) Command() *cobra.Command {
	result := tvc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (tvc *TriggerVoteCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	tvc.playerToKick = args[0]
	return nil
}

// Run is used to build the RunE function for the cobra command
func (tvc *TriggerVoteCmd) Run() error {
	cfg, err := tvc.Configurator.Load()
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
		return errors.New("user is not in a room")
	}

	action := fmt.Sprintf("trigger a vote to kick player \"%s\" from the game in room \"%s\"", tvc.playerToKick, cfg.Room)
	if !util.ConfirmationPrompt(action) {
		fmt.Println("Operation cancelled. No action taken.")
		return nil
	}
	if err := tvc.Client.TriggerVoteKick(cfg.Room, tvc.playerToKick); err != nil {
		return err
	}

	fmt.Println("You've successfully triggered a vote to kick \"" + tvc.playerToKick + "\".")
	fmt.Println("You can use the get-game command to check the vote status.")
	return nil
}

func (tvc *TriggerVoteCmd) buildCommand() *cobra.Command {
	var triggerVoteCmd = &cobra.Command{
		Use:     "trigger-vote [player-to-kick]",
		Aliases: []string{"tv"},
		Short:   "Triggers a democratic vote to kick the provided player from the game.",
		Long:    `Triggers a democratic vote to kick the provided player from the game. Returns error if a game is not running, the player is not in the game or there is currently an ongoing vote.`,
		PreRunE: cmd.PreRunE(tvc),
		RunE:    cmd.RunE(tvc),
	}
	return triggerVoteCmd
}
