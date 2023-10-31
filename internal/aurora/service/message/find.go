package message

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (m *messageService) FindMessages(ctx context.Context, opts *FindMessagesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchMessagesOpts{
		IDs:    opts.IDs,
		ChatID: opts.ChatID,
		UserID: opts.UserID,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	messages, err := m.messageRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := m.messageRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(messages))
	for _, message := range messages {
		if err != nil {
			return nil, nil, err
		}
		results = append(results, &Data{
			MessageData: message,
		})
	}
	return results, total, nil
}

type FindMessagesOpts struct {
	IDs    []string
	UserID string
	ChatID string
}
