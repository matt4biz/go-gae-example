package main

import (
	"log"
	"os"

	"cloud.google.com/go/profiler"
)

func initMonitoring() error {
	var err error

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	service := os.Getenv("GAE_SERVICE")
	config := profiler.Config{ProjectID: projectID, Service: service}

	if err = profiler.Start(config); err != nil {
		log.Fatalf("profiling: %s", err.Error())
	}

	return nil
}
