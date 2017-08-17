package model

import (
	"fmt"
	"time"

	"app/shared/database"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ID        uint32    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	StatusID  uint8     `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   uint8     `db:"deleted"`
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	var err error

	result := User{}

	err = database.SQL.Get(&result, "SELECT id, password, status_id, username FROM user WHERE email = ? LIMIT 1", email)
	return result, standardizeError(err)
}

// UserID returns the user id
func (u *User) UserID() string {
	return fmt.Sprintf("%v", u.ID)
}

// UserCreate creates user
func UserCreate(username, email, password string) error {
	_, err := database.SQL.Exec("INSERT INTO user (username, email, password) VALUES (?,?,?)", username, email, password)
	return standardizeError(err)
}
