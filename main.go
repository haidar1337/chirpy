package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func main() {
	port := "8080"
	mux := http.NewServeMux()
	cfg := apiConfig{
		fileserverHits: 0,
	}
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir("./")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleNumberOfRequests)
	mux.HandleFunc("GET /api/reset", cfg.resetCounter)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server started on port: %s", port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK) // 200
	w.Write([]byte("OK"))
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "No request body")
		return
	}
	
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type returnVal struct {
		Cleaned_body string `json:"cleaned_body"`
	}
	const replacement = "****"
	out := returnVal{
		Cleaned_body: params.Body,
	}
	splitted := strings.Split(params.Body, " ")
	for i := 0; i < len(splitted); i++ {
		word := splitted[i]
		lowered := strings.ToLower(word)
		if lowered == "kerfuffle" || lowered == "fornax" || lowered == "sharbert" {
			out.Cleaned_body = strings.Replace(out.Cleaned_body, word, replacement, -1)
		}
	}
	
	respondWithJSON(w, 200, out)
}

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