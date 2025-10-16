package main

import (
	"fmt"
	"log"
	"log/slog"
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

func getTargetUrl() string {
	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		log.Fatal("TARGET_URL environment variable is not set")
	}
	return targetURL
}

func main() {
	godotenv.Load()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

    mux := http.NewServeMux()

	targetURL := getTargetUrl()
    proxy, err := setupProxy(targetURL)

    if err != nil {
        log.Fatal(err)
    }

    logger.Info("reverse proxy configured", "target", targetURL)

	mux.Handle("/", proxy)
	mux.HandleFunc("/hello", handleHello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Info("Starting server", "port", port)
	log.Fatal(server.ListenAndServe())
}