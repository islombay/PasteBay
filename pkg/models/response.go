package models

import "time"

type ResponseAddPaste struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Alias   string `json:"alias"`
}

type ResponseGetPasteAuthor struct {
}

type ResponseGetPaste struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Author     string `json:"author"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ViewsCount uint   `json:"viewed"`
}

type ResponseRegister struct {
	Message bool `json:"ok"`
}
