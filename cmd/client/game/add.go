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
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/spf13/cobra"
)

// AddCmd is a wrapper for the story-builder add command
type AddCmd struct{}

// Command builds and returns a cobra command that will be added to the root command
func (ac *AddCmd) Command() *cobra.Command {
	result := ac.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (ac *AddCmd) Run() error {
	fmt.Println("add called")
	return nil
}

func (ac *AddCmd) buildCommand() *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "",
		Long:  ``,
		RunE:  cmd.RunE(ac),
	}
	return addCmd
}