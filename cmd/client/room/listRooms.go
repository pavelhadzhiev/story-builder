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

// ListRoomsCmd is a wrapper for the story-builder list-rooms command
type ListRoomsCmd struct{}

// Command builds and returns a cobra command that will be added to the root command
func (lrc *ListRoomsCmd) Command() *cobra.Command {
	result := lrc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (lrc *ListRoomsCmd) Run() error {
	fmt.Println("list-rooms called")
	return nil
}

func (lrc *ListRoomsCmd) buildCommand() *cobra.Command {
	var listRoomsCmd = &cobra.Command{
		Use:   "list-rooms",
		Short: "",
		Long:  ``,
		RunE:  cmd.RunE(lrc),
	}
	return listRoomsCmd
}