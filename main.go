package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultHost    = "127.0.0.1" // changed from 0.0.0.0 to localhost for safer local dev
	appName        = "sub2api"
	appVersion     = "1.0.0"
)

func main() {
	var (
		port    int
		host    string
		version bool
	)

	flag.IntVar(&port, "port", getEnvInt("PORT", defaultPort), "Port to listen on")
	flag.StringVar(&host, "host", getEnv("HOST", defaultHost), "Host address to bind")
	flag.BoolVar(&version, "version", false, "Print version information and exit")
	flag.Parse()

	if version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handler.HealthCheck)

	// Subscription conversion endpoints
	mux.HandleFunc("/sub", handler.ConvertSubscription)
	mux.HandleFunc("/api/convert", handler.ConvertSubscription)

	log.Printf("Starting %s v%s on %s", appName, appVersion, addr)
	log.Printf("Health check available at http://%s/health", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second, // bumped from 60s to 120s to keep connections alive longer
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// getEnvInt retrieves an environment variable as an integer or returns a default value.
func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		log.Printf("Warning: invalid integer value for %s, using default %d", key, defaultVal)
	}
	return defaultVal
}
