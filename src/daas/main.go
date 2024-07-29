package main

import (
	"log/slog"
	"net/http"
	"os"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func main() {
	mux := &http.ServeMux{}

	_, envVarSet := os.LookupEnv("REDIS_HOST")
	if !envVarSet {
		log.Error("environment variable REDIS_HOST is not set - exiting")
		os.Exit(1)
	}

	log.Info("Starting DaaS service...")

	mux.HandleFunc("/", handleDashboard)
	mux.HandleFunc("/alerts", handleWebhook)
	mux.HandleFunc("/events", handleEvents)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/issue/{id}", handleIssue)
	if err := http.ListenAndServe(":5000", mux); err != nil {
		log.Error(err.Error())
	}
}
