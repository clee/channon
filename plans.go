package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/unrolled/render"
	"github.com/robfig/cron"
	"github.com/mholt/binding"
	"github.com/zenazn/goji/web"
)

func (plan *Plan) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&plan.Name: "name",
		&plan.Steps: "steps",
		&plan.Notification: "notify",
		&plan.Trigger: "trigger",
	}
}

func (plan *Plan) Execute() {
	for _, step := range plan.Steps {
		log.Printf("running plan %s (step %s)\n", plan.Name, step.Name)
		// don't actually run anything yet
	}
}

func listPlansHandler(resp http.ResponseWriter, req *http.Request) {
	for _, plan := range PlanManagerGetPlans() {
		fmt.Fprintf(resp, "plan: %s\n", plan.Name);
		fmt.Fprintf(resp, "number of steps: %d\n", len(plan.Steps));
		for i := 0; i < len(plan.Steps); i++ {
			fmt.Fprintf(resp, "\tstep: %s\n", plan.Steps[i].Name)
			fmt.Fprintf(resp, "\tpayload: %s\n", plan.Steps[i].Payload)
		}
	}
}

func addPlanHandler(resp http.ResponseWriter, req *http.Request) {
	plan := new(Plan)
	errs := binding.Bind(req, plan)
	if errs.Handle(resp) {
		return
	}
	err := PlanManagerAddPlan(plan)
	if err != nil {
		http.Error(resp, err.Error(), 500)
		return
	}
	fmt.Fprintf(resp, "plan: %s\n", plan.Name);
	// This is a horrible hack. Needs to be moved out into a separate cron manager.
	if plan.Trigger.Type == "scheduled" {
		fmt.Printf("adding a cron scheduled plan\n")
		c := cron.New()
		c.AddFunc(plan.Trigger.Schedule, plan.Execute)
		c.Start()
	} else {
		fmt.Printf("trigger type is %s\n", plan.Trigger.Type)
	}
	fmt.Fprintf(resp, "number of steps: %d\n", len(plan.Steps));
	for i := 0; i < len(plan.Steps); i++ {
		fmt.Fprintf(resp, "step: %s\n", plan.Steps[i].Name)
		fmt.Fprintf(resp, "payload: %s\n", plan.Steps[i].Payload)
	}
}

func deletePlanHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	planName := c.URLParams["planName"]
	PlanManagerDeletePlan(planName)
	ren := render.New(render.Options{})
	ren.JSON(w, http.StatusOK, map[string]string{"deleted": planName})
}

func getPlanHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	planName := c.URLParams["planName"]
	plan := PlanManagerGetPlan(planName)
	ren := render.New(render.Options{})
	ren.JSON(w, http.StatusOK, plan)
}
