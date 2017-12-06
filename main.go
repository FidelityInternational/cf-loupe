package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/FidelityInternational/cf-loupe/cf"
)

func main() {
	cfClients, err := cf.BuildClientsFromEnvironment(os.Environ())
	if err != nil {
		log.Fatal(err)
	}

	router := BuildRouter(cfClients, time.Now)
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	webServer := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting app on port %s\n", port)
	log.Fatal(webServer.ListenAndServe())
}
