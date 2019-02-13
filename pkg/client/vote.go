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

// TriggerVoteKick triggers a democratic vote to kick the player with the provided username from the game in the room with the provided room name.
// Returns error if room doesn't exist, game is not started, the player is not in the game, or another vote is currently ongoing.
func (client *SBClient) TriggerVoteKick(roomName, playerToKick string) error {
	response, err := client.call(http.MethodPost, "/vote/"+roomName+"/"+playerToKick, nil, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 202:
		return nil
	case 404:
		defer response.Body.Close()
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("could not trigger vote: %s", string(errorMessage))
	case 409:
		return errors.New("there is already an ongoing vote")
	default:
		return errors.New("something went really wrong :(")
	}
}

// SubmitVote tells the server that the user agrees with the ongoing vote.
// Returns error if room doesn't exist, game is not started or no vote is currently running.
func (client *SBClient) SubmitVote(roomName string) error {
	response, err := client.call(http.MethodPut, "/vote/"+roomName, nil, nil)
	if err != nil {
		return fmt.Errorf("error during http request: %e", err)
	}
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		return errors.New("cannot vote: user is not part of the game")
	case 404:
		defer response.Body.Close()
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("cannot vote: %s", string(errorMessage))
	case 409:
		return errors.New("cannot vote: user has already voted once")
	default:
		return errors.New("something went really wrong :(")
	}
}
