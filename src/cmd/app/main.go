package main

import (
	"assignment2/internal/types"
	"assignment2/internal/utils"
	"assignment2/internal/web"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}
	utils.ResetUptime()
	s := web.NewService(path.Join("res", types.CSVFilePath), web.UseRestCountries{}, web.WithoutFirestore{})
	log.Fatal(http.ListenAndServe(":"+port, web.SetupRoutes(port, s)))
}
