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

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pavelhadzhiev/story-builder/pkg/api/game"
)

// GetGame retrieves the game of the room with the provided name.
// Returns error if room doesn't exist or game is not started.
func (client *SBClient) GetGame() (*game.Game, error) {
	roomName := client.config.Room
	response, err := client.call(http.MethodGet, "/gameplay/"+roomName, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		defer response.Body.Close()
		game := &game.Game{}
		if err := json.NewDecoder(response.Body).Decode(game); err != nil {
			return nil, fmt.Errorf("failed to deserialize response from server: %e", err)
		}
		return game, nil
	case 404:
		return nil, errors.New("room \"" + roomName + "\" doesn't exist or no games have been started")
	default:
		return nil, errors.New("something went really wrong :(")
	}
}

// AddEntry adds the provided entry in the game of the room with the provided name on behalf of the user.
// Returns error if room doesn't exist, game is not started or it's not the users turn.
func (client *SBClient) AddEntry(entry string) error {
	roomName := client.config.Room
	headers := make(map[string]string)
	headers["Entry-Text"] = entry
	response, err := client.call(http.MethodPost, "/gameplay/"+roomName, nil, headers)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 400:
		return errors.New("missing Entry-Text header from request")
	case 403:
		defer response.Body.Close()
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("illegal entry: %s", string(errorMessage))
	case 404:
		return errors.New("room \"" + roomName + "\" doesn't exist or no games have been started")
	default:
		return errors.New("something went really wrong :(")
	}
}

// StartGame triggers a game in the room with the provided name.
// Returns error if room doesn't exist, a game is already running or the user doesn't have the required permissions.
func (client *SBClient) StartGame(timeLimit, maxLength, entriesCount int) error {
	if timeLimit < 0 {
		return errors.New("cannost start game: negative time limit value")
	}
	if maxLength < 0 {
		return errors.New("cannost start game: negative max length value")
	}
	if entriesCount < 0 {
		return errors.New("cannost start game: negative entires value")
	}
	roomName := client.config.Room
	if client.config.Room != roomName {
		return errors.New("cannot start game: requires user to be joined in the room")
	}
	headers := make(map[string]string)
	headers["Time-Limit"] = fmt.Sprint(timeLimit)
	headers["Max-Length"] = fmt.Sprint(maxLength)
	headers["Entries-Count"] = fmt.Sprint(entriesCount)
	response, err := client.call(http.MethodPost, "/manage-games/"+roomName, nil, headers)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		return errors.New("cannot start game: requires admin access")
	case 404:
		return errors.New("room \"" + roomName + "\" doesn't exist")
	case 409:
		return errors.New("a game is already running in \"" + roomName + "\"")
	default:
		return errors.New("something went really wrong :(")
	}
}

// EndGame ends a running game. Once called it will set the game's remaining entries to the provided number and the game will effectively end after the count is reached.
// Returns error if room doesn't exist, no game is running or the user doesn't have the required permissions.
func (client *SBClient) EndGame(entriesCount int) error {
	roomName := client.config.Room
	if entriesCount < 0 {
		return errors.New("cannost start game: negative entires value")
	}
	if client.config.Room != roomName {
		return errors.New("cannot end game: requires user to be joined in the room")
	}
	headers := make(map[string]string)
	headers["Entries-Count"] = fmt.Sprint(entriesCount)
	response, err := client.call(http.MethodDelete, "/manage-games/"+roomName, nil, headers)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 202:
		return nil
	case 403:
		return errors.New("cannot end game: requires admin access")
	case 404:
		return errors.New("room \"" + roomName + "\" doesn't exist")
	case 409:
		return errors.New("there is no running game in \"" + roomName + "\"")
	default:
		return errors.New("something went really wrong :(")
	}
}
