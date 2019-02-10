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

package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadUsername is a utility function, used to read usernames from the terminal.
func ReadUsername() (string, error) {
	fmt.Print("Username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(username, "\n"), nil
}

// ReadPassword is a utility function, used to read passwords securely from the terminal.
func ReadPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword((int)(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(password), nil
}
