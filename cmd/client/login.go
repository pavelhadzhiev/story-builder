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

package client

import (
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/spf13/cobra"
)

// LoginCmd is a wrapper for the story-builder login command
type LoginCmd struct{}

// Command builds and returns a cobra command that will be added to the root command
func (lc *LoginCmd) Command() *cobra.Command {
	result := lc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (lc *LoginCmd) Run() error {
	fmt.Println("login called")
	return nil
}

func (lc *LoginCmd) buildCommand() *cobra.Command {
	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "",
		Long:  ``,
		RunE:  cmd.RunE(lc),
	}
	return loginCmd
}