// Package backjob launches a worker to make the jobs running
package backjob

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
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
	go s.RunDelayedTasks(ctx)
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

func (s *Server) RunDelayedTasks(ctx context.Context) {
	for {
		now := time.Now().Unix()
		tasks, err := s.rdb.ZRangeByScore(ctx, DefaultDelayedQueue, &redis.ZRangeBy{Min: "-inf", Max: strconv.FormatInt(now, 10)}).Result()
		if err != nil {
			panic(err)
		}

		for _, task := range tasks {

			_, err := s.rdb.ZRem(ctx, DefaultDelayedQueue, task).Result()
			if err != nil {
				fmt.Println("something wrong with removing task from delayedQueue")
				continue
			}

			_, err1 := s.rdb.LPush(ctx, DefaultNormalQueue, task).Result()
			if err1 != nil {
				fmt.Println("something wrong with pushing task to the normal Queue")
				continue
			}

		}
		//pop from DefaultDelayedQueue
		//add to the default Queue

		fmt.Printf("%v\n", tasks)
		time.Sleep(1 * time.Second)
	}
}
