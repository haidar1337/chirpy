package main

import (
	"net/http"

	"github.com/haidar1337/chirpy/internal/auth"
	"github.com/haidar1337/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, req *http.Request) {
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
	user, err := db.FindRefreshToken(bearerToken)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 500, "Failed to create JWT token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, response{
		Token: jwtToken,
	})
}
