package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arshiamohammdii/backjob"
)

type EmailPayload struct {
	To string
}

func main() {
	fmt.Println()
	emailPayload := EmailPayload{To: "arshiamoh@gmail.com"}
	payloadJSON, err := json.Marshal(emailPayload)
	if err != nil {
		panic(err)
	}

	t1 := backjob.NewTask("email", payloadJSON)
	client := backjob.NewBackJobClient(backjob.ClientOptions{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
	})

	client.Enqueue(t1, backjob.ProcessIn(10*time.Second))
	// client.Enqueue(t1, Every(1 minute)) -> cron job
	// client.Enqueue(t1, In(5 minutes))

	// priority queues
	// divide tasks based on their priority

	server := backjob.NewServer(backjob.ClientOptions{
		Address: "localhost:6379",
		DB:      0,
	})

	server.RegisterHandler("email", func(t *backjob.Task) error {
		fmt.Println("sending email")
		return nil
	})

	ctx := context.Background()
	if err := server.Run(ctx); err != nil {
		panic(err)
	}

	client.Ping(context.Background())
}
