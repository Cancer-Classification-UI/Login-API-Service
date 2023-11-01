package auth

import (
	"context"
	"net/http"
	"time"

	"ccu/api"
	db "ccu/db"
	mAPI "ccu/model/api"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetSignIn godoc
// @Summary      Checks if credentials are correct
// @Description  Checks for a matching username and password hash in the database
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        username      query string    true "username of the account"
// @Param        password_hash query string    true "hashed account password"
// @Success      200  {array}   mAPI.SignInResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /signin [get]
func GetSignIn(w http.ResponseWriter, r *http.Request) {
	log.Info("In signin handler -------------------------")
	if r.Method != http.MethodGet {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	password_hash := r.URL.Query().Get("password_hash")

	if username == "" {
		api.Respond(w, "Invalid Username Parameter", http.StatusBadRequest)
		return
	}

	if password_hash == "" {
		api.Respond(w, "Invalid Password Parameter", http.StatusBadRequest)
		return
	}

	success, name := CheckCredentials(username, password_hash)

	if !success {
		name = ""
	}

	response := mAPI.SignInResponse{
		DateCreated: time.Now(),
		Success:     success,
		Name:        name,
	}

	api.RespondOK(w, response)
}

// Insert Credentials Code Here
func CheckCredentials(Username string, PasswordHash string) (bool, string) {
	// Checks for a specific username in the login Database
	coll := db.CLIENT.Database("login-api-db").Collection("users")

	// Search a database for a certain document
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "username", Value: Username}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Warning("No document was found with the title", Username)
		return false, ""
	} else if err != nil {
		log.Error("Error while validating password", err)
		return false, ""
	}

	// Probably should make sure we can convert name to string.
	return result["password"] == PasswordHash, result["name"].(string)
}
