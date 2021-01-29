package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if err := initMonitoring(); err != nil {
		log.Fatal(err.Error())
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
