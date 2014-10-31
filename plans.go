package main

import (
	"fmt"
	"net/http"
	"github.com/unrolled/render"
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

func listPlansHandler(pm *PlanManager) (func(web.C, http.ResponseWriter, *http.Request)) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		psl := pm.PlansSummarized()
		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, psl)
	}
}

func addPlanHandler(pm *PlanManager) (func(web.C, http.ResponseWriter, *http.Request)) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		plan := new(Plan)
		errs := binding.Bind(r, plan)
		plan.run_update = make(chan int)
		if errs.Handle(w) {
			return
		}

		err := pm.AddPlan(plan)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Printf("trigger type is %s\n", plan.Trigger.Type)
	}
}

func deletePlanHandler(pm *PlanManager) (func (web.C, http.ResponseWriter, *http.Request)) {
	return func (c web.C, w http.ResponseWriter, r *http.Request) {
		planName := c.URLParams["planName"]
		pm.DeletePlan(planName)
		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, map[string]string{"deleted": planName})
	}
}

func getPlanHandler(pm *PlanManager) (func (web.C, http.ResponseWriter, *http.Request)) {
	return func (c web.C, w http.ResponseWriter, r *http.Request) {
		planName := c.URLParams["planName"]
		plan := pm.plans[planName]
		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, plan)
	}
}
