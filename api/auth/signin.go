package auth

import (
	"net/http"
	"time"

	"ccu/api"
	mAPI "ccu/model/api"

	log "github.com/sirupsen/logrus"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostSignIn godoc
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
// @Router       /signin-auth [post]
func PostSignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.Form.Get("username")
	password_hash := r.Form.Get("password_hash")

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

// Sample user data (replace with database)
var users = []User{
	{Username: "user1", Password: "$2a$10$L3ufbDKLP0Cn8qKfsl5BcOzcsk.MDl0zDa.OYRf6.PSdbK7LiTaeO"},  // Hashed password for "password1"
	{Username: "user2", Password: "$2a$10$ewRmCuJzJGOS7TafJzVxjOzpiTn7PZi0CV8hXwq2p3.o4uHHDFwJ8u"}, // Hashed password for "password2"
}

// Insert Credentials Code Here
func CheckCredentials(Username string, PasswordHash string) bool {
	user, exists := userDatabase[username] //add data base connection
	if !exists {
		return false // User not found
	}

	return user.PasswordHash == PasswordHash //change with password hash variable in data base
}
