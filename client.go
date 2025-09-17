package backjob

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Client struct {
	tasks chan *Task
	cron  *cron.Cron
}

func NewClient(bufferSize int) *Client {
	return &Client{tasks: make(chan *Task, bufferSize), cron: cron.New(cron.WithSeconds())}
}

func (c *Client) Enqueue(task *Task) error {
	c.tasks <- task
	return nil

}

func (c *Client) EnqueueEvery(d time.Duration, task *Task) error {
	go func() {
		tick := time.NewTicker(d)
		defer tick.Stop()
		for range tick.C {
			_ = c.Enqueue(task)
		}
	}()
	return nil
}
