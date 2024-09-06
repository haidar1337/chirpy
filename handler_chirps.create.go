package main

import (
	"encoding/json"
	"net/http"

	"github.com/haidar1337/chirpy/internal/database"
)

func handlePostChirp(w http.ResponseWriter, req *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to create database")
		return
	}

	type parameter struct {
		Body string `json:"body"`
	}
	p := parameter{}
	err = json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		respondWithError(w, 500, "Failed to unmarshal request body")
		return
	}

	chirp, err := db.CreateChirp(p.Body)
	if err != nil && err.Error() == "Chirp is too long" {
		respondWithError(w, 400, err.Error())
		return
	} else if err != nil {
		respondWithError(w, 500, "Failed to create chirp")
		return
	}

	respondWithJSON(w, 201, chirp)
}