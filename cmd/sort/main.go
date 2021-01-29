package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	router.Use(routeMiddleware)

	router.HandleFunc("/insert", insertHandler).Methods(http.MethodGet)
	router.HandleFunc("/qsort", qsortHigh).Methods(http.MethodGet)
	router.HandleFunc("/qsortm", qsortMiddle).Methods(http.MethodGet)
	router.HandleFunc("/qsort3", qsortMedian).Methods(http.MethodGet)
	router.HandleFunc("/qsorti", qsortInsert).Methods(http.MethodGet)
	router.HandleFunc("/qsortf", qsortFlag).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
