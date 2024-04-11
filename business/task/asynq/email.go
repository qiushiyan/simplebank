package asynqamanger

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/hibiken/asynq"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"go.uber.org/zap"
)

type EmailProcessor struct {
	log *zap.SugaredLogger
}

func (p *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	p.log.Info("email processing by asynq ...")
	payload := t.Payload()
	p.log.Info(fmt.Sprintf("email payload: %s", payload))
	return nil
}

func (m *AsynqManager) NewEmailDeliveryTask(to, subject, template string) (*asynq.Task, error) {
	payload := taskcommon.EmailDeliveryPayload{
		To:       to,
		Subject:  subject,
		Template: template,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(taskcommon.TypeEmailDelivery, b), nil
}
