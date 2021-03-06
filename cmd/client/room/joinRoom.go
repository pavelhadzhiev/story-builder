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
	"github.com/spf13/cobra"
)

// JoinRoomCmd is a wrapper for the story-builder join-room command
type JoinRoomCmd struct {
	*cmd.Context

	name string
}

// Command builds and returns a cobra command that will be added to the root command
func (jrc *JoinRoomCmd) Command() *cobra.Command {
	result := jrc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (jrc *JoinRoomCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	jrc.name = args[0]
	return nil
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (jrc *JoinRoomCmd) RequiresConnection() *cmd.Context {
	return jrc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (jrc *JoinRoomCmd) RequiresAuthorization() {}

// Run is used to build the RunE function for the cobra command
func (jrc *JoinRoomCmd) Run() error {
	cfg, err := jrc.Configurator.Load()
	if err != nil {
		return err
	}
	if cfg.Room != "" {
		return errors.New("user is already in a room")
	}

	if err := jrc.Client.JoinRoom(jrc.name); err != nil {
		return err
	}

	cfg.Room = jrc.name
	if err := jrc.Configurator.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("You've successfully joined room \"%s\".\n", jrc.name)
	return nil
}

func (jrc *JoinRoomCmd) buildCommand() *cobra.Command {
	var joinRoomCmd = &cobra.Command{
		Use:     "join-room [name]",
		Aliases: []string{"jr"},
		Short:   "Joins the room with the provided name.",
		Long:    `Joins the room with the provided name. If the room doesn't exist or the player is banned from it, an error is returned.`,
		PreRunE: cmd.PreRunE(jrc),
		RunE:    cmd.RunE(jrc),
	}

	return joinRoomCmd
}
