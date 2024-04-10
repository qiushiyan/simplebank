package asynqamanger

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type EmailProcessor struct {
	log *zap.SugaredLogger
}

func (p *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Println("email processing by asynq ...")
	return nil
}
