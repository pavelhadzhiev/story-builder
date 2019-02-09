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

// var user = "admin" // default username
// var passwd = "Abcd1234" // default password

type SBDatabase struct {
	database *sql.DB // a database variable to close on program exit
	username string
	password string
}

func NewSBDatabase() *SBDatabase {
	return &SBDatabase{
		username: "admin",
		password: "Abcd1234",
	}
}

func (sbdb *SBDatabase) InitializeDB() error {
	sbdb.username = "admin"    // default username
	sbdb.password = "Abcd1234" // default password

	var config = mysql.Config{
		User:   sbdb.username,
		Passwd: sbdb.password,
	}
	var err error
	sbdb.database, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
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

func (sbdb *SBDatabase) SetUsername(username string) {
	sbdb.username = username
}

func (sbdb *SBDatabase) SetPassword(password string) {
	sbdb.password = password
}

func (sbdb *SBDatabase) CloseDB() {
	sbdb.database.Close()
}
