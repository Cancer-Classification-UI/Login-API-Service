package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"ccu/api"
	db "ccu/db"
	mAPI "ccu/model/api"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

// PostPasswordChangeEmail godoc
// @Summary      Password change email functionality
// @Description  Checks for database for email and then sends a reset code to the email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        email         query string    true "email of the user"
// @Success      200  {array}   mAPI.PasswordChangeResponse
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /password-change-email [post]
func PostPasswordChangeEmail(w http.ResponseWriter, r *http.Request) {
	log.Info("In password-change email handler -------------------------")

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
		log.Error("Error generating reset code", err)
		api.Respond(w, "Error generating reset code", http.StatusInternalServerError)
		return
	}

	// Store the reset code in the database
	err = StoreResetCode(email, code)
	if err != nil {
		log.Error("Error storing reset code", err)
		api.Respond(w, "Error storing reset code", http.StatusInternalServerError)
		return
	}

	// Send the password reset email
	err = SendPasswordResetEmail(email, code)
	if err != nil {
		log.Error("Error sending reset email", err)
		api.Respond(w, "Error sending reset email", http.StatusInternalServerError)
		return
	}

	response := mAPI.PasswordChangeResponse{
		DateCreated: time.Now(),
		Success:     true,
	}

	api.RespondOK(w, response)
}

// GenerateResetCode creates a 6-digit reset code
func GenerateResetCode() (string, error) {
	const min = 0      // Minimum 6-digit number
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
	// Check for valid email
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result struct{}
	emailCollection := db.CLIENT.Database("login-api-db").Collection("users")

	err := emailCollection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&result)
	if err != nil {
		return err
	}

	// Get the collection where reset codes are stored
	collection := db.CLIENT.Database("login-api-db").Collection("password-reset-codes")

	// Define the filter to match documents with a specific email
	filter := bson.M{"email": email}

	// Delete any previous email codes
	_, err = collection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Create a PasswordResetCode instance
	resetCode := bson.M{
		"email":     email,
		"code":      code,
		"createdAt": time.Now().Add(5 * time.Minute), // Note: Added a timestamp for the expiration of the reset code.
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

	// Get the base URL and endpoint from the environment variables
	baseURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if baseURL == "" {
		log.Warning("NOTIFICATION_SERVICE_URL not specified in .env, defaulting to http://localhost:8087")
		baseURL = "http://localhost:8087"
	}

	endpoint := os.Getenv("NOTIFICATION_SERVICE_API_ENDPOINT")
	if endpoint == "" {
		log.Warning("NOTIFICATION_SERVICE_API_ENDPOINT not specified in .env, defaulting to /api/v1")
		endpoint = "/api/v1"
	}

	sendUrl := baseURL + endpoint + "/send-code"

	// Email address and code to be sent in the request body
	form := url.Values{}
	form.Add("email", email)
	form.Add("code", code)

	// Create a request with POST method, specifying the URL and request body
	_, err := http.PostForm(sendUrl, form)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	return nil
}
