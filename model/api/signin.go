package model

import "time"

type SignInRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type SignInResponse struct {
	DateCreated time.Time `json:"date_created"`
	Success     bool      `json:"success"`
	Name        string    `json:"name"`
}
