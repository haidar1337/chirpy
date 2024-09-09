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
	var chirps []database.Chirp
	authorIdQuery := req.URL.Query().Get("author_id")
	sortingQuery := req.URL.Query().Get("sort")
	if authorIdQuery != "" {
		authorId, err := strconv.Atoi(authorIdQuery)
		if err != nil {
			respondWithError(w, 400, "Invalid author id, must be an integer")
			return
		}
		if sortingQuery != "" {
			chirps, err = db.GetChirpsByAuthorID(authorId, sortingQuery)
		} else {
			chirps, err = db.GetChirpsByAuthorID(authorId, "asc")
		}
		if err != nil {
			respondWithError(w, 400, "sorting query must be desc or asc")
			return
		}
	} else {
		if sortingQuery != "" {
			chirps, err = db.GetChirps(sortingQuery)
		} else {
			chirps, err = db.GetChirps("asc")
		}
		if err != nil {
			respondWithError(w, 400, "Sorting query must be asc or desc")
			return
		}
	}

	respondWithJSON(w, 200, chirps)
}
