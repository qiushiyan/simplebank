package asynqamanger

import (
	"context"
	"fmt"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/hibiken/asynq"
	"github.com/qiushiyan/simplebank/business/email"
	taskcommon "github.com/qiushiyan/simplebank/business/task/common"
	"go.uber.org/zap"
)

// Helper for creating a new email delivery task.
func (m *AsynqManager) NewEmailDeliveryTask(payload email.SenderPayload) (*asynq.Task, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		taskcommon.TypeEmailDelivery,
		b,
		asynq.ProcessIn(60*time.Second),
	), nil
}

type EmailProcessor struct {
	log    *zap.SugaredLogger
	sender *email.GmailSender
}

func NewEmailProcessor(
	log *zap.SugaredLogger,
	senderAddress, senderPassword string,
) *EmailProcessor {
	return &EmailProcessor{
		log:    log,
		sender: email.NewGmailSender(senderAddress, senderPassword),
	}
}

func (p *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload email.SenderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", err)
	}

	err := p.sender.Send(&payload)

	if err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	p.log.Infow("completed email task", "payload", payload)
	return nil
}
