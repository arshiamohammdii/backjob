// Package backjob launches a worker to make the jobs running
package backjob

import (
	"fmt"
	"sync"
)

type Worker struct {
	Client      *Client
	Concurrency int
}

type Handler func(task *Task) error

func (w *Worker) Run(handler Handler) {
	var wg sync.WaitGroup

	wg.Add(w.Concurrency)
	for i := 0; i < w.Concurrency; i++ {
		go func() {
			defer wg.Done()
			for task := range w.Client.tasks {
				if err := handler(task); err != nil {
					fmt.Printf("Task %s failed", task.TaskName)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		fmt.Println("all workers finished")
	}()
}
