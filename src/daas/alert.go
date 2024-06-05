package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Alert struct {
	Status      string   `json:"status"`
	Application string   `json:"application"`
	Reports     []Report `json:"reports"`
	URL         string   `json:"url"`
}

type Report struct {
	Name    string `json:"name"`
	Check   string `json:"check"`
	Message string `json:"message"`
}

type Context struct {
	issueID      string
	namespace    string
	resourceType string
	name         string
	url          string
	summary      string
}

func (a *Alert) logAlert() {
	var reports string
	var reportInstances int

	appSplit := strings.Split(a.Application, ":")
	appNamespace := appSplit[0]
	appType := strings.ToLower(appSplit[1])
	appName := appSplit[2]

	for i, report := range a.Reports {
		r := fmt.Sprintf(" reason[%s]: %s (%s) -- %s -|", strconv.Itoa(i+1), report.Name, report.Check, report.Message)
		reports = reports + r
		reportInstances++
	}

	message := fmt.Sprintf("status: %s --"+" application: %s/%s/%s --"+" Coroot dashboard url: %s --"+" %s alert summary:%s", a.Status, appNamespace, appType, appName, a.URL, strconv.Itoa(reportInstances), reports)

	log.Info(message)
}

func (a *Alert) setContext() Context {
	var issues string

	urlSplit := strings.Split(a.URL, "=")
	id := urlSplit[1]

	appSplit := strings.Split(a.Application, ":")
	appNamespace := appSplit[0]
	appType := strings.ToLower(appSplit[1])
	appName := appSplit[2]

	for i, report := range a.Reports {
		issue := fmt.Sprintf("Issue %s is a %s issue based on %s because %s. ", strconv.Itoa(i+1), report.Name, strings.ToLower(report.Check), report.Message)
		issues = issues + issue
	}

	alertSummary := fmt.Sprintf("%s issue(s) exists. %s", strconv.Itoa(len(a.Reports)), issues)

	context := Context{
		issueID:      id,
		namespace:    appNamespace,
		resourceType: appType,
		name:         appName,
		url:          a.URL,
		summary:      alertSummary,
	}

	return context
}
