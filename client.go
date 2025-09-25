package backjob

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const DefaultNormalQueue = "queue"

type ClientOptions struct {
	Address  string
	Password string
	DB       int
}

type BackJobClient struct {
	rdb *redis.Client
}

func NewBackJobClient(options ClientOptions) *BackJobClient {
	return &BackJobClient{rdb: redis.NewClient(&redis.Options{
		Addr:     options.Address,
		Password: options.Password,
		DB:       options.DB,
	})}
}

func (c *BackJobClient) Enqueue(task *Task, d ...time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	if len(d) == 0 {
		_, err := c.rdb.LPush(ctx, DefaultNormalQueue, data).Result()
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("this feature is not implemented yet")
	}
	return nil
}

func (c *BackJobClient) Ping(ctx context.Context) {
	pong, err := c.rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("could not connect to server: %s", err)
	}
	fmt.Printf("succesfully connected to redis server: %s", pong)
}
