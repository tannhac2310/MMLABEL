package message

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (m *messageService) FindOne(ctx context.Context, opts *FindMessagesOpts) (*repository.MessageData, error) {
	filter := &repository.SearchMessagesOpts{
		IDs:    opts.IDs,
		ChatID: opts.ChatID,
		UserID: opts.UserID,
		Limit:  1,
		Offset: 0,
	}
	return m.messageRepo.SearchOne(ctx, filter)
}
