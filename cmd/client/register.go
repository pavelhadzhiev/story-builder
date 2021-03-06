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
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/client"
	"github.com/pavelhadzhiev/story-builder/pkg/util"
	"github.com/spf13/cobra"
)

// RegisterCmd is a wrapper for the story-builder register command
type RegisterCmd struct {
	*cmd.Context

	username string
	password string
}

// Command builds and returns a cobra command that will be added to the root command
func (rc *RegisterCmd) Command() *cobra.Command {
	result := rc.buildCommand()

	return result
}

// RequiresConnection makes sure that the configured server is valid and online before executing the command logic
func (rc *RegisterCmd) RequiresConnection() *cmd.Context {
	return rc.Context
}

// Run is used to build the RunE function for the cobra command
func (rc *RegisterCmd) Run() error {
	cfg, err := rc.Configurator.Load()
	if err != nil {
		return err
	}
	if cfg.Authorization != "" {
		return errors.New("users is already logged in")
	}
	if rc.username == "" {
		username, err := util.ReadUsername()
		if err != nil {
			return err
		}
		rc.username = username
	}
	if rc.password == "" {
		password, err := util.ReadPassword()
		if err != nil {
			return err
		}
		rc.password = password
	}
	cfg.Authorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(rc.username+":"+rc.password))
	rc.Configurator.Save(cfg)

	if err := client.NewSBClient(cfg).Register(); err != nil {
		cfg.Authorization = ""
		rc.Configurator.Save(cfg)
		return err
	}

	fmt.Printf("You've registered successfully! Welcome, %s.\n", rc.username)
	return nil
}

func (rc *RegisterCmd) buildCommand() *cobra.Command {
	var registerCmd = &cobra.Command{
		Use:     "register",
		Short:   "Registers a non-existing user with the provided username and password.",
		Long:    `Registers a non-existing user with the provided username and password. If the user does exists, the request is rejected. Requires a valid connection to a server, a username and a password. If any of these are missing a sufficient error message is provided. To connect to a server check the connect command.`,
		PreRunE: cmd.PreRunE(rc),
		RunE:    cmd.RunE(rc),
	}

	registerCmd.Flags().StringVarP(&rc.username, "username", "u", "", "username to register")
	registerCmd.Flags().StringVarP(&rc.password, "password", "p", "", "password to register with")

	return registerCmd
}
