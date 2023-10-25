package auth

import (
	"net/http"
	"time"

	"ccu/api"
	mAPI "ccu/model/api"

	log "github.com/sirupsen/logrus"
	"regexp"
)

// PostCreateAccount godoc
// @Summary      Creates account for the user
// @Description  Checks for a unique username and then registers the account in the database
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        username      query string    true "username of the account"
// @Param        password_hash query string    true "hashed account password"
// @Param        email         query string    true "email of the user"
// @Success      200  {array}   mAPI.SignInResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /create-account-auth [post]
func PostCreateAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.Form.Get("username")
	password_hash := r.Form.Get("password_hash")
	email := r.Form.Get("email")

	if username == "" {
		api.Respond(w, "Invalid Username Parameter", http.StatusBadRequest)
		return
	}

	if password_hash == "" {
		api.Respond(w, "Invalid Password Parameter", http.StatusBadRequest)
		return
	}

	regex := regexp.MustCompile("^..*@.*.\\.(com|net|org)$")
	
	if email == "" || !regex.MatchString(email) {
		api.Respond(w, "Invalid Email Parameter", http.StatusBadRequest)
		return
	}

	log.Info("In createaccount handler -------------------------")
	response := mAPI.CreateAccountResponse{
		Id:          "CREATEACCOUNT",
		DateCreated: time.Now(),
		Success:     CredentialsExist(username),
		Username:    username,
	}

	api.RespondOK(w, response)
}

// Insert Credentials Code Here
func CredentialsExist(Username string) bool {
	for i := 0; i < len(users); i++ {
		if users[i].Username == Username {
			return false
		}
	}
	return true
}
