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

	port string
}

// Command builds and returns a cobra command that will be added to the root command
func (hc *HostCmd) Command() *cobra.Command {
	result := hc.buildCommand()

	return result
}

// Run is used to build the RunE function for the cobra command
func (hc *HostCmd) Run() error {
	if hc.username == "" && hc.password == "" { // set default user if not provided
		hc.username = "admin"
		hc.password = "Abcd1234"
	}
	if hc.port == "" { // set default port if not provided
		hc.port = "8080"
	}

	hc.database = db.NewSBDatabase(hc.username, hc.password)
	defer hc.database.CloseDB()
	err := hc.database.InitializeDB()
	if err != nil {
		return err
	}

	if portNumber, err := strconv.Atoi(hc.port); err != nil || portNumber < 0 {
		return fmt.Errorf("provided port (%s) is not valid", hc.port)
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
		Use:   "host",
		Short: "Hosts a server at the specified port.",
		Long:  `Hosts a server at the specified port. If the port is invalid or the server cannot be started, a sufficient errog message is returned.`,
		RunE:  cmd.RunE(hc),
	}

	serverCmd.Flags().StringVarP(&hc.port, "port", "", "", `Port to host server on. Default value is "8080".`)
	serverCmd.Flags().StringVarP(&hc.username, "username", "u", "", `Username to access database with. Default value is "admin"`)
	serverCmd.Flags().StringVarP(&hc.password, "password", "p", "", `Password to access database with. Default value is "Abcd1234"`)

	return serverCmd
}
