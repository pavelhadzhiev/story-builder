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

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// SBDatabase represents the database layer for the story builder server
type SBDatabase struct {
	database *sql.DB // a database variable to close on program exit
	username string
	password string
}

// NewSBDatabase returns a pointer to a new instance of SBDatabase, using the provided credentials.
func NewSBDatabase(username string, password string) *SBDatabase {
	return &SBDatabase{
		username: username,
		password: password,
	}
}

// InitializeDB connects to a local MySQL server using the configured user, creates a database named "storybuilder" and creates all necessary for the story builder API tables inside. Recourses are created only if they do not exist.
func (sbdb *SBDatabase) InitializeDB() error {
	var config = mysql.Config{
		User:   sbdb.username,
		Passwd: sbdb.password,
	}
	var err error
	sbdb.database, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}

	_, err = sbdb.database.Exec("create database if not exists storybuilder")
	if err != nil {
		return err
	}

	_, err = sbdb.database.Exec("use storybuilder")
	if err != nil {
		log.Printf("%q\n", err)
		return err
	}

	_, err = sbdb.database.Exec("create table if not exists users (username varchar(255) not null primary key, password varchar(255) not null)")
	if err != nil {
		return err
	}

	return nil
}

// CloseDB shuts down the connection to the database.
func (sbdb *SBDatabase) CloseDB() {
	sbdb.database.Close()
}
