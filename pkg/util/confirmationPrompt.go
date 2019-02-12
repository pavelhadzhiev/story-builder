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
)

// ConfirmationPrompt is a utility function, used to prompt the user for confirmation on potentially dangerous commands.
// Returns true if the user types yes and false if they don't.
func ConfirmationPrompt(action string) bool {
	fmt.Printf("Are you sure you want to %s? (yes/no)\n", action)
	reader := bufio.NewReader(os.Stdin)
	if input, err := reader.ReadString('\n'); err != nil || strings.Compare(strings.ToLower(input), "yes\n") != 0 {
		return false
	}
	return true
}
