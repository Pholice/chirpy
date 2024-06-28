package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     []byte    `json:"password"`
	ExpiresAt    time.Time `json:"expire"`
	RefreshToken string    `json:"token"`
	Red          bool      `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return User{}, errors.New("COULD NOT ENCRYPT PASSWORD")
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, errors.New("USER ALREADY EXISTS")
		}
	}
	refresh, expires, err := generateRefresh()
	if err != nil {
		return User{}, errors.New("COULD NOT GENERATE REFRESH TOKEN")
	}
	newUser := User{
		ID:           id,
		Email:        email,
		Password:     hashedPW,
		ExpiresAt:    expires,
		RefreshToken: refresh,
		Red:          false,
	}
	dbStructure.Users[id] = newUser

	refreshToken := RefreshToken{
		UserID:    id,
		Token:     refresh,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	dbStructure.RefreshTokens[refresh] = refreshToken
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, nil
	}
	return User{ID: newUser.ID, Email: newUser.Email, RefreshToken: refresh, Red: false}, nil
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

func (db *DB) GetUserEmail(email string) (User, error) {
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

func (db *DB) GetUserID(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("COULD NOT LOAD DB")
	}
	for _, user := range dbStructure.Users {
		if user.ID == id {
			return user, nil
		}
	}
	// user, ok := dbStructure.Users[id]
	return User{}, errors.New("COULD NOT FIND USER")
}

func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("COULD NOT LOAD DB")
	}
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return User{}, errors.New("COULD NOT ENCRYPT PASSWORD")
	}
	updatedUser := User{
		ID:       id,
		Email:    email,
		Password: hashedPW,
	}
	dbStructure.Users[id] = updatedUser
	db.writeDB(dbStructure)
	return updatedUser, nil
}

func (db *DB) UpdateUserType(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("COULD NOT LOAD DB")
	}
	user, err := db.GetUserID(id)
	if err != nil {
		return User{}, errors.New("COULD NOT GET USER ID")
	}
	updatedUser := User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		ExpiresAt:    user.ExpiresAt,
		RefreshToken: user.RefreshToken,
		Red:          true,
	}
	dbStructure.Users[id] = updatedUser
	db.writeDB(dbStructure)
	return updatedUser, nil
}

func generateRefresh() (string, time.Time, error) {
	length := 256
	lengthByte := make([]byte, length)
	_, err := rand.Read(lengthByte)
	if err != nil {
		return "", time.Now(), err
	}
	expires := time.Now().UTC().Add(time.Hour * 1440)
	refreshHex := hex.EncodeToString(lengthByte)
	return refreshHex, expires, nil
}
