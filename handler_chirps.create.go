package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/haidar1337/chirpy/internal/auth"
	"github.com/haidar1337/chirpy/internal/database"
)

func (cfg *apiConfig) handlePostChirp(w http.ResponseWriter, req *http.Request) {
	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to create database")
		return
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Log in to create a chirp")
		return
	}
	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 400, "JWT token expired")
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

	authorId, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, 500, "Couldn't convert user id")
		return
	}
	chirp, err := db.CreateChirp(p.Body, authorId)
	if err != nil && err.Error() == "Chirp is too long" {
		respondWithError(w, 400, err.Error())
		return
	} else if err != nil {
		respondWithError(w, 500, "Failed to create chirp")
		return
	}

	respondWithJSON(w, 201, chirp)
}
