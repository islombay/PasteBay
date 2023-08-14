package models

import "time"

type RequestAddPaste struct {
	Title          string    `json:"title"`
	PasteType      uint8     `json:"type"`
	Content        string    `json:"content" binding:"required"`
	ExpireTime     time.Time `json:"expire_time"`
	ViewsLimit     uint      `json:"view_limit"`
	AccessPassword string    `json:"password"`
}

type RequestGetPaste struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}
