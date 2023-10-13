package models

type RequestAddPaste struct {
	Title                  string `json:"title"`
	PasteType              uint8  `json:"type"`
	Content                string `json:"content" binding:"required"`
	ExpireTimeMilliseconds int64  `json:"expire_time"` // get the amount of milleseconds
	ViewsLimit             uint   `json:"view_limit"`
	AccessPassword         string `json:"password"`
}

type RequestGetPaste struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

type UserRegister struct {
	Username  string `json:"username" binding:"required"`
	EmailAddr string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
