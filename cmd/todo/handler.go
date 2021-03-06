package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const (
	url      = "https://jsonplaceholder.typicode.com"
	cloudKey = "X-Cloud-Trace-Context"
)

type todo struct {
	UserID    int    `json:"userID"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Trace     string `json:"trace"`
	Completed bool   `json:"completed"`
}

func getTraceID(r *http.Request) string {
	if traceKey := r.Header.Get(cloudKey); traceKey != "" {
		return strings.Split(traceKey, "/")[0]
	}

	return ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	traceID := getTraceID(r)
	logger := log.New(log.Writer(), traceID+" ", log.LstdFlags|log.Lmsgprefix)

	if key == "" {
		logger.Println("empty request")
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	// we do it this way so we can ensure we've created the leak
	// because we're not using the default client with pooling

	req, err := http.NewRequest("GET", url+"/todos/"+key, nil)

	if err != nil {
		logger.Println("request:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cli := &http.Client{Transport: &http.Transport{}}
	resp, err := cli.Do(req)

	if err != nil {
		logger.Println("request:", err.Error())
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// right here we need defer resp.Body.Close()
	// without which we will leak goroutines & fds

	if resp.StatusCode != http.StatusOK {
		logger.Println("not found:", key)
		http.NotFound(w, r)
		return
	}

	var item todo

	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		logger.Println("decode:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.Trace = traceID

	logger.Printf("found %s: %s", key, item.Title)

	tmpl := template.New("todo")

	tmpl.Parse(form)
	tmpl.Execute(w, item)
}

var form = `
<h1>Todo #{{.ID}}</h1>
<div>{{printf "User %d" .UserID}}</div>
<div>{{printf "%s (completed: %t)" .Title .Completed}}</div>
<div>{{printf "[Trace %s]" .Trace}}</div>`
