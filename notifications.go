package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func (n *Notification) Execute(run *Run) {
	/*
	 * Bail out early if this trigger doesn't need to execute for this run.
	 */
	switch n.Target {
	case "every-failure":
		if run.Status != "failure" {
			log.Printf("not running notification! %s trigger not met", n.Target)
			return
		}
	case "result-changed":
		if run.Id == 0 {
			log.Printf("not running notification! %s trigger not met", n.Target)
			return
		}
		if run.Status == run.plan.Runs[run.Id-1].Status {
			log.Printf("not running notification! %s trigger not met", n.Target)
			return
		}
	}

	path, err := os.Getwd()
	notificationPath := fmt.Sprintf("%s/plans/%s/notify-%s", path, run.plan.Name, n.Target)

	cmd := exec.Command(notificationPath)
	cmd.Stdout = nil
	cmd.Stderr = nil

	env := os.Environ()
	env = append(env, fmt.Sprintf("CHANNON_TRIGGER=%s/trigger", run.path))
	for index := range run.plan.Steps {
		stepPath := fmt.Sprintf("%s/step%d", run.path, index)
		env = append(env, fmt.Sprintf("CHANNON_STEP%d_STDOUT='%s.out'", index, stepPath))
		env = append(env, fmt.Sprintf("CHANNON_STEP%d_STDERR='%s.err'", index, stepPath))
	}

	err = cmd.Run()
	if err != nil {
		log.Printf("err was not nil! the %s notification for plan %s failed.", n.Target, run.plan.Name)
		log.Printf(err.Error())
	}
}
