package models

type RequestAddPaste struct {
	Title          string `json:"title"`
	Content        string `json:"content" binding:"required"`
	ExpireTime     int64  `json:"expire_time"` // get the amount of seconds
	ViewsLimit     int    `json:"view_limit"`
	AccessPassword string `json:"password"`
}

type RequestGetPaste struct {
	Password string `json:"password" binding:"required"`
}

type RequestRegister struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RequestLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
