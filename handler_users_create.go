package main

import (
	"encoding/json"
	"net/http"

	"github.com/haidar1337/chirpy/internal/database"
)

func createUser(w http.ResponseWriter, req *http.Request) {
	type parameter struct {
		Email string `json:"email"`
	}
	p := parameter{}
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		respondWithError(w, 500, "Failed to read request body")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to create database")
		return
	}

	user, err := db.CreateUser(p.Email)
	if err != nil {
		respondWithError(w, 500, "Failed to create user")
		return
	}

	respondWithJSON(w, 201, user)
}