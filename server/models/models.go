package models

type RegisterInfos struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email_address"`
	Birthdate string `json:"birthdate"`
}

type LoginInfos struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
