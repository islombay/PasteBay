package models

type RequestAddPaste struct {
	Title          string `json:"title"`
	IsPrivate      bool   `json:"is_private"`
	Content        string `json:"content" binding:"required"`
	ExpireTime     int64  `json:"expire_time"` // get the amount of seconds
	ViewsLimit     int    `json:"view_limit"`
	AccessPassword string `json:"password"`
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
