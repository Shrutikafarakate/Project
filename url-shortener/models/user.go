package models

import (
	"database/sql"
	"errors"
	"url-shortener/db"
	"url-shortener/utils"
)

type User struct {
	ID       int
	Username string
	Password string
}

func CreateUser(username, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	return err
}

func AuthenticateUser(username, password string) (*User, error) {
	row := db.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", username)

	var id int
	var hashedPassword string
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(password, hashedPassword) {
		return nil, errors.New("incorrect password")
	}

	return &User{ID: id, Username: username}, nil
}
