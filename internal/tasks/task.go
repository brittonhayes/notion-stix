package tasks

import "github.com/hibiken/asynq"

const (
	TypeDatabaseCreate = "database:create"
)

type Queue struct {
	*asynq.Client
}

func NewQueue(url string, password string) *Queue {
	client := asynq.NewClient(&asynq.RedisClientOpt{
		Addr:     url,
		Password: password,
	})
	return &Queue{client}
}
