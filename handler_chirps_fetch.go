package main

import (
	"net/http"
	"strconv"

	"github.com/haidar1337/chirpy/internal/database"
)

func fetchChirpByID(w http.ResponseWriter, req *http.Request) {
	reqArgument := req.PathValue("chirpId")
	if reqArgument == "" {
		respondWithError(w, 400, "please provide a chirp id")
		return
	}

	id, err := strconv.Atoi(reqArgument)
	if err != nil {
		respondWithError(w, 400, "please provide a valid chirp id")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to create database")
		return
	}

	chirp, err := db.GetChirp(id)
	if err != nil && err.Error() == "Chirp not found" {
		respondWithError(w, 404, err.Error())
		return
	} else if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 200, chirp)
}

func fetchChirps(w http.ResponseWriter, req *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to create database")
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Failed to read chirps")
		return
	}

	respondWithJSON(w, 200, chirps)
}