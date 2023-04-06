package web

import (
	"log"
	"net/http"
)

func SetupRoutes(port string) *http.ServeMux {
	mux := http.ServeMux{}

	// TODO: create const
	mux.HandleFunc("/", DefaultHandler)
	mux.HandleFunc(RenewablesCurrentPath, EnergyCurrentHandler)
	mux.HandleFunc(RenewablesHistoryPath, EnergyHistoryHandler)
	mux.HandleFunc(NotificationsPath, NotificationHandler)
	mux.HandleFunc(StatusPath, StatusHandler)

	domainNamePort := "http://localhost:" + port

	log.Println("Started services on:")
	log.Println(domainNamePort + RenewablesCurrentPath)
	log.Println(domainNamePort + RenewablesHistoryPath)
	log.Println(domainNamePort + NotificationsPath)
	log.Println(domainNamePort + StatusPath)

	return &mux
}
