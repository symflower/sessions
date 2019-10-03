package model

import (
	"database/sql"
)

// User holds a user.
type User struct {
	Mail     string
	Password string
}

// UserInit initializes the users table.
func UserInit(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE users(mail TEXT, password TEXT)")

	return err
}

// UserByMail returns a user identified by its mail address, or nil if the user does not exist.
func UserByMail(db *sql.DB, mail string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT mail, password FROM users WHERE mail = '"+mail+"'").Scan(&user.Mail, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

// UserAdd adds a user to the database.
func UserAdd(db *sql.DB, mail string, password string) error {
	_, err := db.Exec("INSERT INTO users(mail, password) VALUES('" + mail + "', '" + password + "')")

	return err
}

// UserLogin returns if a password belongs to a user.
func UserLogin(db *sql.DB, mail string, password string) (bool, error) {
	err := db.QueryRow("SELECT mail FROM users WHERE mail = '" + mail + "' AND password = '" + password + "'").Scan(&mail)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
