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

// StartGameCmd is a wrapper for the story-builder start-game command
type StartGameCmd struct{}

// Command builds and returns a cobra command that will be added to the root command
func (sgc *StartGameCmd) Command() *cobra.Command {
	result := sgc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (sgc *StartGameCmd) Run() error {
	fmt.Println("start-game called")
	return nil
}

func (sgc *StartGameCmd) buildCommand() *cobra.Command {
	var startGameCmd = &cobra.Command{
		Use:   "start-game",
		Short: "",
		Long:  ``,
		RunE:  cmd.RunE(sgc),
	}
	return startGameCmd
}
