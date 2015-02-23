package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"time"
	"strconv"
	"os/exec"
	"net/http"
	"path/filepath"
	"encoding/json"
	"github.com/unrolled/render"
	"github.com/zenazn/goji/web"
)

func (run *Run) save() {
	runMetadata, _ := os.Create(filepath.Join(run.path, "run.json"))
	json.NewEncoder(runMetadata).Encode(run)
	defer runMetadata.Close()
}

func nextRunID(runs map[uint]*Run) uint {
	log.Printf("calling nextRunID")
	var m uint = 0
	for i, _ := range runs {
		if i >= m {
			m = i + 1
		}
	}
	log.Printf("\treturning: %d", m)
	return m
}

func (run *Run) updateStatus(status string) {
	go func() {
		run.Status = status
		run.save()
		run.plan.run_update <- 0
	}()
	<- run.plan.run_update
}

func (run *Run) finished() {
	go func() {
		run.Duration = time.Now().Sub(run.Start)
		run.save()
		run.plan.run_update <- 0
	}()
	<- run.plan.run_update
}

/*
 * Execute the steps for this plan.
 */
func (run *Run) Execute() {

	run.updateStatus("executing")

	for index, step := range run.plan.Steps {
		log.Printf("running step %s\n", step.Name)
		path, _ := os.Getwd()
		stepPath := fmt.Sprintf("%s/plans/%s/step%d", path, run.plan.Name, index)
		logPrefix := fmt.Sprintf("%s/step%d", run.path, index)

		/*
		 * We want to capture the stdout and stderr from each step, so set that up here.
		 */
		stdout, err := os.Create(logPrefix + ".out")
		if err != nil {
			run.updateStatus("failure")
			log.Printf("cannot create stdout log for run! out of disk space or inodes?\n")
			break
		}

		stderr, err := os.Create(logPrefix + ".err")
		if err != nil {
			run.updateStatus("failure")
			log.Printf("cannot create stderr log for run! out of disk space or inodes?\n")
			break
		}

		cmd := exec.Command(stepPath)
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		/*
		 * Grab the current environment, and add an env var pointing at the trigger that
		 * caused this run. The body from an HTTP trigger is already dumped into this file.
		 * If the run was triggered by crontab-style execution, the trigger file will just
		 * contain the bare string "scheduled" (without quotes).
		 */
		env := os.Environ()
		env = append(env, fmt.Sprintf("CHANNON_TRIGGER=%s/trigger", run.path))
		cmd.Env = env

		err = cmd.Run()
		stdout.Close()
		stderr.Close()
		if err != nil {
			log.Printf("err was not nil! shit")
			log.Printf(err.Error())
			run.updateStatus("failure")
			break
		}
	}

	run.finished()

	if run.Status != "failure" {
		run.updateStatus("success")
	}

	for _, n := range run.plan.Notifications {
		n := n // the shadow knows. http://golang.org/doc/faq#closures_and_goroutines
		go n.Execute(run)
	}
}

/*
 * This handler will trigger a run from the current plan.
 */
func addRunHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		planName := c.URLParams["planName"]
		plan := pm.plans[planName]

		path, _ := os.Getwd()
		path = filepath.Join(path, "plans", planName, "runs", strconv.FormatUint(uint64(nextRunID(plan.Runs)), 10))

		// Save the trigger's body before we do anything else
		os.MkdirAll(path, 0755)
		trigger, _ := os.Create(filepath.Join(path, "trigger"))
		io.Copy(trigger, r.Body)
		defer trigger.Close()

		newRunID := nextRunID(plan.Runs)
		newRun := Run{Id: newRunID, Status: "pending", Trigger: "post", Start: time.Now(), plan: plan, path: path}
		newRun.save()

		go func() {
			plan.Runs[newRunID] = &newRun
			go newRun.Execute()
			plan.run_update <- 0
		}()
		<- plan.run_update

		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, map[string]string{"runID": fmt.Sprintf("%d", newRunID)})
	}
}

/*
 * Get the list of runs for a plan
 */
func listRunsHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		planName := c.URLParams["planName"]
		plan := pm.plans[planName]
		runList := make([]*Run, 0)
		for _, r := range plan.Runs {
			runList = append(runList, r)
		}

		go func() {
			ren := render.New(render.Options{})
			ren.JSON(w, http.StatusOK, runList)
			plan.run_update <- 0
		}()
		<- plan.run_update
	}
}

/*
 * Get the info for a specific run
 */
func getRunHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		planName := c.URLParams["planName"]
		runID, err := strconv.ParseUint(c.URLParams["runID"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		plan := pm.plans[planName]
		run := plan.Runs[uint(runID)]

		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, run)
	}
}

/*
 * Delete a specific run
 */
func deleteRunHandler(pm *PlanManager) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		planName := c.URLParams["planName"]
		runID, err := strconv.ParseUint(c.URLParams["runID"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

		plan := pm.plans[planName]
		run := plan.Runs[uint(runID)]

		go func() {
			os.RemoveAll(run.path)
			delete(plan.Runs, uint(runID))
			plan.run_update <- 0
		}()
		<- plan.run_update
		ren := render.New(render.Options{})
		ren.JSON(w, http.StatusOK, map[string]string{"deleted": fmt.Sprintf("%d", runID)})
	}
}
