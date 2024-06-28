package database

import (
	"errors"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, author int) (Chirp, error) {
	dbstructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbstructure.Chirps) + 1
	newChirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: author,
	}
	dbstructure.Chirps[id] = newChirp
	err = db.writeDB(dbstructure)
	if err != nil {
		return Chirp{}, nil
	}
	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, errors.New("COULD NOT LOAD DB")
	}
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("COULD NOT FIND CHIRP")
	}
	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return errors.New("COULD NOT LOAD DB")
	}
	delete(dbStructure.Chirps, id)
	err = db.writeDB(dbStructure)
	if err != nil {
		return errors.New("COULD SAVE DB")
	}
	return nil
}
