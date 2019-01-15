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
)

// StartStoryBuilderServer starts a story builder server at localhost:<port>
func StartStoryBuilderServer(port int) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/rooms/", RoomHandler)
	http.HandleFunc("/register/", RegistrationHandler)
	http.HandleFunc("/login/", LoginHandler)
	http.HandleFunc("/healthcheck/", healthcheck)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	return srv
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
