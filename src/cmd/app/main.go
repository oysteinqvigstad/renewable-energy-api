package main

import (
	"assignment2/internal/app"
	"assignment2/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080" // TODO: use const
	}

	// TODO: create const
	http.HandleFunc("/", handlers.DefaultHandler)

	log.Println("Starting server on port " + port + " ...")
	log.Println("http://localhost:" + port)

	app.ResetUptime()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
