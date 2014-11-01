package main

import (
	"os"
	"fmt"
	"log"
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
		if run.Status == run.plan.Runs[run.Id - 1].Status {
			log.Printf("not running notification! %s trigger not met", n.Target)
			return
		}
	}

	notificationPath := fmt.Sprintf("%s/notification-%s", run.path, n.Target)
	exe, err := os.Create(notificationPath)
	if err != nil {
		log.Printf("cannot create notification script! out of disk space or inodes?\n")
		log.Printf(err.Error())
		return
	}

	exe.WriteString(n.Payload)
	exe.Chmod(0755)
	exe.Close()

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
