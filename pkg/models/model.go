package models

import (
	"time"
)

type PasteModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    uint

	Title          string
	PasteType      uint8 // 2 - private 1 - public
	BlobSrc        string
	ExpireTime     time.Time
	ViewsLimit     uint
	AccessPassword string
}

type UserModel struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	LastLogin    time.Time
	Username     string
	EmailAddr    string
	PasswordHash string
}
