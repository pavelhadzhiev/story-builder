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

// RegistrationHandler is an http handler for the story builder's registration endpoint
func (server *SBServer) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register/" {
		w.WriteHeader(404)
		return
	}
	switch r.Method {
	case http.MethodPost:
		// Registers user (writes them in DB)
		fmt.Print("Register endpoint was called: ", r, "\n")
		w.Write([]byte("Let's say you've registered.\n"))
	default:
		w.WriteHeader(405)
	}
}
