package main

import (
	"log"
	"net/http"

	logging "github.com/mbugua/enterprise/pkg/logger"
)

func main() {

	// Create the main handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! This is you guardian promting you..."))
	})

	// Wrap the main handler with LoggerMiddleware
	http.Handle("/", logging.LoggerMiddleware(handler))

	log.Fatal(http.ListenAndServe(":8060", nil))
}
