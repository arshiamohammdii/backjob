package backjob

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientOptions struct {
	Address  string
	Password string
	Db       int
}

type backJobClient struct {
	rdb *redis.Client
}

func NewbackJobClient(options ClientOptions) *backJobClient {
	return &backJobClient{rdb: redis.NewClient(&redis.Options{
		Addr:     options.Address,
		Password: options.Password,
		DB:       options.Db,
	})}
}

func (c *backJobClient) Enqueue(task *Task) error {
	return nil
}

func (c *backJobClient) EnqueueEvery(d time.Duration, task *Task) error {
	return nil
}
