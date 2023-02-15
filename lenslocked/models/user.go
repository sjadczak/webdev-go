package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w\n", err)
	}

	user := User{
		Email:        email,
		PasswordHash: string(passwordHash),
	}
	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`, user.Email, user.PasswordHash)

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w\n", err)
	}

	return &user, nil
}
