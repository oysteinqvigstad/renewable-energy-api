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
	http.HandleFunc(handlers.RenewablesCurrentPath, handlers.EnergyCurrentHandler)
	http.HandleFunc(handlers.RenewablesHistoryPath, handlers.EnergyHistoryHandler)
	http.HandleFunc(handlers.NotificationsPath, handlers.NotificationHandler)
	http.HandleFunc(handlers.StatusPath, handlers.StatusHandler)

	domainNamePort := "http://localhost:" + port

	log.Println("Started services on:")
	log.Println(domainNamePort + handlers.RenewablesCurrentPath)
	log.Println(domainNamePort + handlers.RenewablesHistoryPath)
	log.Println(domainNamePort + handlers.NotificationsPath)
	log.Println(domainNamePort + handlers.StatusPath)

	app.ResetUptime()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
