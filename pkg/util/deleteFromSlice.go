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

// DeleteFromSlice is a utility function to remove an object from a slice at the provided index.
func DeleteFromSlice(slice []string, index int) []string {
	if len(slice) > index+1 {
		return append(slice[:index], slice[index+1:]...)
	} else {
		return slice[:index]
	}
}
