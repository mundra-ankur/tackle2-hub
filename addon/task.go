package addon

import (
	"encoding/json"
	"fmt"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/tasking"
)

//
// Task API.
type Task struct {
	// hub API client.
	client *Client
	// Addon Secret
	secret *tasking.Secret
	// Task report.
	report api.TaskReport
}

//
// Application returns the application associated with the task.
func (h *Task) Application() (r *api.Application, err error) {
	id := h.secret.Hub.Application
	if id == nil {
		err = NotFound{}
		return
	}
	r = &api.Application{}
	path := Params{api.ID: *id}.inject(api.ApplicationRoot)
	err = h.client.Get(path, r)
	return
}

//
// Data returns the addon data.
func (h *Task) Data() (d map[string]interface{}) {
	d = h.secret.Addon.(map[string]interface{})
	return
}

//
// DataWith populates the addon data object.
func (h *Task) DataWith(object interface{}) (err error) {
	b, _ := json.Marshal(h.secret.Addon)
	err = json.Unmarshal(b, object)
	return
}

//
// Started report addon started.
func (h *Task) Started() {
	h.deleteReport()
	h.report.Status = tasking.Running
	h.pushReport()
	Log.Info("Addon reported started.")
	return
}

//
// Succeeded report addon succeeded.
func (h *Task) Succeeded() {
	h.report.Status = tasking.Succeeded
	h.report.Completed = h.report.Total
	h.pushReport()
	Log.Info("Addon reported: succeeded.")
	return
}

//
// Failed report addon failed.
// The reason can be a printf style format.
func (h *Task) Failed(reason string, x ...interface{}) {
	h.report.Status = tasking.Failed
	h.report.Error = fmt.Sprintf(reason, x...)
	h.pushReport()
	Log.Info(
		"Addon reported: failed.",
		"error",
		h.report.Error)
	return
}

//
// Activity report addon activity.
// The description can be a printf style format.
func (h *Task) Activity(entry string, x ...interface{}) {
	entry = fmt.Sprintf(entry, x...)
	h.report.Activity = append(
		h.report.Activity,
		entry)
	h.pushReport()
	Log.Info(
		"Addon reported: activity.",
		"activity",
		h.report.Activity)
	return
}

//
// Total report addon total items.
func (h *Task) Total(n int) {
	h.report.Total = n
	h.pushReport()
	Log.Info(
		"Addon updated: total.",
		"total",
		h.report.Total)
	return
}

//
// Increment report addon completed (+1) items.
func (h *Task) Increment() {
	h.report.Completed++
	h.pushReport()
	Log.Info(
		"Addon updated: total.",
		"total",
		h.report.Total)
	return
}

//
// Completed report addon completed (N) items.
func (h *Task) Completed(n int) {
	h.report.Completed = n
	h.pushReport()
	Log.Info("Addon reported: completed.")
	return
}

//
// deleteReport deletes the task report.
func (h *Task) deleteReport() {
	params := Params{
		api.ID: h.secret.Hub.Task,
	}
	path := params.inject(api.TaskReportRoot)
	err := h.client.Delete(path)
	if err != nil {
		panic(err)
	}
}

//
// Bucket returns the bucket path.
func (h *Task) Bucket() (b string) {
	r := &api.Task{}
	params := Params{
		api.ID: h.secret.Hub.Task,
	}
	path := params.inject(api.TaskRoot)
	err := h.client.Get(path, r)
	if err != nil {
		panic(err)
	}
	b = r.Bucket
	return
}

//
// pushReport create/update the task report.
func (h *Task) pushReport() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()
	params := Params{
		api.ID: h.secret.Hub.Task,
	}
	path := params.inject(api.TaskReportRoot)
	if h.report.ID == 0 {
		err = h.client.Post(path, &h.report)
	} else {
		err = h.client.Put(path, &h.report)
	}
	return
}
