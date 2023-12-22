package worker

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/zhang2092/mediahls/internal/db"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskConvertHLS(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)

	config := asynq.Config{
		Concurrency: 2, // 最大并发数量
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Printf("type: %s\n", task.Type())
			log.Printf("payload: %s\n", task.Payload())
			log.Printf("process task failed\n")
		}),
	}
	server := asynq.NewServer(redisOpt, config)
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskConvertHLS, processor.ProcessTaskConvertHLS)
	return processor.server.Start(mux)
}
