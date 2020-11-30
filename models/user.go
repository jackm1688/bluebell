package models

type User struct {
	UserId   int64  `json:"user_id,string" db:"user_id"`
	Username string `json:"username" db:"username" binding:"required,lt=6,lte=30"`
	Password string `json:"password" db:"password" binding:"required,lte=8"`
	Token    string `json:"token"`
}
