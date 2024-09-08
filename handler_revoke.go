package main

import (
	"net/http"

	"github.com/haidar1337/chirpy/internal/auth"
	"github.com/haidar1337/chirpy/internal/database"
)

func handlerRevoke(w http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 400, "Couldn't find token in authorization header")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to load database")
		return
	}

	err = db.RevokeToken(bearerToken)
	if err != nil {
		respondWithError(w, 400, "User not found or the token already expired")
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
