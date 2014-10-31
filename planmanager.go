package main

import (
	"errors"
)

type PlanManager struct {
	plans map[string]*Plan
	lock chan int
}

func NewPlanManager() *PlanManager {
	pm := PlanManager{}
	pm.plans = make(map[string]*Plan, 0)
	pm.lock = make(chan int)
	return &pm
}

func (pm *PlanManager) AddPlan(plan *Plan) error {
	if _, ok := pm.plans[plan.Name]; ok {
		return errors.New("The plan name is already taken!")
	}

	go func() {
		pm.plans[plan.Name] = plan
		pm.lock <- 0
	}()
	<- pm.lock

	return nil
}

func (pm *PlanManager) UpdatePlan(plan *Plan) error {
	go func() {
		pm.plans[plan.Name] = plan
		pm.lock <- 0
	}()
	<- pm.lock

	return nil
}

func (pm *PlanManager) GetPlans() []*Plan {
	plans := make([]*Plan, len(pm.plans))
	for _, plan := range pm.plans {
		plans = append(plans, plan)
	}
	return plans
}

func (pm *PlanManager) PlansSummarized() (psl PlanSummaryList) {
	for _, plan := range pm.plans {
		if plan.Name != "" {
			psl.Names = append(psl.Names, plan.Name)
		}
	}
	return psl
}

func (pm *PlanManager) DeletePlan(name string) {
	delete(pm.plans, name)
}
