// Package backjob launches a worker to make the jobs running
package backjob

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Handler func(task *Task) error
type Server struct {
	rdb      *redis.Client
	handlers map[string]Handler
}

func NewServer(opts ClientOptions) *Server {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.Address,
		Password: opts.Password,
		DB:       opts.DB,
	})

	return &Server{
		rdb:      rdb,
		handlers: make(map[string]Handler),
	}
}
func (s *Server) RegisterHandler(taskName string, h Handler) {
	s.handlers[taskName] = h
}

func (s *Server) Run(ctx context.Context) error {
	//pop task from redis, this should happend in loop
	for {
		res, err := s.rdb.BRPop(ctx, 0*time.Second, DefaultNormalQueue).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return err
		}

		payload := res[1]

		var task Task
		if err := json.Unmarshal([]byte(payload), &task); err != nil {
			fmt.Println("Invalid task:", err)
			continue
		}

		if handler, ok := s.handlers[task.TaskName]; ok {
			go func() {
				if err := handler(&task); err != nil {
					fmt.Println("something wrong with the task")
				}
			}()
		} else {
			fmt.Println("No handler for:", task.TaskName)
		}

	}
}
