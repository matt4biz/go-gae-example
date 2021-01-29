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
	RouteKey    tag.Key
	Requests    = stats.Int64("sort/requests", "The number of requests", "1")
	RequestTime = stats.Int64("sort/request-time", "The duration of requests (ms)", "1")
)

func initMonitoring() error {
	var err error

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	service := os.Getenv("GAE_SERVICE")
	config := profiler.Config{ProjectID: projectID, Service: service}

	if err = profiler.Start(config); err != nil {
		log.Fatalf("profiling: %s", err.Error())
	}

	if RouteKey, err = tag.NewKey("route"); err != nil {
		return fmt.Errorf("user key: %w", err)
	}

	RequestView := view.View{
		Name:        "sort/requests",
		Measure:     Requests,
		Description: "Number of requests",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{RouteKey},
	}

	RequestTimeView := view.View{
		Name:        "sort/request-time",
		Measure:     RequestTime,
		Description: "Duration of requests (s)",
		Aggregation: view.Distribution(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16),
		TagKeys:     []tag.Key{RouteKey},
	}

	if err = view.Register(&RequestView, &RequestTimeView); err != nil {
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
