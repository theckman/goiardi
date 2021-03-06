/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package user

import (
	"database/sql"
	"log"

	"github.com/theckman/goiardi/config"
	"github.com/theckman/goiardi/datastore"
)

func checkForUserSQL(dbhandle datastore.Dbhandle, name string) (bool, error) {
	_, err := datastore.CheckForOne(dbhandle, "users", name)
	if err == nil {
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, err
	}
	return false, nil
}

func (u *User) fillUserFromSQL(row datastore.ResRow) error {
	var email sql.NullString
	err := row.Scan(&u.Username, &u.Name, &u.Admin, &u.pubKey, &email, &u.passwd, &u.salt)
	if err != nil {
		return err
	}
	if !email.Valid {
		u.Email = ""
	} else {
		u.Email = email.String
	}
	return nil
}

func getUserSQL(name string) (*User, error) {
	user := new(User)
	var sqlStatement string
	if config.Config.UseMySQL {
		sqlStatement = "select name, displayname, admin, public_key, email, passwd, salt FROM users WHERE name = ?"
	} else if config.Config.UsePostgreSQL {
		sqlStatement = "select name, displayname, admin, public_key, email, passwd, salt FROM goiardi.users WHERE name = $1"
	}
	stmt, err := datastore.Dbh.Prepare(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(name)
	err = user.fillUserFromSQL(row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) deleteSQL() error {
	tx, err := datastore.Dbh.Begin()
	if err != nil {
		return err
	}
	if config.Config.UseMySQL {
		_, err = tx.Exec("DELETE FROM users WHERE name = ?", u.Username)
	} else if config.Config.UsePostgreSQL {
		_, err = tx.Exec("DELETE FROM goiardi.users WHERE name = $1", u.Username)
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func numAdminsSQL() int {
	var numAdmins int
	var sqlStatement string
	if config.Config.UseMySQL {
		sqlStatement = "SELECT count(*) FROM users WHERE admin = 1"
	} else if config.Config.UsePostgreSQL {
		sqlStatement = "SELECT count(*) FROM goiardi.users WHERE admin = TRUE"
	}
	stmt, err := datastore.Dbh.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&numAdmins)
	if err != nil {
		log.Fatal(err)
	}
	return numAdmins
}

func getListSQL() []string {
	var userList []string
	var sqlStatement string
	if config.Config.UseMySQL {
		sqlStatement = "SELECT name FROM users"
	} else if config.Config.UsePostgreSQL {
		sqlStatement = "SELECT name FROM goiardi.users"
	}
	rows, err := datastore.Dbh.Query(sqlStatement)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		rows.Close()
		return userList
	}
	userList = make([]string, 0)
	for rows.Next() {
		var userName string
		err = rows.Scan(&userName)
		if err != nil {
			log.Fatal(err)
		}
		userList = append(userList, userName)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return userList
}

func allUsersSQL() []*User {
	var users []*User
	var sqlStatement string
	if config.Config.UseMySQL {
		sqlStatement = "SELECT name, displayname, admin, public_key, email, passwd, salt FROM users"
	} else if config.Config.UsePostgreSQL {
		sqlStatement = "SELECT name, displayname, admin, public_key, email, passwd, salt FROM goiardi.users"
	}

	stmt, err := datastore.Dbh.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, qerr := stmt.Query()
	if qerr != nil {
		if qerr == sql.ErrNoRows {
			return users
		}
		log.Fatal(qerr)
	}
	for rows.Next() {
		us := new(User)
		err = us.fillUserFromSQL(rows)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, us)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users
}
