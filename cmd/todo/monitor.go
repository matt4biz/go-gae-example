package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	UserKey  tag.Key
	Requests = stats.Int64("todo/requests", "The number of requests", "1")
)

func initMonitoring() error {
	var err error

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	service := os.Getenv("GAE_SERVICE")
	config := profiler.Config{ProjectID: projectID, Service: service}

	if err = profiler.Start(config); err != nil {
		log.Fatalf("profiling: %s", err.Error())
	}

	if UserKey, err = tag.NewKey("user"); err != nil {
		return fmt.Errorf("user key: %w", err)
	}

	RequestView := view.View{
		Name:        "todo/requests",
		Measure:     Requests,
		Description: "Number of requests",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{UserKey},
	}

	if err = view.Register(&RequestView); err != nil {
		return fmt.Errorf("register view: %w", err)
	}

	opts := stackdriver.Options{ProjectID: projectID}
	exp, err := stackdriver.NewExporter(opts)

	if err != nil {
		return fmt.Errorf("exporter: %w", err)
	}

	view.RegisterExporter(exp)
	view.SetReportingPeriod(60 * time.Second)
	return nil
}
