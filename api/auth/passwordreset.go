package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"time"

	"ccu/api"
	_ "ccu/db"
	mAPI "ccu/model/api"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostForgotPassword godoc
// @Summary      Password Reset for user
// @Description  Checks for database for email and then sends a reset code to the email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        email         query string    true "email of the user"
// @Success      200  {array}   mAPI.PasswordResetResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /password-reset [post]
func PostForgotPassword(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling forgot password request")

	// Make sure only POST requests are processed
	if r.Method != http.MethodPost {
		api.Respond(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the email from the request
	r.ParseForm()
	email := r.Form.Get("email")
	emailRegex := regexp.MustCompile(`^..*@.*\.(com|net|org)$`)
	if !emailRegex.MatchString(email) {
		api.Respond(w, "Invalid Email Parameter", http.StatusBadRequest)
		return
	}

	// Generate a reset code
	code, err := GenerateResetCode() // Now it will generate a 6-digit number
	if err != nil {
		api.Respond(w, "Error generating reset code", http.StatusInternalServerError)
		return
	}

	// Store the reset code in the database
	err = StoreResetCode(email, code)
	if err != nil {
		api.Respond(w, "Error storing reset code", http.StatusInternalServerError)
		return
	}

	// Send the password reset email
	err = SendPasswordResetEmail(email, code)
	if err != nil {
		api.Respond(w, "Error sending reset email", http.StatusInternalServerError)
		return
	}

	api.RespondOK(w, fmt.Sprintf("Password reset email sent to %s", email))
	response := mAPI.PasswordResetResponse{
		DateCreated: time.Now(),
		Success:     true,
	}

	api.RespondOK(w, response)
}

// GenerateResetCode creates a 6-digit reset code
func GenerateResetCode() (string, error) {
	const min = 100000 // Minimum 6-digit number
	const max = 999999 // Maximum 6-digit number

	// Generate a random number within the range.
	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return "", err
	}

	num := n.Int64() + min

	// Return as a string.
	return fmt.Sprintf("%06d", num), nil // Ensure the string has 6 digits
}

// StoreResetCode in the database for later verification
func StoreResetCode(email string, code string) error {
	// Connect to the database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("m")) // dont think this is right
	if err != nil {
		log.Errorf("Failed to connect to the database: %v", err)
		return err
	}
	defer client.Disconnect(context.TODO())

	// Get the collection where reset codes are stored
	collection := client.Database("database-name-here").Collection("password_reset_codes")

	// Create a PasswordResetCode instance
	resetCode := bson.M{
		"email":     email,
		"code":      code,
		"createdAt": time.Now(), // Note: Added a timestamp for the creation of the reset code.
	}

	// Insert the reset code into the database
	_, err = collection.InsertOne(context.Background(), resetCode)
	if err != nil {
		log.Errorf("Failed to insert reset code into the database: %v", err)
		return err
	}

	return nil
}

// SendPasswordResetEmail calls an external API to send the reset code
func SendPasswordResetEmail(email string, code string) error {
	// Placeholder for sending email, replace with actual call to external API
	log.Infof("Sending password reset code %s to email %s", code, email)

	// Example POST request to the notification API (replace with actual code)
	// response, err := http.Post("notification-api-url", "application/json", bytes.NewBufferString(payload))
	// if err != nil {
	//     return err
	// }
	// defer response.Body.Close()

	return nil
}
