package main

func PlansSummarized(plans []*Plan) PlanSummaryList {
	psl := PlanSummaryList{Names: []string{}}
	for _, plan := range plans {
		psl.Names = append(psl.Names, plan.Name)
	}
	return psl
}
