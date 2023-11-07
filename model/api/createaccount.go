package model

import "time"

type CreateAccountResponse struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"date_created"`
	Success     bool      `json:"success"`
	Username    string    `json:"username"`
}
type CreateAccountDatabase struct {
	Username     string `bson:"username"`
	PasswordHash string `bson:"password"`
	Email        string `bson:"email"`
	Name         string `bson:"name"`
}
