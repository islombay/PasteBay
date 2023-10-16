package models

import (
	"time"
)

type PasteModel struct {
	ID        int
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Author    int64

	Title          string
	BlobSrc        string    `db:"blob_src"`
	ExpireTime     time.Time `db:"expire_time"`
	ViewsLimit     int64     `db:"views_limit"`
	ViewsCount     uint      `db:"views_counted"`
	AccessPassword string    `db:"access_password"`
}

type UserModel struct {
	ID           uint
	CreatedAt    time.Time
	LastLogin    time.Time
	Username     string
	EmailAddr    string
	PasswordHash string
}
