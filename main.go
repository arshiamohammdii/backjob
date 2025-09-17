package main

import (
	"encoding/json"
	"fmt"
)

// task enqueues right away
// task gets done (seconds/minutes/hours) later
// task gets done at a particualar time
// task gets done a a particualr time every /day/month/year

type EmailTaskPayload struct {
	Address string
}

func handler(task *Task) error {
	switch task.TaskName {
	case "email":
		var payload struct{ Address string }
		json.Unmarshal(task.Payload, &payload)
		fmt.Printf("sending email to: %s\n", payload.Address)
		return nil
	default:
		return fmt.Errorf("something went wrong executing the tasks")
	}
}

func main() {
	fmt.Println()
	payload, err := json.Marshal(EmailTaskPayload{Address: "something@gmail.com"})
	if err != nil {
		panic(err)
	}
	task := NewTask("email", payload)
	task2 := NewTask("email", payload)
	client := NewClient(20)
	client.Enqueue(task)
	client.Enqueue(task2)
	// close(client.tasks)

	worker := Worker{client: client, Concurrency: 2}
	worker.Run(handler)
}
