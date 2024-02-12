package tasks

import "github.com/hibiken/asynq"

const (
	TypeDatabaseCreate = "database:create"
)

type Queue struct {
	Client *asynq.Client
	Server *asynq.Server
}

func NewQueue(url string, password string) *Queue {
	options := &asynq.RedisClientOpt{
		Addr:     url,
		Password: password,
		DB:       0,
	}
	client := asynq.NewClient(options)
	srv := asynq.NewServer(
		options,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
		},
	)

	return &Queue{Client: client, Server: srv}
}

func NewMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}
