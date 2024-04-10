package asynqamanger

import (
	"github.com/hibiken/asynq"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"go.uber.org/zap"
)

type AsynqManager struct {
	log       *zap.SugaredLogger
	client    *asynq.Client
	server    *asynq.Server
	inspector *asynq.Inspector
}

func New(log *zap.SugaredLogger, redisAddr string) *AsynqManager {
	redisOpt := asynq.RedisClientOpt{Addr: redisAddr}

	client := asynq.NewClient(redisOpt)
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 5,
			Logger:      &Logger{log: log},
		},
	)
	inspector := asynq.NewInspector(redisOpt)

	return &AsynqManager{
		log:       log,
		server:    server,
		client:    client,
		inspector: inspector,
	}
}

func (m *AsynqManager) Start() error {
	mux := asynq.NewServeMux()
	mux.Handle(taskcommon.TypeEmailDelivery, &EmailProcessor{log: m.log})

	return m.server.Run(mux)
}

func (m *AsynqManager) Close() error {
	return m.client.Close()
}

func (m *AsynqManager) CreateTask(taskType string, payload []byte) (string, error) {
	task := asynq.NewTask(taskType, payload)
	info, err := m.client.Enqueue(task)

	if err != nil {
		m.log.Warnw("task created error", "error", err)
		return "", err
	}

	m.log.Infow(
		"task created",
		"task_id",
		info.ID,
		"type",
		task.Type,
		"queue",
		info.Queue,
		"payload",
		task.Payload,
	)

	return info.ID, nil
}

func (m *AsynqManager) GetTaskStatus(id string) (string, error) {
	task, err := m.inspector.GetTaskInfo("default", id)
	if err != nil {
		return "", err
	}

	return task.State.String(), nil
}

func (m *AsynqManager) CancelTask(id string) error {
	return m.inspector.CancelProcessing(id)
}
