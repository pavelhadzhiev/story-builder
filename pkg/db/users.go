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

package db

import "errors"

const GET_USER_BY_USERNAME = "select * from users where username = ?"

// UserExists returns true if the provided username is already taken according to the server database.
func (sbdb *SBDatabase) UserExists(username string) (bool, error) {
	stmt, err := sbdb.database.Prepare(GET_USER_BY_USERNAME)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Err(); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// LoginUser returns true if the provided user exists and the password matches the one that is saved for that username in the server database.
func (sbdb *SBDatabase) LoginUser(username, password string) error {
	stmt, err := sbdb.database.Prepare(GET_USER_BY_USERNAME)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		return err
	}
	defer rows.Close()
	var user, pass string
	if rows.Next() {
		err := rows.Scan(&user, &pass)
		if err != nil {
			return err
		}
		if pass != password {
			return errors.New("password incorrect")
		}
		return nil
	}
	return errors.New("user not found")
}

// RegisterUser registers a new user to the server with the provided username and password.
func (sbdb *SBDatabase) RegisterUser(username, password string) error {
	_, err := sbdb.database.Exec("insert into users(username, password) values(?, ?)", username, password)
	if err != nil {
		return err
	}
	return nil
}
