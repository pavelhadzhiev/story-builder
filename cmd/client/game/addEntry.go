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
	"github.com/spf13/cobra"
)

// AddEntryCmd is a wrapper for the story-builder add command
type AddEntryCmd struct {
	*cmd.Context

	entry string
}

// Command builds and returns a cobra command that will be added to the root command
func (aec *AddEntryCmd) Command() *cobra.Command {
	result := aec.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (aec *AddEntryCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	aec.entry = args[0]
	return nil
}

// Run is used to build the RunE function for the cobra command
func (aec *AddEntryCmd) Run() error {
	cfg, err := aec.Configurator.Load()
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

	if err := aec.Client.AddEntry(cfg.Room, aec.entry); err != nil {
		return err
	}

	fmt.Println("You've successfully submitted your entry.")
	fmt.Println("You can use the get-game command to check the story status.")
	return nil
}

func (aec *AddEntryCmd) buildCommand() *cobra.Command {
	var addEntryCmd = &cobra.Command{
		Use:     "add-entry [entry]",
		Aliases: []string{"add"},
		Short:   "Adds an entry to the current game if its your turn.",
		Long:    `Adds an entry to the current game if its your turn. If the game is not started or finished or its not your turn, returns error.`,
		PreRunE: cmd.PreRunE(aec),
		RunE:    cmd.RunE(aec),
	}
	return addEntryCmd
}
