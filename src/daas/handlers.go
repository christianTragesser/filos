package main

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"
)

//go:embed templates/dashboard.html.tmpl
var dashboardTemplate []byte

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	renderDashboardTemplate(w, dashboardTemplate)
}

func renderDashboardTemplate(w http.ResponseWriter, tmpl []byte) {
	t, _ := template.New("dashboard").Parse(string(tmpl))
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//go:embed templates/events.html.tmpl
var eventTemplate []byte

func handleEvents(w http.ResponseWriter, r *http.Request) {
	renderEventTemplate(w, eventTemplate)
}

func renderEventTemplate(w http.ResponseWriter, tmpl []byte) {
	keys, err := getEventKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, _ := template.New("events").Parse(string(tmpl))
	if err := t.Execute(w, keys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//go:embed templates/issue.html.tmpl
var issueTemplate []byte

func handleIssue(w http.ResponseWriter, r *http.Request) {
	issueID := r.PathValue("id")

	renderIssueTemplate(w, issueTemplate, issueID)
}

func renderIssueTemplate(w http.ResponseWriter, tmpl []byte, issueID string) {
	report, err := getIssueReport(issueID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := strings.Replace(report, "\n", "<br>", -1)

	issue := map[string]string{
		"id":     issueID,
		"report": html,
	}

	t, _ := template.New("issue").Parse(string(tmpl))
	if err := t.Execute(w, issue); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
