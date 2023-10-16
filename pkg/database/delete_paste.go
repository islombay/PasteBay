package database

func (db *Database) DeletePaste(id int) (string, error) {
	var blob_src string
	err := db.db.QueryRow("SELECT blob_src FROM pastes WHERE id = $1", id).Scan(&blob_src)
	if err != nil {
		return "", err
	}

	tx, err := db.db.Begin()
	if err != nil {
		return "", err
	}

	sql := "DELETE FROM pastes WHERE id = $1"
	_, err = tx.Exec(sql, id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec("DELETE FROM hashConnector WHERE pasteID = $1", id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return "", err
	}
	return blob_src, nil
}
