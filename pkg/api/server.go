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

package api

import (
	"fmt"
	"net/http"

	"github.com/pavelhadzhiev/story-builder/pkg/api/rooms"

	"github.com/pavelhadzhiev/story-builder/pkg/db"
)

// SBServer implements the story builder server API. It contains a database and some configurations. Use the Start and Shutdown methods to manage.
type SBServer struct {
	Database db.UserDatabase
	Rooms    []rooms.Room
	Online   []string

	srv *http.Server
}

// NewSBServer returns a story builder server configured for localhost:<port> that will use the provided database
func NewSBServer(sbdb *db.SBDatabase, port int) (sbServer *SBServer) {
	sbServer = &SBServer{
		Database: sbdb,
		Rooms:    make([]rooms.Room, 0),
		Online:   make([]string, 0),

		srv: &http.Server{Addr: fmt.Sprintf(":%d", port)},
	}

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/healthcheck/", sbServer.HealthcheckHandler)

	http.HandleFunc("/register/", sbServer.RegistrationHandler)
	http.HandleFunc("/login/", sbServer.LoginHandler)
	http.HandleFunc("/logout/", sbServer.LogoutHandler)

	http.HandleFunc("/rooms/", sbServer.RoomHandler)
	http.HandleFunc("/join-room/", sbServer.JoinRoomHandler)
	http.HandleFunc("/leave-room/", sbServer.LeaveRoomHandler)

	http.HandleFunc("/vote/", sbServer.VoteHandler)
	http.HandleFunc("/gameplay/", sbServer.GameplayHandler)
	http.HandleFunc("/manage-games/", sbServer.ManageGamesHandler)

	http.HandleFunc("/admin/", sbServer.PromoteAdminHandler)
	http.HandleFunc("/admin/ban/", sbServer.BanHandler)
	http.HandleFunc("/admin/kick/", sbServer.KickHandler)

	return
}

// Start starts an HTTP server, using the available configuration
func (sbServer *SBServer) Start() {
	go func() {
		if err := sbServer.srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

// Shutdown stops the HTTP server gracefully. It's supposed to be used in a defer statement in the process that starts the server, in case of errors.
func (sbServer *SBServer) Shutdown() {
	if err := sbServer.srv.Shutdown(nil); err != nil {
		panic(err) // failed or timed out while shutting down the server
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}
