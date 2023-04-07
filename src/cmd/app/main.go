package main

import (
	"assignment2/internal/db"
	"assignment2/internal/utils"
	"assignment2/internal/web"
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
	db.GlobalRenewableDB.ParseCSV("res/renewable-share-energy.csv")
	utils.ResetUptime()
	log.Fatal(http.ListenAndServe(":"+port, web.SetupRoutes(port)))
}
