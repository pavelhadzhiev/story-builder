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

package server

import (
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/config"
	"github.com/spf13/cobra"
)

// DisconnectCmd is a wrapper for the story-builder connect command
type DisconnectCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (dc *DisconnectCmd) Command() *cobra.Command {
	result := dc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (dc *DisconnectCmd) Run() error {
	cfg, err := dc.Configurator.Load()
	if err != nil {
		return err
	}
	if err = cfg.ValidateConnection(); err != nil {
		dc.Configurator.Save(&config.SBConfiguration{})
		return fmt.Errorf("there is no valid connection with a server: %v. Existing configuration will be reset", err)
	}
	if cfg.Room != "" {
		dc.Client.LeaveRoom(cfg.Room)
		cfg.Room = ""
	}
	if cfg.Authorization != "" {
		//dc.Client.Logout()
		cfg.Authorization = ""
	}
	cfg.URL = ""
	dc.Configurator.Save(cfg)

	fmt.Println("You've disconnected from the server.")
	return nil
}

func (dc *DisconnectCmd) buildCommand() *cobra.Command {
	var disconnectCmd = &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnects from the connected server.",
		Long:  `Disconnects from the connected server. If there is none, a sufficient error message is returned.`,
		RunE:  cmd.RunE(dc),
	}
	return disconnectCmd
}
