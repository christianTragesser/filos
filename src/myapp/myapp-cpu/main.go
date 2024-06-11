package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func cpuSpike() {
	cmd := exec.Command(
		"stress-ng", "--cpu", "1",
		"--timeout", "120",
	)

	for {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Error("Failed to display cmd stdout")
		}

		if err := cmd.Start(); err != nil {
			log.Error("Failed to start command")
		}

		data, err := io.ReadAll(stdout)
		if err != nil {
			log.Error("Failed to read command stdout")
		}

		if err := cmd.Wait(); err != nil {
			log.Error("Failed command")
		}

		fmt.Println(string(data))

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
	go cpuSpike()

	http.HandleFunc("/", handleApp)
	http.HandleFunc("/health", handleHealth)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Error(err.Error())
	}
}
