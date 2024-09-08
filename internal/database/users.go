package database

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	ExpirationDate jwt.NumericDate `json:"expiry_date"`
	Token          string          `json:"refresh_token"`
}

type User struct {
	ID           int          `json:"id"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	RefreshToken RefreshToken `json:"refresh_token"`
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	for _, u := range dbStructure.Users {
		users = append(users, User{ID: u.ID, Email: u.Email})
	}

	return users, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("User not found")
}

func (db *DB) GetUserByID(id int) (User, error) {
	structure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range structure.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return User{}, errors.New("User not found")
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("User not found")
	}

	user.Email = email
	user.Password = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateRefreshToken(id int, token string) error {
	structure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		return err
	}

	updated := User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		RefreshToken: RefreshToken{
			ExpirationDate: *jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 1440)),
			Token:          token,
		},
	}
	structure.Users[user.ID] = updated
	err = db.writeDB(structure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserRefreshToken(id int) (string, error) {
	user, err := db.GetUserByID(id)
	if err != nil {
		return "", err
	}

	return user.RefreshToken.Token, nil
}

func (db *DB) FindRefreshToken(token string) (User, error) {
	structure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range structure.Users {
		if user.RefreshToken.Token == token && user.RefreshToken.ExpirationDate.Sub(time.Now().UTC()) > 0 {
			return user, nil
		}
	}

	return User{}, errors.New("Refresh token not found or expired")
}

func (db *DB) RevokeToken(token string) error {
	user, err := db.FindRefreshToken(token)
	if err != nil {
		return err
	}

	structure, err := db.loadDB()
	if err != nil {
		return err
	}
	updated := User{
		ID:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: RefreshToken{},
	}
	structure.Users[user.ID] = updated
	err = db.writeDB(structure)
	if err != nil {
		return err
	}

	return nil
}
