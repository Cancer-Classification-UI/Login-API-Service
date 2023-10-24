package model

import "time"

type CreateAccountResponse struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"date_created"`
	Success     bool      `json:"success"`
	Username    string    `json:"username"`
}
