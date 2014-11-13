package main

import (
	"net/http"
	"github.com/unrolled/render"
	"github.com/mholt/binding"
	"github.com/zenazn/goji/web"
)

func (plan *Plan) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&plan.Name: binding.Field{
			Form: "name",
			Required: true,
		},
		&plan.Steps: binding.Field{
			Form: "steps",
			Required: true,
		},
		&plan.Notifications: "notifications",
		&plan.Triggers: "triggers",
	}
}

func NewPlan() *Plan {
	plan := new(Plan)
	plan.run_update = make(chan int)
	return plan
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
		plan := NewPlan()
		errs := binding.Bind(r, plan)
		if errs.Handle(w) {
			return
		}

		err := pm.AddPlan(plan)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

func putPlanHandler(pm *PlanManager) (func (web.C, http.ResponseWriter, *http.Request)) {
	return func (c web.C, w http.ResponseWriter, r *http.Request) {
		plan := NewPlan()
		errs := binding.Bind(r, plan)
		if errs.Handle(w) {
			return
		}

		oldPlanName := c.URLParams["planName"]
		newPlanName := plan.Name

		if oldPlanName != newPlanName {
			err := pm.RenamePlan(oldPlanName, newPlanName)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		err := pm.UpdatePlan(plan)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, map[string]string{"updated": newPlanName})
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
