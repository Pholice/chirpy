package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type ResponseUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string, password string) (ResponseUser, error) {
	dbstructure, err := db.loadDB()
	if err != nil {
		return ResponseUser{}, err
	}

	id := len(dbstructure.Users) + 1
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return ResponseUser{}, errors.New("COULD NOT ENCRYPT PASSWORD")
	}

	for _, user := range dbstructure.Users {
		if user.Email == email {
			return ResponseUser{}, errors.New("USER ALREADY EXISTS")
		}
	}

	newUser := User{
		ID:       id,
		Email:    email,
		Password: hashedPW,
	}
	dbstructure.Users[id] = newUser
	err = db.writeDB(dbstructure)
	if err != nil {
		return ResponseUser{}, nil
	}

	return ResponseUser{ID: newUser.ID, Email: newUser.Email}, nil
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

func (db *DB) GetUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("COULD NOT LOAD DB")
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	// user, ok := dbStructure.Users[id]
	return User{}, errors.New("COULD NOT FIND USER")
}
