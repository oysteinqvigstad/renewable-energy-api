package web

import (
	"assignment2/internal/datastore"
	"log"
	"net/http"
)

func SetupRoutes(port string, energyData datastore.RenewableDB) *http.ServeMux {
	mux := http.ServeMux{}

	// TODO: create const
	mux.HandleFunc("/", DefaultHandler)
	mux.HandleFunc(RenewablesCurrentPath, EnergyCurrentHandler(energyData))
	mux.HandleFunc(RenewablesHistoryPath, EnergyHistoryHandler(energyData))
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
