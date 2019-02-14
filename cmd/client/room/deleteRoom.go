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

package room

import (
	"errors"
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// DeleteRoomCmd is a wrapper for the story-builder delete-room command
type DeleteRoomCmd struct {
	*cmd.Context

	name string
}

// Command builds and returns a cobra command that will be added to the root command
func (drc *DeleteRoomCmd) Command() *cobra.Command {
	result := drc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (drc *DeleteRoomCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	drc.name = args[0]
	return nil
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (drc *DeleteRoomCmd) RequiresConnection() *cmd.Context {
	return drc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (drc *DeleteRoomCmd) RequiresAuthorization() {}

// Run is used to build the RunE function for the cobra command
func (drc *DeleteRoomCmd) Run() error {
	if drc.name == "" {
		return errors.New("room name is empty")
	}

	action := fmt.Sprintf("delete room \"%s\"", drc.name)
	if !util.ConfirmationPrompt(action) {
		fmt.Println("Operation cancelled. No action taken.")
		return nil
	}
	if err := drc.Client.DeleteRoom(drc.name); err != nil {
		return err
	}

	fmt.Printf("You've successfully deleted room \"%s\".\n", drc.name)
	return nil
}

func (drc *DeleteRoomCmd) buildCommand() *cobra.Command {
	var deleteRoomCmd = &cobra.Command{
		Use:     "delete-room [name]",
		Aliases: []string{"dr"},
		Short:   "Deletes the game room with the provided name. Requires the user to be the creator of the room.",
		Long:    `Deletes the game room with the provided name. Requires the user to be the creator of the room. Returns an error if a room with this name doesn't exist or the issuer isn't the room creator.`,
		PreRunE: cmd.PreRunE(drc),
		RunE:    cmd.RunE(drc),
	}

	return deleteRoomCmd
}
