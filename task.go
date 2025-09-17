package main

import "encoding/json"

type Task struct {
	TaskName string          `json:"taskName"`
	Payload  json.RawMessage `json:"payload"`
}

func NewTask(taskName string, payload json.RawMessage) *Task {
	return &Task{TaskName: taskName, Payload: payload}
}
