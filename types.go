package main

import (
	"time"
)

type Run struct {
	Id uint `json:"id"`
	Status string `json:"status"`
	Trigger string `json:"trigger"`
	Start time.Time `json:"start"`
	Duration time.Duration `json:"duration"`
	plan *Plan
	path string
}

type trigger struct {
	Schedule string `json:"schedule"`
	Type string `json:"type"`
}

type Notification struct {
	Target string `json:"target"`
	Payload string `json:"payload"`
}

type step struct {
	Name string `json:"name"`
	Payload string `json:"payload"`
}

type Plan struct {
	Name string `json:"name"`
	Triggers []trigger `json:"triggers"`
	Notifications []Notification `json:"notifications"`
	Steps []step `json:"steps"`
	Runs []*Run `json:"runs"`
	run_update chan int
}

type PlanSummaryList struct {
	Names []string `json:"names"`
}
