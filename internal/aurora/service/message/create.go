package message

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (m *messageService) CreateMessage(ctx context.Context, opt *CreateMessageOpts) (string, error) {
	now := time.Now()

	return now.String(), nil
}

type CreateMessageOpts struct {
	ChatID            string
	ChatTitle         string
	EntityType        enum.ChatEntityType
	EntityID          string
	Extranet          bool
	Message           string
	Attachments       []*model.Attachment
	AuthorID          string
	Responsibility    []*model.ZaloResponsibility
	ExternalMessageID string
}
