package model

import "time"

type PasswordChangeResponse struct {
	DateCreated time.Time `json:"date_created"`
	Success     bool      `json:"success"`
}
