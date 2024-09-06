package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
	"sync"
)

type DB struct {
	path string
	mux *sync.RWMutex
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
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

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
	}
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}


func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	f, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStructure := DBStructure{}
	err = json.Unmarshal(f, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(&dbStructure)
	if err != nil {
		return err
	}

	db.mux.RLock()
	defer db.mux.RUnlock()
	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}

	return nil
}