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

	"github.com/pavelhadzhiev/story-builder/pkg/util"
)

// LogoutHandler is an http handler for the story builder's logout endpoint
func (server *SBServer) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout/" {
		w.WriteHeader(404)
		return
	}

	switch r.Method {
	case http.MethodPost:
		username, password, err := util.ExtractCredentialsFromAuthorizationHeader(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}

		if err := server.Database.LoginUser(username, password); err != nil {
			w.WriteHeader(401)
			w.Write([]byte("Could not authenticate user."))
			return
		}

		for index, user := range server.Online {
			if user == username {
				server.Online = append(server.Online[:index], server.Online[index+1:]...)
			}
		}
		w.Write([]byte("Successfully logged out. See you soon, " + username + "!"))
	default:
		w.WriteHeader(405)
	}
}
