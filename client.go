package main

type Client struct {
	tasks chan *Task
}

func NewClient(bufferSize int) *Client {
	return &Client{tasks: make(chan *Task, bufferSize)}
}

func (c *Client) Enqueue(task *Task) error {
	c.tasks <- task
	return nil
}
