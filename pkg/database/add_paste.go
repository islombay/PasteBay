package database

import (
	"PasteBay/pkg/utils/logger/sl"
	"PasteBay/pkg/utils/random"
	"time"
)

type BodyAddPaste struct {
	Author     int
	Title      string
	ExpireTime int64
	ViewsLimit int
	BlobSrc    string
	Password   string
}

func (db *Database) AddPaste(body BodyAddPaste) (int64, string, error) {
	insertSQL := `INSERT INTO pastes (
                    created_at,
                    updated_at,
                    author,
                    title,
                    expire_time,
                    views_limit,
                    blob_src,
                    access_password
                    ) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	if body.Author == 0 {
		body.Author = -1
	}
	if body.ViewsLimit == 0 {
		body.ViewsLimit = -1
	}

	created_time := time.Now()
	expireTime := created_time.Add(time.Duration(body.ExpireTime) * time.Second)

	tx, err := db.db.Begin()
	if err != nil {
		db.Log.Error("could not start transaction on adding paste")
		return 0, "", err
	}
	_, err = tx.Exec(
		insertSQL,
		created_time,
		created_time,
		body.Author,
		body.Title,
		expireTime,
		body.ViewsLimit,
		body.BlobSrc,
		body.Password,
	)
	if err != nil {
		db.Log.Error("transaction add paste error into table pastes", sl.Err(err))
		if err := tx.Rollback(); err != nil {
			db.Log.Error("transaction rollback error", sl.Err(err))
		}
		return 0, "", err
	}

	query2 := `INSERT INTO hashConnector(
                          pasteHash,
                          pasteID
	) VALUES($1, $2)`

	var insertedID int64

	err = tx.QueryRow("SELECT id FROM pastes WHERE blob_src = $1 AND created_at = $2", body.BlobSrc, created_time).Scan(&insertedID)
	if err != nil {
		db.Log.Error("could not get the id of row", sl.Err(err))
		return 0, "", err
	}

	randomAlias := random.NewRandomString()

	_, err = tx.Exec(query2, randomAlias, insertedID)
	if err != nil {
		db.Log.Error("transaction add paste error into table pastes", sl.Err(err))
		if err := tx.Rollback(); err != nil {
			db.Log.Error("transaction rollback error", sl.Err(err))
		}
		return 0, "", err
	}

	err = tx.Commit()
	if err != nil {
		db.Log.Error("could not commit transaction on add paste", sl.Err(err))
		if err := tx.Rollback(); err != nil {
			db.Log.Error("transaction rollback error", sl.Err(err))
		}
		return 0, "", err
	}
	return insertedID, randomAlias, nil
}
