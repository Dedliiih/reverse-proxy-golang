package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received:", r.Method, r.URL.Path)
	
	w.Write([]byte("Hello World"))
}

func setupProxy(targetURL string) (*httputil.ReverseProxy, error) {
    log.Println("Setting up reverse proxy to NestJS server at http://localhost:3000")

	target, err := url.Parse(targetURL)

	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy, nil
}

func main() {
	godotenv.Load()
    mux := http.NewServeMux()
    proxy, err := setupProxy(os.Getenv("TARGET_URL"))

    if err != nil {
        log.Fatal(err)
    }

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