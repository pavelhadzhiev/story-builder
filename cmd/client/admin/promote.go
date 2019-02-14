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

package admin

import (
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// PromoteCmd is a wrapper for the story-builder promote command
type PromoteCmd struct {
	*cmd.Context

	user string
}

// Command builds and returns a cobra command that will be added to the root command
func (pc *PromoteCmd) Command() *cobra.Command {
	result := pc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (pc *PromoteCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires a single arg")
	}

	pc.user = args[0]
	return nil
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (pc *PromoteCmd) RequiresConnection() *cmd.Context {
	return pc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (pc *PromoteCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (pc *PromoteCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (pc *PromoteCmd) Run() error {
	action := fmt.Sprintf("promote user \"%s\" to admin", pc.user)
	if !util.ConfirmationPrompt(action) {
		fmt.Println("Operation cancelled. No action taken.")
		return nil
	}
	if err := pc.Client.PromoteAdmin(pc.user); err != nil {
		return err
	}

	fmt.Printf("You've promoted \"%s\" to admin.\n", pc.user)
	return nil
}

func (pc *PromoteCmd) buildCommand() *cobra.Command {
	var PromoteCmd = &cobra.Command{
		Use:     "promote [user]",
		Short:   "An admin command that promotes the user provided as argument to an admin.",
		Long:    `An admin command that promotes the user provided as argument to an admin in the current room. Returns error if you don't have admin access for the room.`,
		PreRunE: cmd.PreRunE(pc),
		RunE:    cmd.RunE(pc),
	}

	return PromoteCmd
}
