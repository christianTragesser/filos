package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var alert Alert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, "DaaS failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	alert.logAlert()
	context := alert.setContext()
	go runFilosInstance(context)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert received (200)"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func main() {
	log.Info("Starting DaaS service...")

	http.HandleFunc("/alerts", handleWebhook)
	http.HandleFunc("/health", handleHealth)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Error(err.Error())
	}
}
