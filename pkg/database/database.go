package database

import (
	"PasteBay/pkg/utils/logger/sl"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

type Database struct {
	db *sql.DB
}

type DatabaseLoad struct {
	Host, Port, DBName, SSLMode, User, Password string
	Log                                         *slog.Logger
}

func InitDatabase(cfg DatabaseLoad) *Database {
	db, err := sql.Open("postgres",
		fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
		))
	if err != nil {
		cfg.Log.Error("could not open database", sl.Err(err))
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		cfg.Log.Error("could not ping database", sl.Err(err))
		os.Exit(1)
	}

	InitTables(db, cfg.Log)

	return &Database{
		db: db,
	}
}

func InitTables(db *sql.DB, log *slog.Logger) {
	createUsersTable := `CREATE TABLE IF NOT EXISTS users(
    	id SERIAL PRIMARY KEY,
    	created_at DATE,
    	last_login DATE,
    	username VARCHAR(200),
    	email_addr VARCHAR(200),
    	pwd_hash VARCHAR(255)
	)
    `
	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Error("could not create table users", sl.Err(err))
		os.Exit(1)
	}

	createPasteTable := `CREATE TABLE IF NOT EXISTS pastes(
    	id SERIAL PRIMARY KEY,
    	created_at DATE,
    	updated_at DATE,
    	author INT DEFAULT -1,
    	
    	title VARCHAR(200),
    	is_private BOOLEAN DEFAULT false,
    	expire_time DATE,
    	views_limit INT DEFAULT -1,
    	
    	blob_src VARCHAR(255),
    	access_password VARCHAR(255)    	
	)`
	_, err = db.Exec(createPasteTable)
	if err != nil {
		log.Error("could not create table pastes", sl.Err(err))
		os.Exit(1)
	}

	createHashConnectorTable := `CREATE TABLE IF NOT EXISTS hashConnector(
    	pasteHash VARCHAR(255) UNIQUE NOT NULL,
    	pasteID INT NOT NULL
	)`
	_, err = db.Exec(createHashConnectorTable)
	if err != nil {
		log.Error("could not create table hashConnector", sl.Err(err))
		os.Exit(1)
	}
}

func (db *Database) Close() error {
	return db.db.Close()
}
