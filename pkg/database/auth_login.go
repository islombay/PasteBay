package database

import (
	"PasteBay/pkg/utils/auth"
	"errors"
)

const (
	DBNotFound = "not_founc"
)

func (db *Database) LoginValid(user, pwd string) (bool, error) {
	pwd = auth.GenerateHash(pwd)
	sql := "SELECT COUNT(*) FROM users WHERE username = $1 AND pwd_hash = $2"

	var res int
	err := db.db.QueryRow(sql, user, pwd).Scan(&res)
	if err != nil {
		return false, err
	}
	if res != 1 {
		sql = "SELECT COUNT(*) FROM users WHERE username = $1"
		err = db.db.QueryRow(sql, user).Scan(&res)
		if err != nil {
			return false, err
		}
		if res == 0 {
			return false, errors.New(DBNotFound)
		}
		return false, nil
	}
	return true, nil
}
