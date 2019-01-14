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

// CreateRoomCmd is a wrapper for the story-builder create-room command
type CreateRoomCmd struct{}

// Command builds and returns a cobra command that will be added to the root command
func (crc *CreateRoomCmd) Command() *cobra.Command {
	result := crc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (crc *CreateRoomCmd) Run() error {
	fmt.Println("create-room called")
	return nil
}

func (crc *CreateRoomCmd) buildCommand() *cobra.Command {
	var createRoomCmd = &cobra.Command{
		Use:   "create-room",
		Short: "",
		Long:  ``,
		RunE:  cmd.RunE(crc),
	}
	return createRoomCmd
}
