module github.com/haidar1337/chirpy

go 1.23.0

replace github.com/haidar1337/chirpy/internal/database => ./internal/database

require github.com/haidar1337/chirpy/internal/database v1.2.3

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.27.0 // indirect
)
