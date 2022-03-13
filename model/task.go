package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type TaskReport struct {
	Model
	Status    string
	Error     string
	Total     int
	Completed int
	Activity  JSON
	TaskID    uint `gorm:"<-:create;uniqueIndex"`
	Task      *Task
}

type Task struct {
	Model
	Bucket
	Name        string `gorm:"index"`
	Addon       string `gorm:"index"`
	Locator     string `gorm:"index"`
	Image       string
	Isolated    bool
	Data        JSON
	Started     *time.Time
	Terminated  *time.Time
	Status      string
	Error       string
	Job         string
	Report      *TaskReport `gorm:"constraint:OnDelete:CASCADE"`
	TaskGroupID *uint
	TaskGroup   *TaskGroup
}

func (m *Task) Reset() {
	m.Started = nil
	m.Terminated = nil
	m.Report = nil
}

func (m *Task) BeforeCreate(db *gorm.DB) (err error) {
	err = m.Bucket.BeforeCreate(db)
	m.Reset()
	return
}

type TaskGroup struct {
	Model
	Bucket
	Name  string
	Addon string
	Data  JSON
	Tasks []Task `gorm:"constraint:OnDelete:CASCADE"`
}

func (m *TaskGroup) BeforeCreate(db *gorm.DB) (err error) {
	err = m.Bucket.BeforeCreate(db)
	if err != nil {
		return
	}
	err = m.Propagate()
	return
}

func (m *TaskGroup) BeforeUpdate(*gorm.DB) (err error) {
	err = m.Propagate()
	return
}

func (m *TaskGroup) BeforeDelete(db *gorm.DB) (err error) {
	err = m.Bucket.BeforeDelete(db)
	if err != nil {
		return
	}
	return
}

//
// Propagate group data into the task.
func (m *TaskGroup) Propagate() (err error) {
	for i := range m.Tasks {
		task := &m.Tasks[i]
		if task.Addon == "" {
			task.Addon = m.Addon
		}
		if m.Data == nil {
			continue
		}
		a := Map{}
		err = json.Unmarshal(m.Data, &a)
		if err != nil {
			return
		}
		b := Map{}
		err = json.Unmarshal(task.Data, &b)
		if err != nil {
			return
		}
		task.Data, _ = json.Marshal(m.merge(a, b))
	}

	return
}

//
// merge maps B into A.
// The B map is the authority.
func (m *TaskGroup) merge(a, b Map) (out Map) {
	if a == nil || b == nil {
		return
	}
	out = Map{}
	//
	// Merge-in elements found in B and in A.
	for k, v := range a {
		out[k] = v
		if bv, found := b[k]; found {
			out[k] = bv
			if av, cast := v.(Map); cast {
				if bv, cast := bv.(Map); cast {
					out[k] = m.merge(av, bv)
				} else {
					out[k] = bv
				}
			}
		}
	}
	//
	// Add elements found only in B.
	for k, v := range b {
		if _, found := a[k]; !found {
			out[k] = v
		}
	}

	return
}

//
// Map alias.
type Map = map[string]interface{}
