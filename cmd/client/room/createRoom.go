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
	"fmt"

	"github.com/pavelhadzhiev/story-builder/pkg/util"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/spf13/cobra"
)

// CreateRoomCmd is a wrapper for the story-builder create-room command
type CreateRoomCmd struct {
	*cmd.Context

	name string
}

// Command builds and returns a cobra command that will be added to the root command
func (crc *CreateRoomCmd) Command() *cobra.Command {
	result := crc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (crc *CreateRoomCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	crc.name = args[0]
	return nil
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (crc *CreateRoomCmd) RequiresConnection() *cmd.Context {
	return crc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (crc *CreateRoomCmd) RequiresAuthorization() {}

// Run is used to build the RunE function for the cobra command
func (crc *CreateRoomCmd) Run() error {
	cfg, err := crc.Configurator.Load()
	if err != nil {
		return err
	}

	user, err := util.ExtractUsernameFromAuthorizationHeader(cfg.Authorization)
	if err != nil {
		return err
	}
	if err := crc.Client.CreateNewRoom(rooms.NewRoom(crc.name, user)); err != nil {
		return err
	}

	fmt.Printf("Room \"%s\" was successfully created by \"%s\".\n", crc.name, user)
	return nil
}

func (crc *CreateRoomCmd) buildCommand() *cobra.Command {
	var createRoomCmd = &cobra.Command{
		Use:     "create-room [name]",
		Aliases: []string{"cr"},
		Short:   "Creates a game room with the provided name.",
		Long:    `Creates a game room with the provided name. Returns an error if a room with this name already exists.`,
		PreRunE: cmd.PreRunE(crc),
		RunE:    cmd.RunE(crc),
	}

	return createRoomCmd
}
