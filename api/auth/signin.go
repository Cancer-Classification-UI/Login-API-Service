package auth

import (
	"math/rand"
	"net/http"
	"time"

	"ccu/api"
	mAPI "ccu/model/api"

	log "github.com/sirupsen/logrus"
)

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

// Insert Credentials Code Here
func CheckCredentials(Username string, PasswordHash string) bool {
	return rand.Intn(2) == 1
}
