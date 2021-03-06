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
	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// VoteCmd is a wrapper for the story-builder submit vote command
type VoteCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (vc *VoteCmd) Command() *cobra.Command {
	result := vc.buildCommand()

	return result
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (vc *VoteCmd) RequiresConnection() *cmd.Context {
	return vc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (vc *VoteCmd) RequiresAuthorization() {}

// RequiresRoom marks the command to require the configuration to have a user logged in.
func (vc *VoteCmd) RequiresRoom() {}

// Run is used to build the RunE function for the cobra command
func (vc *VoteCmd) Run() error {
	action := "submit your approval of the currently running vote"
	if !util.ConfirmationPrompt(action) {
		fmt.Println("Operation cancelled. No action taken.")
		return nil
	}
	if err := vc.Client.SubmitVote(); err != nil {
		return err
	}

	fmt.Println("You've successfully cast your vote.")
	fmt.Println("You can use the get-game command to check the vote status.")
	return nil
}

func (vc *VoteCmd) buildCommand() *cobra.Command {
	var submitVoteCmd = &cobra.Command{
		Use:     "vote",
		Short:   "Submits your approval of the ongoing voting.",
		Long:    `Submits your approval of the ongoing voting. Returns error if a game is not running or there is no currently ongoing vote.`,
		PreRunE: cmd.PreRunE(vc),
		RunE:    cmd.RunE(vc),
	}
	return submitVoteCmd
}
