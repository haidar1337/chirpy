package main

import (
	"encoding/json"
	"net/http"
)


func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	type errorVal struct{
		Error string `json:"error"`
	}

	respondWithJSON(w, statusCode, errorVal{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("content-type", "application/json")
	data, err := json.Marshal(&payload)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}