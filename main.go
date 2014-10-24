package main

import (
	"github.com/zenazn/goji"
)

func main() {
	PlanManagerLoad()
	goji.Get("/plans", listPlansHandler)
	goji.Post("/plans", addPlanHandler)
	goji.Delete("/plans/:planName", deletePlanHandler)
	goji.Get("/plans/:planName", getPlanHandler)
	goji.Serve()
}
