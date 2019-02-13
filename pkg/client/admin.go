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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// BanPlayer bans the provided player from entering the room with the provided name.
// Returns error if either the room or the player doesn't exist or if the issuer doesn't have admin access for the room.
func (client *SBClient) BanPlayer(roomName, player string) error {
	response, err := client.call(http.MethodDelete, "/admin/ban/"+roomName+"/"+player, nil, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		return errors.New("user does not have permissions to ban in room \"" + roomName + "\"")
	case 404:
		defer response.Body.Close()
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("could not ban player: %s", string(errorMessage))
	case 409:
		return errors.New("player is already banned")
	default:
		return errors.New("something went really wrong :(")
	}
}

// KickPlayer kicks the provided player from the current game in the room with the provided name.
// Returns error if either the room, the game or the player doesn't exist, the player is not in the game, or the issuer doesn't have admin access for the room.
func (client *SBClient) KickPlayer(roomName, player string) error {
	response, err := client.call(http.MethodDelete, "/admin/kick/"+roomName+"/"+player, nil, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		return errors.New("user does not have permissions to kick in room \"" + roomName + "\"")
	case 404:
		defer response.Body.Close()
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("could not kick player: %s", string(errorMessage))
	default:
		return errors.New("something went really wrong :(")
	}
}

// PromoteAdmin gives admin permissions to the provider user in the room with the provided name.
// Returns error if the room or user doesn't exist or the issuer doesn't have admin access.
func (client *SBClient) PromoteAdmin(roomName, user string) error {
	response, err := client.call(http.MethodPost, "/admin/"+roomName+"/"+user, nil, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		return errors.New("user does not have permissions to promote in room \"" + roomName + "\"")
	case 404:
		defer response.Body.Close()
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("could not promote user: %s", string(errorMessage))
	default:
		return errors.New("something went really wrong :(")
	}
}
