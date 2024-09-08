package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/haidar1337/chirpy/internal/auth"
	"github.com/haidar1337/chirpy/internal/database"
)

type Webhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerUpgradeWebhook(w http.ResponseWriter, req *http.Request) {
	webhook, err := handleWebhook(req, cfg.polkaKey)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	if webhook.Event != "user.upgraded" {
		respondWithError(w, 204, "only handling upgrade events")
		return
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Failed to load database")
		return
	}
	user, err := db.GetUserByID(webhook.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "User not found")
		return
	}
	err = db.UpgradeUser(user)
	if err != nil {
		respondWithError(w, 500, "Failed to upgrade user")
		return
	}

	respondWithJSON(w, 204, struct{}{})
}

func handleWebhook(req *http.Request, key string) (Webhook, error) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		return Webhook{}, errors.New("No authorization header")
	}
	if apiKey != key {
		return Webhook{}, errors.New("Wrong API Key")
	}
	type parameters struct {
		Data struct {
			UserID int `json:"user_id"`
		} `json:"data"`
		Event string `json:"event"`
	}
	params := parameters{}
	err = json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		return Webhook{}, err
	}

	return Webhook{
		Event: params.Event,
		Data:  params.Data,
	}, nil
}
