package database

import (
	"errors"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	dbstructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbstructure.Users) + 1
	newUser := User{
		ID:    id,
		Email: email,
	}
	dbstructure.Users[id] = newUser
	err = db.writeDB(dbstructure)
	if err != nil {
		return User{}, nil
	}
	return newUser, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("COULD NOT LOAD DB")
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("COULD NOT FIND USER")
	}
	return user, nil
}
