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

// ConnectCmd is a wrapper for the story-builder connect command
type ConnectCmd struct {
	*cmd.Context

	host string
}

// Command builds and returns a cobra command that will be added to the root command
func (cc *ConnectCmd) Command() *cobra.Command {
	result := cc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (cc *ConnectCmd) Run() error {
	cfg, err := cc.Configurator.Load()
	if err != nil {
		return err
	}
	cfg.URL = cc.host
	if err = cfg.ValidateConnection(); err != nil {
		return err
	}
	cc.Configurator.Save(cfg)

	fmt.Println("connect called")
	return nil
}

func (cc *ConnectCmd) buildCommand() *cobra.Command {
	var connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connects to a healthy server with the provided host.",
		Long:  `Connects to a healthy server with the provided host. If the URL is invalid, or the specified server doesn't have a responing healthcheck, the request is rejected.`,
		RunE:  cmd.RunE(cc),
	}

	connectCmd.Flags().StringVarP(&cc.host, "host", "", "", "Host")

	return connectCmd
}
