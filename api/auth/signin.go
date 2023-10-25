package auth

import (
	"context"
	"net/http"
	"time"

	"ccu/api"
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

	log.Info("In signin handler -------------------------")
	response := mAPI.SignInResponse{
		Id:          "SIGNIN",
		DateCreated: time.Now(),
		Success:     CheckCredentials(username, password_hash),
		Username:    username,
	}

	api.RespondOK(w, response)
}

var client *mongo.Client

func SetClientSignIn(c *mongo.Client) {
	client = c
}

// Insert Credentials Code Here
func CheckCredentials(Username string, PasswordHash string) bool {

	//checks for a specific username in the login Database
	coll := client.Database("loginDB").Collection("user_login")

	//search a database for a certain document
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"username", Username}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Warning("No document was found with the title", Username)
		return false
	}
	if err != nil {
		panic(err)
	}

	//ping to check if mongo is successfully connected
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	return result["password"] == PasswordHash
}
