package model

import "time"

type PasswordResetResponse struct {
	DateCreated time.Time `json:"date_created"`
	Success     bool      `json:"success"`
}
