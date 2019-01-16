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
	"github.com/spf13/cobra"
)

// DisconnectCmd is a wrapper for the story-builder connect command
type DisconnectCmd struct {
	*cmd.Context
}

// Command builds and returns a cobra command that will be added to the root command
func (cc *DisconnectCmd) Command() *cobra.Command {
	result := cc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (cc *DisconnectCmd) Run() error {
	cfg, err := cc.Configurator.Load()
	if err != nil {
		return err
	}
	if err = cfg.ValidateConnection(); err != nil {
		return fmt.Errorf("there is no valid connection with a server: %v", err)
	}
	cfg.URL = ""
	cc.Configurator.Save(cfg)

	fmt.Println("disconnect called")
	return nil
}

func (cc *DisconnectCmd) buildCommand() *cobra.Command {
	var disconnectCmd = &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnects from the connected server.",
		Long:  `Disconnects from the connected server. If there is none, a sufficient error message is returned.`,
		RunE:  cmd.RunE(cc),
	}
	return disconnectCmd
}
