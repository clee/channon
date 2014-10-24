package main

import (
	"errors"
	"sync"
)

type planmanager struct {
	once sync.Once
	plans map[string]*Plan
	lock chan int
}

var PlanManager planmanager

func PlanManagerLoad() {
	PlanManager.once.Do(func() {
		PlanManager.plans = make(map[string]*Plan, 0)
		PlanManager.lock = make(chan int)
	})
}

func PlanManagerAddPlan(plan *Plan) error {
	PlanManagerLoad()
	for i := 0; i < len(PlanManager.plans); i++ {
		if _, ok := PlanManager.plans[plan.Name]; ok {
			return errors.New("The plan name is already taken!")
		}
	}
	go func() {
		PlanManager.plans[plan.Name] = plan
		PlanManager.lock <- 0
	}()
	<- PlanManager.lock
	return nil
}

func PlanManagerGetPlans() []*Plan {
	plans := make([]*Plan, len(PlanManager.plans))
	for _, plan := range PlanManager.plans {
		plans = append(plans, plan)
	}
	return plans
}

func PlanManagerGetPlan(name string) *Plan {
	return PlanManager.plans[name]
}

func PlanManagerDeletePlan(name string) {
	delete(PlanManager.plans, name)
}
