package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arshiamohammdii/backjob"
)

type EmailPayload struct {
	To string
}

func main() {
	fmt.Println()
	emailPayload := EmailPayload{To: "arshiamoh@gmail.com"}
	payloadJson, err := json.Marshal(emailPayload)
	if err != nil {
		panic(err)
	}

	t1 := backjob.NewTask("email", payloadJson)
	client := backjob.NewBackJobClient(backjob.ClientOptions{
		Address:  "localhost:6379",
		Password: "",
		DB:       0})

	client.Enqueue(t1)

	server := backjob.NewServer(backjob.ClientOptions{
		Address: "localhost:6379",
		DB:      0})

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
