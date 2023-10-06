package models

type Author struct {
	Id       int    `json:"mo_id"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
