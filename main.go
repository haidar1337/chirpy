package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		os.Remove("./database.json")
	}

	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load env file: %w", err)
		return
	}

	port := "8080"
	mux := http.NewServeMux()
	jwtSecret := os.Getenv("JWT_SECRET")
	cfg := apiConfig{
		fileserverHits: 0,
		jwtSecret:      jwtSecret,
	}
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir("./")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleNumberOfRequests)
	mux.HandleFunc("GET /api/reset", cfg.resetCounter)

	mux.HandleFunc("POST /api/chirps", handlePostChirp)
	mux.HandleFunc("GET /api/chirps", fetchChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", fetchChirpByID)

	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("PUT /api/users", cfg.handlerUsersUpdate)

	mux.HandleFunc("POST /api/refresh", cfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", handlerRevoke)

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
