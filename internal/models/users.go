package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// to hold a specific user
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// to interact with the database table
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	// 12 is a reasonable minimum (cost)
	// returns a 60-char long password hash.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	q := `INSERT INTO users (name, email, hashed_password, created_at)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(q, name, email, string(hashedPassword))
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		// for all other errors, return them as is.
		return err
	}

	return nil
}

// if the usere is not registered, or the credentials don't match, return `ErrInvalidCredentials`
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var (
		id             int
		hashedPassword []byte
	)

	q := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := m.DB.QueryRow(q, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// if the credentials match, return the user ID.
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	q := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)` // 0/1

	err := m.DB.QueryRow(q, id).Scan(&exists)
	return exists, err
}
