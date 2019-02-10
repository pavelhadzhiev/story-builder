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

// DeleteRoomCmd is a wrapper for the story-builder delete-room command
type DeleteRoomCmd struct {
	*cmd.Context

	roomName string
}

// Command builds and returns a cobra command that will be added to the root command
func (drc *DeleteRoomCmd) Command() *cobra.Command {
	result := drc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (drc *DeleteRoomCmd) Run() error {
	cfg, err := drc.Configurator.Load()
	if err != nil {
		return err
	}
	if err := cfg.ValidateConnection(); err != nil {
		return fmt.Errorf("there is no valid connection with a server: %v", err)
	}
	if cfg.Authorization == "" {
		return errors.New("users is not logged in")
	}
	if drc.roomName == "" {
		return errors.New("room name is empty")
	}

	if err := drc.Client.DeleteRoom(drc.roomName); err != nil {
		return err
	}

	fmt.Printf("You've successfully deleted room \"%s\".\n", drc.roomName)
	return nil
}

func (drc *DeleteRoomCmd) buildCommand() *cobra.Command {
	var deleteRoomCmd = &cobra.Command{
		Use:   "delete-room",
		Aliases: []string{"dr"},
		Short: "",
		Long:  ``,
		RunE:  cmd.RunE(drc),
	}

	deleteRoomCmd.Flags().StringVarP(&drc.roomName, "name", "n", "", "name of the room to delete")

	return deleteRoomCmd
}
