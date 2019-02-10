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

// ListRoomsCmd is a wrapper for the story-builder list-rooms command
type ListRoomsCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (lrc *ListRoomsCmd) Command() *cobra.Command {
	result := lrc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (lrc *ListRoomsCmd) Run() error {
	cfg, err := lrc.Configurator.Load()
	if err != nil {
		return err
	}
	if err := cfg.ValidateConnection(); err != nil {
		return fmt.Errorf("there is no valid connection with a server: %v", err)
	}
	if cfg.Authorization == "" {
		return errors.New("users is not logged in")
	}

	rooms, err := lrc.Client.GetAllRooms()
	if err != nil {
		return err
	}

	fmt.Println("All rooms in the current server are:")
	for _, room := range rooms {
		fmt.Println(room)
	}
	return nil
}

func (lrc *ListRoomsCmd) buildCommand() *cobra.Command {
	var listRoomsCmd = &cobra.Command{
		Use:     "list-rooms",
		Aliases: []string{"list"},
		Short:   "Lists all room in the current server.",
		Long:    `Lists all room in the current server.`,
		PreRunE: cmd.PreRunE(lrc),
		RunE:    cmd.RunE(lrc),
	}
	return listRoomsCmd
}
