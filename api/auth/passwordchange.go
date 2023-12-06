package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"ccu/api"
	db "ccu/db"
	mAPI "ccu/model/api"

	"regexp"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo"
)

// GetPasswordChangeVerifyCode godoc
// @Summary      Validates a given reset code
// @Description  Check to see if given reset code is a valid
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        password_hash query string    true "new hashed account password"
// @Param        email         query string    true "email of the user"
// @Param        code          query string    true "password reset code"
// @Success      200  {array}   mAPI.PasswordChangeResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /password-change-verify-code [get]
func GetPasswordChangeVerifyCode(w http.ResponseWriter, r *http.Request) {
	log.Info("In password-change-verify-code handler -------------------------")
	r.ParseForm()
	if r.Method != http.MethodGet {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	code := r.Form.Get("code")

	if code == "" {
		api.Respond(w, "Invalid Reset Code Parameter", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	if email == "" {
		api.Respond(w, "Invalid Email Parameter", http.StatusBadRequest)
		return
	}

	response := mAPI.PasswordChangeResponse{
		DateCreated: time.Now(),
		Success:     VerifyResetCode(email, code),
	}

	api.RespondOK(w, response)
}

// VerifyResetCode checks if the given reset code is valid # UNFINISHED PLEASE FIX
func VerifyResetCode(Email string, Code string) bool {
	coll := db.CLIENT.Database("login-api-db").Collection("password-reset-codes")

	// Search for a matching email
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: Email}}).Decode(&result)
	if err != nil {
		return false
	}

	timeCreated := result["createdAt"].(primitive.DateTime)
	parsedTime := time.Unix(int64(timeCreated)/1000, int64(timeCreated)%1000*int64(time.Millisecond))
	// Check if the current time is not before the parsed time
	if !time.Now().Before(parsedTime) {
		return false
	}

	return true
}

// PostPasswordChange godoc
// @Summary      Allows users to change their password with a valid reset code
// @Description  Checks for a reset code match and then changes the password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        password_hash query string    true "new hashed account password"
// @Param        email         query string    true "email of the user"
// @Param        code          query string    true "password reset code"
// @Success      200  {array}   mAPI.PasswordChangeResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /password-change [post]
func PostPasswordChange(w http.ResponseWriter, r *http.Request) {
	log.Info("In password-change handler -------------------------")
	r.ParseForm()
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	password_hash := r.Form.Get("password_hash")
	email := r.Form.Get("email")
	code := r.Form.Get("code")

	if code == "" {
		api.Respond(w, "Invalid Reset Code Parameter", http.StatusBadRequest)
		return
	}

	if password_hash == "" {
		api.Respond(w, "Invalid Password Parameter", http.StatusBadRequest)
		return
	}

	regex := regexp.MustCompile(`^..*@.*.\.(com|net|org)$`)

	if email == "" || !regex.MatchString(email) {
		api.Respond(w, "Invalid Email Parameter", http.StatusBadRequest)
		return
	}

	err := ChangePassword(password_hash, email, code)
	if err != nil {
		api.Respond(w, "Error changing password", http.StatusInternalServerError)
		return
	}

	response := mAPI.PasswordChangeResponse{
		DateCreated: time.Now(),
		Success:     true,
	}

	api.RespondOK(w, response)
}

// Insert Credentials Code Here
func ChangePassword(PasswordHash string, Email string, Code string) error {
	coll := db.CLIENT.Database("login-api-db").Collection("password-reset-codes")

	// Search for a matching email
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: Email}}).Decode(&result)
	if err != nil {
		return err
	}

	timeCreated := result["createdAt"].(primitive.DateTime)
	parsedTime := time.Unix(int64(timeCreated)/1000, int64(timeCreated)%1000*int64(time.Millisecond))
	// Check if the current time is not before the parsed time
	if !time.Now().Before(parsedTime) {
		return errors.New("Code has expired")
	}

	testCode := result["code"].(string)
	if Code != testCode {
		return errors.New("Code does not match")
	}

	// Update database
	usersColl := db.CLIENT.Database("login-api-db").Collection("users")

	// Define a filter to find the document with the matching email
	filter := bson.M{"email": Email}

	// Define an update to set the password field to the new password
	update := bson.M{"$set": bson.M{"password": PasswordHash}}

	// Perform the update operation
	_, err = usersColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
