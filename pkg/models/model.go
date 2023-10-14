package models

import (
	"time"
)

type PasteModel struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    uint

	Title          string
	PasteType      uint8 // 2 - private 1 - public
	BlobSrc        string
	ExpireTime     time.Time
	ViewsLimit     uint
	ViewsCount     uint
	AccessPassword string
}

type UserModel struct {
	ID           uint
	CreatedAt    time.Time
	LastLogin    time.Time
	Username     string
	EmailAddr    string
	PasswordHash string
}
