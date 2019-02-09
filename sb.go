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

package main

import (
	"github.com/pavelhadzhiev/story-builder/cmd"
	"github.com/pavelhadzhiev/story-builder/cmd/client"
	"github.com/pavelhadzhiev/story-builder/cmd/client/game"
	"github.com/pavelhadzhiev/story-builder/cmd/client/room"
	"github.com/pavelhadzhiev/story-builder/cmd/server"
)

func main() {
	ctx := &cmd.Context{}
	rootCmd := cmd.BuildRootCommand(ctx)

	commands := []cmd.CommandWrapper{
		&server.HostCmd{},
		&server.ConnectCmd{Context: ctx},
		&server.DisconnectCmd{Context: ctx},
		&client.LoginCmd{Context: ctx},
		&client.LogoutCmd{Context: ctx},
		&client.RegisterCmd{Context: ctx},
		&room.CreateRoomCmd{},
		&room.DeleteRoomCmd{},
		&room.JoinRoomCmd{},
		&room.LeaveRoomCmd{},
		&room.ListRoomsCmd{},
		&game.StartGameCmd{},
		&game.EndGameCmd{},
		&game.AddCmd{},
		&game.StoryCmd{},
		&game.UsersCmd{},
	}
	for _, command := range commands {
		rootCmd.AddCommand(command.Command())
	}
	cmd.Execute(rootCmd)
}
