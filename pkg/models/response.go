package models

import "time"

type ResponsePasteView struct {
	ID        interface{}             `json:"id"`
	CreatedAt time.Time               `json:"created_at"`
	Author    ResponseAuthorPasteView `json:"author"`

	Title      string `json:"title"`
	PasteType  uint8  `json:"paste_type"` // 2 - private 1 - public
	Content    string `json:"content"`
	ViewsCount uint   `json:"views_count"`
}

// this is only for ResponsePasteView
type ResponseAuthorPasteView struct {
	Username string
}
