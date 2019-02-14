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

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/spf13/cobra"
)

// LeaveRoomCmd is a wrapper for the story-builder leave-room command
type LeaveRoomCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (lrc *LeaveRoomCmd) Command() *cobra.Command {
	result := lrc.buildCommand()

	return result
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (lrc *LeaveRoomCmd) RequiresConnection() *cmd.Context {
	return lrc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (lrc *LeaveRoomCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (lrc *LeaveRoomCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (lrc *LeaveRoomCmd) Run() error {
	cfg, err := lrc.Configurator.Load()
	if err != nil {
		return err
	}

	var roomName = cfg.Room
	cfg.Room = ""
	if err := lrc.Configurator.Save(cfg); err != nil {
		return err
	}

	if err := lrc.Client.LeaveRoom(roomName); err != nil {
		fmt.Printf("Something went wrong: %s\n", err)
		fmt.Println("This could've been caused by a problem with the server.")
		fmt.Println("However, your configuration has been adjusted, so you have technically \"left\".")
	} else {
		fmt.Printf("You've successfully left room \"%s\".\n", roomName)
	}

	return nil
}

func (lrc *LeaveRoomCmd) buildCommand() *cobra.Command {
	var leaveRoomsCmd = &cobra.Command{
		Use:     "leave-room",
		Aliases: []string{"lr"},
		Short:   "Leaves the room in the user configuration.",
		Long:    `Leaves the room in the user configuration. It's possible that the room doesn't exist or the player is not in it due to a server error or outage. In this case, the configuratin will be adjusted without any server side changes.`,
		PreRunE: cmd.PreRunE(lrc),
		RunE:    cmd.RunE(lrc),
	}
	return leaveRoomsCmd
}
