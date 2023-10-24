package main

import (
	_ "encoding/json"
	"fmt"
	_ "io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	_ "time"

	auth "ccu/api/auth"
	test "ccu/api/test"

	log "github.com/sirupsen/logrus"

	_ "ccu/docs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/thedevsaddam/gojsonq"
)

// @title           Login API
// @version         1.0
// @description     This service is responsible for login logic. Handles login events and account creation
// @license.name    MIT License
// @license.url     https://opensource.org/license/mit/
// @BasePath  /api/v1
func main() {
	fmt.Println("Starting Login-API microservice...")
	fmt.Println("No logs will be generated here. Please see log.txt file for logging")

	CreateLog()
	SetupLog()
	SetupEndpoint()
}

// Requests
func handleRequests(r *mux.Router) {
	r.HandleFunc("/api/v1/test-no-auth", test.GetTest).Methods("GET")
	r.HandleFunc("/api/v1/signin-auth", auth.PostSignIn).Methods("POST")
	r.HandleFunc("/api/v1/create-account-auth", auth.PostCreateAccount).Methods("POST")
}

// Build log output file
func CreateLog() {
	os.Remove("log.txt")      // remove old log
	file, err := os.OpenFile( // create new log
		"log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666)

	if err != nil {
		fmt.Errorf("Cannot create a log file: ", err)
		os.Exit(1)
	}

	log.SetOutput(file)
}

// Load in .env variables and setup logging
func SetupLog() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file, program will terminate: ", err)
	}

	// Check if we should be logging methods along log messages
	methodLogging := os.Getenv("METHOD_LOGGING")
	if methodLogging == "" {
		log.Warning("METHOD_LOGGING not specified in .env, defaulting to false")
		methodLogging = "false"
	}
	log.SetReportCaller(methodLogging == "true")

	// Trace, Debug, Info, Warn, Error, Fatal, and Panic (oridnal 6 - 0)
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		log.Warning("LOG_LEVEL not specified in .env, defaulting to info")
		logLevel = "info"
	}

	// Parse string to log level and set global log level
	if parsedLevel, err := log.ParseLevel(logLevel); err != nil {
		log.Error("Invalid log level, defaulting to info: ", err)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(parsedLevel)
	}

	log.Info("STARTING LOG...")
	log.Info("LOG_LEVEL: " + logLevel)
	log.Info("METHOD_LOGGING: " + methodLogging)
}

// Setup http as a go routine
func SetupHttp(APP_PORT string, r *mux.Router, wg *sync.WaitGroup) {
	log.Info("Listening and serving on HTTP port", APP_PORT)
	log.Error(http.ListenAndServe(APP_PORT, r))
	Cleanup()
	wg.Done()
}

// Sets up swagger and serves it
func SetupSwagger(APP_PORT string, r *mux.Router, wg *sync.WaitGroup) {
	// Serve Swagger UI at the root URL
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	log.Info("Swagger is served on url: http://localhost" + APP_PORT + "/swagger/")
	wg.Done()
}

// Sets up the end points for the microservice, incl. swagger.
func SetupEndpoint() {
	// Requests
	r := mux.NewRouter()
	handleRequests(r)

	APP_PORT := os.Getenv("APP_PORT")

	// Adds waitgroup to wait for os signal or http server failure
	var wg sync.WaitGroup
	wg.Add(2) // 2 because Swagger and REST API point

	// Listen and serve
	go SetupHttp(APP_PORT, r, &wg)
	go SetupSwagger(APP_PORT, r, &wg)

	// OS signal handler
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go WaitForOSSignal(sig, &wg)

	wg.Wait() // Wait for all routines to finish to finish (Only happens if interrupted or exit or error)
	close(sig)
}

// Waits for os signal as a go routine
func WaitForOSSignal(sig chan os.Signal, wg *sync.WaitGroup) {
	conn := <-sig
	fmt.Println("Received os signal, shutting down: ", conn)
	Cleanup()
	wg.Done()
}

// Performs cleanup of service to make sure no leaks of resources
func Cleanup() {
	fmt.Println("Cleaning up!")
}
