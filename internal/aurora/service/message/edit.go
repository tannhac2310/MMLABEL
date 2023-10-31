package message

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (m *messageService) EditMessage(ctx context.Context, opt *EditMessageOpts) error {
	var err error
	table := model.Message{}
	updater := cockroach.NewUpdater(table.TableName(), model.MessageFieldID, opt.ID)

	updater.Set(model.MessageFieldMessage, opt.Message)
	updater.Set(model.MessageFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update message failed %w", err)
	}
	return nil
}

type EditMessageOpts struct {
	ID      string
	Message string
	UserID  string
}
