package auth

import (
	"context"
	"net/http"
	"time"

	"ccu/api"
	db "ccu/db"
	mAPI "ccu/model/api"

	"regexp"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostCreateAccount godoc
// @Summary      Creates account for the user
// @Description  Checks for a unique username and then registers the account in the database
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        username      query string    true "username of the account"
// @Param        name          query string    true "name of the user"
// @Param        password_hash query string    true "hashed account password"
// @Param        email         query string    true "email of the user"
// @Success      200  {array}   mAPI.SignInResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /create-account [post]
func PostCreateAccount(w http.ResponseWriter, r *http.Request) {
	log.Info("In createaccount handler -------------------------")
	r.ParseForm()
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.Form.Get("username")
	password_hash := r.Form.Get("password_hash")
	email := r.Form.Get("email")
	name := r.Form.Get("name")

	if username == "" {
		api.Respond(w, "Invalid Username Parameter", http.StatusBadRequest)
		return
	}

	if name == "" {
		api.Respond(w, "Invalid Name Parameter", http.StatusBadRequest)
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

	response := mAPI.CreateAccountResponse{
		Id:          "CREATEACCOUNT",
		DateCreated: time.Now(),
		Success:     CredentialsExist(username, password_hash, email, name),
		Username:    username,
	}

	api.RespondOK(w, response)
}

// Insert Credentials Code Here
func CredentialsExist(Username string, PasswordHash string, Email string, Name string) bool {
	//checks for a specific username in the login Database
	coll := db.CLIENT.Database("login-api-db").Collection("users")

	//search a database for a certain document
	var result bson.M
	err := coll.FindOne(context.Background(), bson.D{{Key: "username", Value: Username}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		err = nil
		err := coll.FindOne(context.Background(), bson.D{{Key: "email", Value: Email}}).Decode(&result)
		if err == mongo.ErrNoDocuments {
			row := mAPI.CreateAccountDatabase{
				Username:     Username,
				PasswordHash: PasswordHash,
				Email:        Email,
				Name:         Name,
			}
			err = nil
			_, err := coll.InsertOne(context.Background(), row)
			if err != nil {
				log.Fatal(err)
				return false
			}
			return true
		}
	}
	if err != nil {
		panic(err)
	}
	return false
}
