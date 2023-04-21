package web

import (
	"log"
	"net/http"
)

// SetupRoutes configures the routes for the application and returns a ServeMux with the registered routes.
func SetupRoutes(port string, s *State) *http.ServeMux {
	mux := http.ServeMux{}

	// Registering route handlers
	mux.HandleFunc("/", DefaultHandler)
	mux.HandleFunc(RenewablesCurrentPath, s.EnergyCurrentHandler)
	mux.HandleFunc(RenewablesHistoryPath, s.EnergyHistoryHandler)
	mux.HandleFunc(NotificationsPath, s.NotificationHandler)
	mux.HandleFunc(StatusPath, s.StatusHandler)

	// Constructing the base domain name with the provided port
	domainNamePort := "http://localhost:" + port

	// Logging the available services and their endpoints
	log.Println("Started services on:")
	log.Println(domainNamePort + RenewablesCurrentPath)
	log.Println(domainNamePort + RenewablesHistoryPath)
	log.Println(domainNamePort + NotificationsPath)
	log.Println(domainNamePort + StatusPath)

	return &mux
}
