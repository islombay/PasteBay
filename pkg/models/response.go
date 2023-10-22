package models

import "time"

type ResponseAddPaste struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Alias   string `json:"alias"`
}

type ResponseGetPasteAuthor struct {
	LastLogin time.Time `json:"last_login"`
	Username  string    `json:"username"`
}

type ResponseGetPaste struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Author     interface{} `json:"author,omitempty"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	ViewsCount uint        `json:"viewed"`
}

type ResponseRegister struct {
	Message bool `json:"ok"`
}

type ResponseLogin struct {
	Token   string `json:"token"`
	Message bool   `json:"ok"`
}

type ResponseDeletePaste struct {
	Message bool `json:"ok"`
}
