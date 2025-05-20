package data

import (
	"database/sql"
	"errors"
	"time"
)

// User represents a user in the system
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Do not expose password hash in JSON responses
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AddUser adds a new user to the database
func AddUser(username, email, passwordHash string) (User, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
	stmt, err := DB.Prepare(query)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(username, email, passwordHash)
	if err != nil {
		// Check for unique constraint violation (specific to SQLite)
		if err.Error() == "UNIQUE constraint failed: users.username" || err.Error() == "UNIQUE constraint failed: users.email" {
			return User{}, NewUserExistsError("username or email already exists")
		}
		return User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return User{}, err
	}

	return GetUserByID(id)
}

// GetUserByUsername retrieves a user by their username from the database
func GetUserByUsername(username string) (User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = ?`
	row := DB.QueryRow(query, username)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}
	return u, nil
}

// GetUserByEmail retrieves a user by their email from the database
func GetUserByEmail(email string) (User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = ?`
	row := DB.QueryRow(query, email)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}
	return u, nil
}

// GetUserByID retrieves a user by their ID from the database
func GetUserByID(id int64) (User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = ?`
	row := DB.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}
	return u, nil
}

// UserExistsError is an error type for when a user already exists.
type UserExistsError struct {
	Message string
}

func (e *UserExistsError) Error() string {
	return e.Message
}

func NewUserExistsError(message string) *UserExistsError {
	return &UserExistsError{Message: message}
}

