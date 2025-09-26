package backjob

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const DefaultNormalQueue = "queue"
const DefaultDelayedQueue = "delayed-queue"

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

func ProcessIn(d time.Duration) int64 {
	now := time.Now()
	return now.Add(d).Unix()
}

func (c *BackJobClient) Enqueue(task *Task, t ...int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	if len(t) == 0 {
		_, err := c.rdb.LPush(ctx, DefaultNormalQueue, data).Result()
		if err != nil {
			return err
		}
	} else {
		_, err := c.rdb.ZAdd(ctx, DefaultDelayedQueue, redis.Z{Member: data, Score: float64(t[0])}).Result()
		if err != nil {
			return err
		}
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
