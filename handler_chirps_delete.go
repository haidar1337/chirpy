package main

import (
	"net/http"
	"strconv"

	"github.com/haidar1337/chirpy/internal/auth"
	"github.com/haidar1337/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 400, "Couldn't find bearer token in authorization header")
		return
	}

	reqArgument := req.PathValue("chirpId")
	if reqArgument == "" {
		respondWithError(w, 400, "please provide a chirp id")
		return
	}
	chirpId, err := strconv.Atoi(reqArgument)
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 400, "Couldn't validate JWT token")
		return
	}

	authorId, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, 500, "Failed to convert user id")
	}
	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to load database")
		return
	}

	err = db.DeleteChirp(chirpId, authorId)
	if err != nil {
		respondWithError(w, 403, "You are not the author of this chirp")
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
