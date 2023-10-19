package database

import (
	"PasteBay/pkg/models"
)

func (db *Database) GetUserByUsername(user string) (models.UserModel, error) {
	sql := "SELECT * FROM users WHERE username = $1"
	var res models.UserModel
	err := db.db.Get(&res, sql, user)
	if err != nil {
		return res, err
	}
	return res, nil
}
