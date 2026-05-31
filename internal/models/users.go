package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	// 12 is a reasonable minimum
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	q := `INSERT INTO users (name, email, hashed_password, created_at)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(q, name, email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
