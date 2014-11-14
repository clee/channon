package main

import (
	"github.com/zenazn/goji"
)

func main() {
	planManager := NewPlanManager()
	goji.Get("/plans", listPlansHandler(planManager))
	goji.Post("/plans", addPlanHandler(planManager))

	goji.Get("/plans/:planName", getPlanHandler(planManager))
	goji.Delete("/plans/:planName", deletePlanHandler(planManager))
	goji.Put("/plans/:planName", putPlanHandler(planManager))

	goji.Post("/plans/:planName/runs", addRunHandler(planManager))
	goji.Get("/plans/:planName/runs", listRunsHandler(planManager))

	goji.Get("/plans/:planName/runs/:runID", getRunHandler(planManager))
	goji.Delete("/plans/:planName/runs/:runID", deleteRunHandler(planManager))

	goji.Get("/tags", listTagsHandler(planManager))
	goji.Post("/tags", addTagHandler(planManager))
	goji.Delete("/tags/:tagName", deleteTagHandler(planManager))

	goji.Serve()
}
