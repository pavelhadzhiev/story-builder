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

// LogoutCmd is a wrapper for the story-builder logout command
type LogoutCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (lc *LogoutCmd) Command() *cobra.Command {
	result := lc.buildCommand()

	return result
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (lc *LogoutCmd) RequiresConnection() *cmd.Context {
	return lc.Context
}

// RequiresAuthorization marks the command to require the configuration to have a user logged in.
func (lc *LogoutCmd) RequiresAuthorization() {}

// Run is used to build the RunE function for the cobra command
func (lc *LogoutCmd) Run() error {
	cfg, err := lc.Configurator.Load()
	if err != nil {
		return err
	}
	if cfg.Room != "" {
		lc.Client.LeaveRoom(cfg.Room)
		cfg.Room = ""
	}
	if err := lc.Client.Logout(); err != nil {
		return err
	}
	cfg.Authorization = ""
	lc.Configurator.Save(cfg)

	fmt.Println("You've logged out successfully. See you soon!")
	return nil
}

func (lc *LogoutCmd) buildCommand() *cobra.Command {
	var logoutCmd = &cobra.Command{
		Use:     "logout",
		Short:   "Logs the user out.",
		Long:    `Logs the user out. If a user is not logged in, rejects the request. Requires a valid connection to a server. If it is missing, a sufficient error message is provided. To connect to a server check the connect command.`,
		PreRunE: cmd.PreRunE(lc),
		RunE:    cmd.RunE(lc),
	}
	return logoutCmd
}
