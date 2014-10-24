package main

type trigger struct {
	Schedule string `json:"schedule"`
	Type string `json:"type"`
}

type notify struct {
	Payload string `json:"payload"`
}

type step struct {
	Name string `json:"name"`
	Payload string `json:"payload"`
}

type Plan struct {
	Name string `json:"name"`
	Trigger trigger `json:"trigger"`
	Notification notify `json:"notify"`
	Steps []step `json:"steps"`
}
