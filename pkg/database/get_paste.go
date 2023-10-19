package database

import (
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/logger/sl"
)

func (db *Database) GetPaste(alias string) (models.PasteModel, error) {
	var res models.PasteModel
	sql := `SELECT pasteID FROM hashConnector WHERE pasteHash = $1`
	var pasteID string

	err := db.db.QueryRow(sql, alias).Scan(&pasteID)
	if err != nil {
		//db.Log.Error("[DB] - could not find the paste alias", sl.Err(err))
		return res, err
	}

	err = db.db.Get(&res, "SELECT * FROM pastes WHERE id = $1", pasteID)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (db *Database) ViewIncreasePaste(id int) error {
	var current int
	sql := "SELECT views_counted FROM pastes WHERE id = $1"

	err := db.db.QueryRow(sql, id).Scan(&current)
	if err != nil {
		db.Log.Error("[DB] - error during getting view count", sl.Err(err))
		return err
	}
	_, err = db.db.Exec("UPDATE pastes SET views_counted = $1 WHERE id = $2", current+1, id)
	if err != nil {
		db.Log.Error("[DB] - error during increasing view count", sl.Err(err))
		return err
	}
	return nil
}
