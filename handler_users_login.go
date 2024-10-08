package main

import (
	"encoding/json"
	"net/http"

	"github.com/haidar1337/chirpy/internal/auth"
	"github.com/haidar1337/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		database.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "failed to load database server")
		return
	}
	user, err := db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}
	hexToken := user.RefreshToken
	if hexToken.Token == "" {
		hexToken, err = auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, 500, "failed to create a refresh token")
			return
		}
		err := db.UpdateRefreshToken(user.ID, hexToken.Token)
		if err != nil {
			respondWithError(w, 500, "failed to save refresh token")
			return
		}
	}

	respondWithJSON(w, http.StatusOK, response{
		User: database.User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: hexToken.Token,
	})
}
