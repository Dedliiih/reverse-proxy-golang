package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL.Path)
	
	w.Write([]byte("Hello World"))
}

func setupProxy() *httputil.ReverseProxy {
    log.Println("Setting up reverse proxy to NestJS server at http://localhost:3000")

	target, err := url.Parse("http://localhost:3000")

	if err != nil {
		log.Fatal("Error parsing target URL:", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

func main() {
    mux := http.NewServeMux()
    proxy := setupProxy()


	mux.Handle("/", proxy)
	mux.HandleFunc("/hello", handleHello)

	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}