package database

import (
	"errors"
	"sort"
	"strings"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	if len(body) > 140 {
		return Chirp{}, errors.New("Chirp is too long")
	}

	const replacement = "****"
	splitted := strings.Split(body, " ")
	for i := 0; i < len(splitted); i++ {
		word := splitted[i]
		lowered := strings.ToLower(word)
		if lowered == "kerfuffle" || lowered == "fornax" || lowered == "sharbert" {
			body = strings.Replace(body, word, replacement, -1)
		}
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp := Chirp{
		ID: len(dbStructure.Chirps) + 1,
		Body: body,
	}
	dbStructure.Chirps[len(dbStructure.Chirps) + 1] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err  := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, Chirp{ID: chirp.ID, Body: chirp.Body})
	}

	sort.Slice(chirps, func(i int, j int) bool {
		if chirps[i].ID < chirps[j].ID {
			return true
		}
		return false
	})

	return chirps, nil
}

func (db *DB) GetChirp(Id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[Id]
	if !ok {
		return Chirp{}, errors.New("Chirp not found")
	}

	return chirp, nil
}