package jaeger

import (
	"fmt"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	serverPort = "8000"
	server     = "httpServer"
	client     = "httpClient"
)

func getTime(w http.ResponseWriter, r *http.Request) {
	log.Print("Received getTime request")
	t := time.Now()
	ts := t.Format("Mon Jan _2 15:04:05 2006")
	io.WriteString(w, fmt.Sprintf("The time is %s", ts))
}

func runHTTPServer() {
	http.HandleFunc("/gettime", getTime)
	http.HandleFunc("/", getTime)
	log.Printf("Starting server on port %s", serverPort)
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", serverPort),
		// use nethttp.Middleware to enable OpenTracing for server
		nethttp.Middleware(opentracing.GlobalTracer(), http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}
