package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mbugua/enterprise/pkg/logging"
)

func main() {
	_err := godotenv.Load("../../.env")
	if _err != nil {
		log.Fatalf("unable to load env vars. Err: %s", _err)
	}
	// Create the log file directory if it doesn't exist
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "../../logs/"
	}
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.info("<< welcome to middleware >>")

	// Create the main handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! This is you guardian promting you..."))
	})

	// Wrap the main handler with LoggerMiddleware
	http.Handle("/", logging.LoggerMiddleware(handler))

	log.Fatal(http.ListenAndServe(":8060", nil))
}
