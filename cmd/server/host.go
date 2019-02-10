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
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/pavelhadzhiev/story-builder/pkg/db"

	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/pkg/api"
	"github.com/spf13/cobra"
)

// HostCmd is a wrapper for the story-builder host command
type HostCmd struct {
	database *db.SBDatabase
	username string
	password string

	port int
}

// Command builds and returns a cobra command that will be added to the root command
func (hc *HostCmd) Command() *cobra.Command {
	result := hc.buildCommand()

	return result
}

// Validate makes sure all required arguments are legal and are provided
func (hc *HostCmd) Validate(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("requires a single arg or no args")
	}

	if hc.username == "" && hc.password == "" { // set default database user if not provided
		hc.username = "admin"
		hc.password = "Abcd1234"
	}

	if len(args) == 0 { // set default server port if not provided
		hc.port = 8080
		return nil
	}

	portString := args[0]
	portNumber, err := strconv.Atoi(portString)
	if err != nil || portNumber < 0 {
		return fmt.Errorf("provided port \"%s\" is not valid", portString)
	}
	hc.port = portNumber
	return nil
}

// Run is used to build the RunE function for the cobra command
func (hc *HostCmd) Run() error {
	hc.database = db.NewSBDatabase(hc.username, hc.password)
	defer hc.database.CloseDB()
	err := hc.database.InitializeDB()
	if err != nil {
		return err
	}

	sbServer := api.NewSBServer(hc.database, hc.port)
	sbServer.Start()
	defer sbServer.Shutdown()

	fmt.Println("Server was started at http://localhost:", hc.port)
	fmt.Print("Press ENTER to shut down...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	return nil
}

func (hc *HostCmd) buildCommand() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:     "host [port]",
		Short:   "Hosts a server at the specified port.",
		Long:    `Hosts a server at the specified port. If the port is invalid or the server cannot be started, a sufficient errog message is returned.`,
		PreRunE: cmd.PreRunE(hc),
		RunE:    cmd.RunE(hc),
	}

	serverCmd.Flags().StringVarP(&hc.username, "username", "u", "", `Username to access database with. Default value is "admin"`)
	serverCmd.Flags().StringVarP(&hc.password, "password", "p", "", `Password to access database with. Default value is "Abcd1234"`)

	return serverCmd
}
