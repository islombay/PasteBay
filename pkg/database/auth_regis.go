package database

import (
	"PasteBay/pkg/utils/auth"
	"errors"
	"time"
)

func (db *Database) CheckUsername(username string) (bool, error) {
	sql := "SELECT COUNT(*) FROM users WHERE username = $1"
	var count int

	err := db.db.QueryRow(sql, username).Scan(&count)
	if err != nil {
		return true, err
	}
	if count > 0 {
		return false, err
	}
	return true, nil
}

func (db *Database) CreateUser(email, user, password string) error {
	sql := `INSERT INTO users (
                   created_at,
                   last_login,
                   username,
                   email_addr,
                   pwd_hash
	) VALUES ($1, $2, $3, $4, $5)`

	if email == "" || user == "" || password == "" {
		return errors.New("empty values")
	}

	now := time.Now()
	_, err := db.db.Exec(sql, now, now, user, email, auth.GenerateHash(password))
	if err != nil {
		return err
	}
	return nil
}
