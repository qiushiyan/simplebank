package asynqamanger

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"go.uber.org/zap"
)

type AsynqManager struct {
	log            *zap.SugaredLogger
	client         *asynq.Client
	server         *asynq.Server
	inspector      *asynq.Inspector
	senderAddr     string
	senderPassword string
}

func New(
	log *zap.SugaredLogger,
	redisAddr string,
	senderAddr string,
	senderPassword string,
) *AsynqManager {
	redisOpt := asynq.RedisClientOpt{Addr: redisAddr}

	client := asynq.NewClient(redisOpt)
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 5,
			Logger:      &Logger{log: log},
			ErrorHandler: asynq.ErrorHandlerFunc(
				func(ctx context.Context, task *asynq.Task, err error) {
					log.Errorw(
						"task processing error",
						"type",
						task.Type,
						"payload",
						string(task.Payload()),
						"error",
						err,
					)
				},
			),
		},
	)
	inspector := asynq.NewInspector(redisOpt)

	return &AsynqManager{
		log:            log,
		server:         server,
		client:         client,
		inspector:      inspector,
		senderAddr:     senderAddr,
		senderPassword: senderPassword,
	}
}

func (m *AsynqManager) Start() error {
	mux := asynq.NewServeMux()
	mux.Handle(
		taskcommon.TypeEmailDelivery,
		NewEmailProcessor(m.log, m.senderAddr, m.senderPassword),
	)

	return m.server.Run(mux)
}

func (m *AsynqManager) Close() error {
	m.inspector.Close()
	return m.client.Close()
}

func (m *AsynqManager) CreateTask(
	ctx context.Context,
	taskType string,
	payload any,
) (string, error) {
	// task := asynq.NewTask(taskType, payload)
	// info, err := m.client.Enqueue(task)

	var err error
	var info *asynq.TaskInfo
	var task *asynq.Task

	switch taskType {
	case taskcommon.TypeEmailDelivery:
		payload, ok := payload.(taskcommon.EmailDeliveryPayload)
		if !ok {
			return "", fmt.Errorf("invalid payload type for email delivery task: %T", payload)
		}
		task, err = m.NewEmailDeliveryTask(payload.To, payload.Username, payload.Subject)
		if err != nil {
			return "", err
		}
	}
	info, err = m.client.EnqueueContext(ctx, task)

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

func (m *AsynqManager) GetTaskState(id string) (*taskcommon.State, error) {
	info, err := m.inspector.GetTaskInfo("default", id)
	if err != nil {
		return &taskcommon.State{}, err
	}

	if info != nil {
		return &taskcommon.State{}, fmt.Errorf("task %s not found", id)
	}

	state := adaptState(info)

	return state, nil
}

func (m *AsynqManager) CancelTask(id string) error {
	return m.inspector.CancelProcessing(id)
}

func adaptState(info *asynq.TaskInfo) *taskcommon.State {
	var state taskcommon.State

	state.Id = info.ID
	state.Type = info.Type

	if info.LastErr != "" {
		state.Error = info.LastErr
		state.Status = taskcommon.StatusFailed
	} else {
		switch info.State {
		case asynq.TaskStateActive, asynq.TaskStatePending:
			state.Status = taskcommon.StatusInProgress
		case asynq.TaskStateCompleted:
			state.Status = taskcommon.StatusCompleted
		default:
			state.Status = taskcommon.StatusOther
		}
	}

	return &state

}
