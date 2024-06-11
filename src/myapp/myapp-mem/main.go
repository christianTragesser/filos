package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func memoryLeak() {
	var memoryHog [][]byte
	counter := 0

	log.Info("Starting memory allocation...")

	for {
		// Allocate 5MB of memory in each iteration
		chunk := make([]byte, 5*1024*1024)
		memoryHog = append(memoryHog, chunk)
		counter++

		memoryTotal := "Allocated " + strconv.Itoa(counter*5) + " MB of memory"
		log.Info(memoryTotal)

		time.Sleep(2 * time.Second)
	}

}

func handleApp(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	w.Write([]byte("ok"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func main() {
	go memoryLeak()

	http.HandleFunc("/", handleApp)
	http.HandleFunc("/health", handleHealth)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Error(err.Error())
	}
}
