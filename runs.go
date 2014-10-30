package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"time"
	"os/exec"
	"net/http"
	"github.com/zenazn/goji/web"
)

func (run *Run) Execute(steps []step) error {
	run.Status = "executing"
	for index, step := range steps {
		log.Printf("running step %s\n", step.Name)
		stepPath := fmt.Sprintf("%s/%d", run.path, index)

		/*
		 * We want to capture the stdout and stderr from each step, so set that up here.
		 */
		stdout, err := os.Create(stepPath + ".out")
		if err != nil {
			run.Status = "failed"
			log.Printf("cannot create stdout log for run! out of disk space or inodes?\n")
			return err
		}
		stderr, err := os.Create(stepPath + ".err")
		if err != nil {
			run.Status = "failed"
			log.Printf("cannot create stderr log for run! out of disk space or inodes?\n")
			return err
		}

		/*
		 * Take the payload from this step, turn it into an executable script, and run it.
		 */
		exe, err := os.Create(stepPath)
		if err != nil {
			run.Status = "failed"
			log.Printf("could create file for payload! out of disk space or inodes?\n")
			return err
		}
		exe.WriteString(step.Payload)
		log.Printf("successfully wrote payload to disk!\n")
		exe.Chmod(0755)
		log.Printf("successfully changed payload permissions!\n")

		cmd := exec.Command(stepPath)
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		/*
		 * Grab the current environment, and add an env var pointing at the trigger that
		 * caused this run. The body from an HTTP trigger is already dumped into this file,
		 * If the run was triggered by crontab-style execution, the trigger file will just
		 * contain the bare string "scheduled" (without quotes).
		 */
		env := os.Environ()
		env = append(env, fmt.Sprintf("CHANNON_TRIGGER=%s/trigger", run.path))
		cmd.Env = env

		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

/* 
 * This handler will trigger a run from the current plan.
 */
func addRunHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		planName := c.URLParams["planName"]
		plan := pm.GetPlan(planName)

		go func() {
			path, _ := os.Getwd()
			path = fmt.Sprintf("%s/runs/%d", path, len(plan.Runs))

			// Save the trigger's body before we do anything else
			os.MkdirAll(path, 0755)
			trigger, _ := os.Create(path + "/trigger")
			io.Copy(trigger, r.Body)
			trigger.Close()

			newRun := Run{Id: uint(len(plan.Runs)), Status: "pending", Trigger: "post", Start: time.Now(), plan: plan, path: path}
			plan.addRun(newRun)
		}()
	}
}

func getRunHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
	}
}


func deleteRunHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
	}
}
