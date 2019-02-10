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

	roomName string
}

// Command builds and returns a cobra command that will be added to the root command
func (jrc *JoinRoomCmd) Command() *cobra.Command {
	result := jrc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (jrc *JoinRoomCmd) Run() error {
	cfg, err := jrc.Configurator.Load()
	if err != nil {
		return err
	}
	if err := cfg.ValidateConnection(); err != nil {
		return fmt.Errorf("there is no valid connection with a server: %v", err)
	}
	if cfg.Authorization == "" {
		return errors.New("users is not logged in")
	}
	if cfg.Room != "" {
		return errors.New("user is already in a room")
	}
	if jrc.roomName == "" {
		return errors.New("room name is empty")
	}

	room, err := jrc.Client.GetRoom(jrc.roomName)
	if err != nil {
		return err
	}
	cfg.Room = room.Name
	if err := jrc.Configurator.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("You've successfully joined room \"%s\".\n", room.Name)
	return nil
}

func (jrc *JoinRoomCmd) buildCommand() *cobra.Command {
	var joinRoomCmd = &cobra.Command{
		Use:   "join-room",
		Aliases: []string{"jr"},
		Short: "",
		Long:  ``,
		RunE:  cmd.RunE(jrc),
	}

	joinRoomCmd.Flags().StringVarP(&jrc.roomName, "name", "n", "", "name of the room to join")

	return joinRoomCmd
}
