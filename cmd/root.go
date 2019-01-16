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
	"os"

	"github.com/pavelhadzhiev/story-builder/pkg/client"
	"github.com/pavelhadzhiev/story-builder/pkg/config"
	"github.com/spf13/cobra"
)

// BuildRootCommand builds the root command for story-builder
func BuildRootCommand(ctx *Context) *cobra.Command {
	var cfgFile string

	var rootCmd = &cobra.Command{
		Use:   "story-builder",
		Short: "Story Builder is a chat room game. All participants take turns in adding a sentence to a story and see the results that they get.",
		Long:  `Story Builder is a chat room game. All participants take turns in adding a sentence to a story and see the results that they get. To use it you have to connect to a server and register or log in with an existing user. Then join a room and start playing. Check out the help page for available commands and what they are used for.`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if ctx.Configurator == nil {
				configurator, err := config.NewViperConfigurator(cfgFile)
				if err != nil {
					return err
				}
				ctx.Configurator = configurator
			}
			config, err := ctx.Configurator.Load()
			if err != nil {
				return err
			}
			ctx.Client = client.NewStoryBuilderClient(config)
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.story-builder.json)")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
