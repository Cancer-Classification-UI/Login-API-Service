package auth
package main

import (
	"math/rand"
	"net/http"
	"time"
         "github.com/gin-gonic/gin"
         "golang.org/x/crypto/bcrypt" 
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

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Sample user data (replace with database)
var users = []User{
    {Username: "user1", Password: "$2a$10$L3ufbDKLP0Cn8qKfsl5BcOzcsk.MDl0zDa.OYRf6.PSdbK7LiTaeO"}, // Hashed password for "password1"
    {Username: "user2", Password: "$2a$10$ewRmCuJzJGOS7TafJzVxjOzpiTn7PZi0CV8hXwq2p3.o4uHHDFwJ8u"}, // Hashed password for "password2"
}

func main() {
    r := gin.Default()

    // API endpoint for login
    r.POST("/api/login", func(c *gin.Context) {
        var req User

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
            return
        }

        if !validateUser(req.Username, req.Password) {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
    })

    // API endpoint for account creation
    r.POST("/api/account", func(c *gin.Context) {
        var req User

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
            return
        }

        if userExists(req.Username) {
            c.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
            return
        }

        // Hash the password before storing it
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create account"})
            return
        }

        users = append(users, User{Username: req.Username, Password: string(hashedPassword)})
        c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
    })

    r.Run(":8080") // Replace with port
}

// Function to validate username and password
func validateUser(username, password string) bool {
    for _, user := range users {
        if user.Username == username {
            err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
            return err == nil
        }
    }
    return false
}

// Function to check if a user with the given username already exists
func userExists(username string) bool {
    for _, user := range users {
        if user.Username == username {
            return true
        }
    }
    return false
}
