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

import "github.com/spf13/cobra"

// CommandWrapper represents a wrapper struct for a cobra command. It has a method to build the cobra struct.
type CommandWrapper interface {
	Command() *cobra.Command
}

// Command represents a cobra command. It has a run method that executes the command's logic.
type Command interface {
	Run() error
}

// ValidatedCommand should be implemented if a validation method is needed.
type ValidatedCommand interface {
	Validate([]string) error
}

// PreRunE is used to execute some generic preparations for the command execution, depending on interfaces the command impements.
// Set this function to the PreRunE property of a cobra command.
func PreRunE(cmd Command) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		if valCmd, ok := cmd.(ValidatedCommand); ok {
			if err := valCmd.Validate(args); err != nil {
				return err
			}
		}
		c.SilenceUsage = true
		return nil
	}
}

// RunE is used to set the RunE property of a cobra command.
func RunE(cmd Command) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		return cmd.Run()
	}
}
