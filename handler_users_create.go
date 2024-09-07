package main

import (
	"encoding/json"
	"net/http"

	"github.com/haidar1337/chirpy/internal/auth"
	"github.com/haidar1337/chirpy/internal/database"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "failed to load database server")
		return
	}
	user, err := db.CreateUser(params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
