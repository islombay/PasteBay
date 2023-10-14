package models

type ResponseAddPaste struct {
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private"`
	Content   string `json:"content"`
	Alias     string `json:"alias"`
}
