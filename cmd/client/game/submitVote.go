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
	"errors"
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/spf13/cobra"
)

// SubmitVoteCmd is a wrapper for the story-builder submit vote command
type SubmitVoteCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (svc *SubmitVoteCmd) Command() *cobra.Command {
	result := svc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (svc *SubmitVoteCmd) Run() error {
	cfg, err := svc.Configurator.Load()
	if err != nil {
		return err
	}
	if err := cfg.ValidateConnection(); err != nil {
		return fmt.Errorf("there is no valid connection with a server: %v", err)
	}
	if cfg.Authorization == "" {
		return errors.New("users is not logged in")
	}
	if cfg.Room == "" {
		return errors.New("user is not in a room")
	}

	if err := svc.Client.SubmitVote(cfg.Room); err != nil {
		return err
	}

	fmt.Println("You've successfully cast your vote.")
	fmt.Println("You can use the get-game command to check the vote status.")
	return nil
}

func (svc *SubmitVoteCmd) buildCommand() *cobra.Command {
	var submitVoteCmd = &cobra.Command{
		Use:     "submit-vote",
		Aliases: []string{"sv"},
		Short:   "Submits your approval of the ongoing voting.",
		Long:    `Submits your approval of the ongoing voting. Returns error if a game is not running or there is no currently ongoing vote.`,
		PreRunE: cmd.PreRunE(svc),
		RunE:    cmd.RunE(svc),
	}
	return submitVoteCmd
}
