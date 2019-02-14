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

	"github.com/pavelhadzhiev/story-builder/pkg/client"
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

// Validate makes sure all required arguments are legal and are provided
func (cc *ConnectCmd) Validate(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("requires a single arg or no arg")
	}
	if len(args) == 0 { // set default server if not provided
		cc.host = "http://localhost:8080"
		return nil
	}

	cc.host = args[0]
	return nil
}

// Run is used to build the RunE function for the cobra command
func (cc *ConnectCmd) Run() error {
	cfg, err := cc.Configurator.Load()
	if err != nil {
		return err
	}
	if cfg.URL != "" {
		return fmt.Errorf("you are already connected to \"" + cfg.URL + "\". You have to disconnect first")
	}
	cfg.URL = cc.host
	cc.Client = client.NewSBClient(cfg)
	if err = cc.Client.HealthCheck(cc.Configurator); err != nil {
		return err
	}

	fmt.Println("You've successfully connected to \"" + cfg.URL + "\"! You can now log in as an existing user or register as a new one.")
	return nil
}

func (cc *ConnectCmd) buildCommand() *cobra.Command {
	var connectCmd = &cobra.Command{
		Use:     "connect [host]",
		Short:   "Connects to a healthy server with the provided host.",
		Long:    `Connects to a healthy server with the provided host. If the URL is invalid, or the specified server doesn't have a responing healthcheck, the request is rejected.`,
		PreRunE: cmd.PreRunE(cc),
		RunE:    cmd.RunE(cc),
	}

	return connectCmd
}
