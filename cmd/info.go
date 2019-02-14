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

package cmd

import (
	"fmt"

	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// InfoCmd is a wrapper for the story-builder info command
type InfoCmd struct {
	*Context
}

// Command builds and returns a cobra command that will be added to the root command
func (ic *InfoCmd) Command() *cobra.Command {
	result := ic.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (ic *InfoCmd) Run() error {
	cfg, err := ic.Configurator.Load()
	if err != nil {
		fmt.Println("You don't have a valid configuration.")
		return nil
	}
	if cfg.URL == "" {
		fmt.Println("You are not connected to a server.")
		return nil
	}
	fmt.Println("Server:", cfg.URL)

	if cfg.Authorization == "" {
		fmt.Println("You are not logged in.")
		return nil
	}
	user, err := util.ExtractUsernameFromAuthorizationHeader(cfg.Authorization)
	if err != nil {
		fmt.Println("You are not logged in.")
		return nil
	}
	fmt.Println("User:", user)

	if cfg.Room == "" {
		fmt.Println("You have not joined a room.")
		return nil
	}
	fmt.Println("Room:", cfg.Room)

	return nil
}

func (ic *InfoCmd) buildCommand() *cobra.Command {
	var infoCmd = &cobra.Command{
		Use:     "info",
		Short:   "Outputs information about the server, user and joined room.",
		Long:    "Outputs information about the server, user and joined room.",
		PreRunE: PreRunE(ic),
		RunE:    RunE(ic),
	}

	return infoCmd
}
