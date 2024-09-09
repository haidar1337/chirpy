package database

import (
	"errors"
	"sort"
	"strings"
)

type Chirp struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author_ID int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {

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
		ID:        len(dbStructure.Chirps) + 1,
		Body:      body,
		Author_ID: authorID,
	}
	dbStructure.Chirps[len(dbStructure.Chirps)+1] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps(sorting string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, Chirp{ID: chirp.ID, Body: chirp.Body, Author_ID: chirp.Author_ID})
	}

	if sorting == "asc" {
		sort.Slice(chirps, func(i int, j int) bool {
			if chirps[i].ID < chirps[j].ID {
				return true
			}
			return false
		})

	} else if sorting == "desc" {
		sort.Slice(chirps, func(i int, j int) bool {
			if chirps[i].ID >= chirps[j].ID {
				return true
			}
			return false
		})
	}

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

func (db *DB) DeleteChirp(id, authorId int) error {
	structure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, err := db.ValidateChirp(id, authorId)
	if err != nil {
		return err
	}

	delete(structure.Chirps, chirp.ID)
	err = db.writeDB(structure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ValidateChirp(id, authorId int) (Chirp, error) {
	chirp, err := db.GetChirp(id)
	if err != nil {
		return Chirp{}, err
	}
	if chirp.Author_ID != authorId {
		return Chirp{}, errors.New("Unauthorized")
	}

	return chirp, nil
}

func (db *DB) GetChirpsByAuthorID(authorId int, sorting string) ([]Chirp, error) {
	structure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	out := make([]Chirp, 0)
	for _, chirp := range structure.Chirps {
		if chirp.Author_ID == authorId {
			out = append(out, chirp)
		}
	}

	if sorting == "asc" {
		sort.Slice(out, func(i int, j int) bool {
			if out[i].ID < out[j].ID {
				return true
			}
			return false
		})
	} else if sorting == "desc" {
		sort.Slice(out, func(i int, j int) bool {
			if out[i].ID >= out[j].ID {
				return true
			}
			return false
		})
	}

	return out, nil
}
