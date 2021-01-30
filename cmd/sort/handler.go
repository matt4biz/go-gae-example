package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

const (
	url      = "https://jsonplaceholder.typicode.com"
	cloudKey = "X-Cloud-Trace-Context"
)

type todo struct {
	UserID    int    `json:"userID"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	User      string `json:"user"`
	Trace     string `json:"trace"`
	Completed bool   `json:"completed"`
}

func getTraceID(r *http.Request) string {
	if traceKey := r.Header.Get(cloudKey); traceKey != "" {
		return strings.Split(traceKey, "/")[0]
	}

	return ""
}

func getLoop(r *http.Request) int {
	loop := 1

	if i, err := strconv.Atoi(r.FormValue("loop")); err == nil {
		loop = i - 1
	}

	return loop
}

func getDelay(r *http.Request) int {
	delay := 8

	if i, err := strconv.Atoi(r.FormValue("delay")); err == nil {
		delay = i
	}

	return delay
}

func routeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := strings.Split(r.URL.Path[1:], "?")[0]
		ctx, _ := tag.New(r.Context(), tag.Insert(RouteKey, route))
		traceID := getTraceID(r)
		logger := log.New(log.Writer(), traceID+" ", log.LstdFlags|log.Lmsgprefix)
		counter := countingWriter{w, 0}
		start := time.Now()

		stats.Record(ctx, Requests.M(1))
		next.ServeHTTP(&counter, r)

		done := time.Since(start).Milliseconds()

		stats.Record(ctx, RequestTime.M(done/1000))
		logger.Println(r.URL.String(), done, counter.size)
	})
}

type countingWriter struct {
	http.ResponseWriter
	size int
}

func (w *countingWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}
