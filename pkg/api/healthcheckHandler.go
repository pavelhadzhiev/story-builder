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

// HealthcheckHandler is an http handler for the story builder's healthcheck endpoint
func (server *SBServer) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/healthcheck/" {
		w.WriteHeader(404)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Returns status code 200 to show server is online and healthy
		fmt.Print("Healthcheck endpoint was called: ", r, "\n")
		w.WriteHeader(200)
	default:
		w.WriteHeader(405)
	}
}
