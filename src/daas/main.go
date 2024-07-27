package main

import (
	"log/slog"
	"net/http"
	"os"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func main() {
	log.Info("Starting DaaS service...")

	_, envVarSet := os.LookupEnv("REDIS_HOST")
	if !envVarSet {
		log.Error("environment variable REDIS_HOST is not set - exiting")
		os.Exit(1)
	}

	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/alerts", handleWebhook)
	http.HandleFunc("/events", handleEvents)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/issue/", handleIssue)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Error(err.Error())
	}
}
