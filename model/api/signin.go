package model

import "time"

type SignInRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type SignInResponse struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"date_created"`
	Success     bool      `json:"success"`
	Username    string    `json:"username"`
}
