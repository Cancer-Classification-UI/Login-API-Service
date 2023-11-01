package db

import (
	"context"
	_ "encoding/json"
	"math"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CLIENT *mongo.Client

// Initialize DB
func InitDB() {
	//get Mongo URI for connection
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI'")
	}

	log.Info("Using MongoDB URI: " + uri)

	//Setup connection to mongo, Maybe setup exponential backoff and reconnect if connection lost?
	client, err := mongo.Connect(context.Background(),
		options.Client().
			ApplyURI(uri).
			SetConnectTimeout(5*time.Second).
			SetServerSelectionTimeout(5*time.Second))

	if err != nil {
		log.Fatal(err)
	}

	// Ping database, if not successful, exponential backoff N times
	retries := 3
	for i := 1; i <= retries; i++ {
		if err := client.Ping(context.Background(), nil); err == nil {
			log.Info("Connected to MongoDB!")
			break
		} else {

			if i == retries {
				log.Fatal("Could not connect to MongoDB, terminating...")
			}

			log.Warn("Could not connect to MongoDB, retrying...")
			time.Sleep(exponentialBackoff(i)) // Exponential backoff
		}
	}

	CLIENT = client
}

func exponentialBackoff(i int) time.Duration {
	return time.Duration(math.Pow(2, float64(i))) * time.Second
}

func Cleanup() {
	if err := CLIENT.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}
