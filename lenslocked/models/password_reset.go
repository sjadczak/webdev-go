package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/sjadczak/webdev-go/lenslocked/rand"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when creating a new Password Reset Token. When looking up a Token
	// this will be left empty, as we only store the hash of the reset token in the
	// database.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// BytesPerToken is used to determin how many bytes to use when generating each
	// PasswordReset token. If this value is not set or is less than the MinBytesPerToken
	// const it will be ignored and MinBytesPerToken will be used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset token is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	// Verify we have a valid email address
	email = strings.ToLower(email)
	var userID int
	row := service.DB.QueryRow(`
	SELECT id
	FROM users
	WHERE email=$1;`, email)
	err := row.Scan(&userID)
	if err != nil {
		// TODO: consider returning a specific error when the users does not exist.
		return nil, fmt.Errorf("create: %w", err)
	}

	// Build PasswordResetToken
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	tokenHash := service.hash(token)

	// Set duration
	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}

	// Add PasswordResetToken into DB
	row = service.DB.QueryRow(`
	INSERT INTO password_resets (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash=$2, expires_at=$3
			WHERE password_resets.user_id=$1
	RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create:  %w", err)
	}

	return &pwReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	// Validate Token & check isn't expired
	tokenHash := service.hash(token)
	var user User
	var pwReset PasswordReset
	row := service.DB.QueryRow(`
	SELECT
		password_resets.id, password_resets.expires_at,
		users.id, users.email, users.password_hash
	FROM password_resets
		JOIN users ON password_resets.user_id = users.id
	WHERE password_resets.token_hash = $1;`, tokenHash)
	err := row.Scan(
		&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("consume: token expired - %v", token)
	}

	// Delete PasswordResetToken
	err = service.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &user, nil
}

func (service *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (service *PasswordResetService) delete(id int) error {
	_, err := service.DB.Exec(`
	DELETE FROM password_resets
	WHERE id=$1;`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
