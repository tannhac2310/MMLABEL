package messagerelation

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (b *messageRelationService) FindMessageRelations(ctx context.Context, opts *FindMessageRelationsOpts, limit, offset int64) ([]*repository.MessageRelationData, *repository.CountResult, error) {
	filter := &repository.SearchMessageRelationOpts{
		IDs:    opts.IDs,
		ChatID: opts.ChatID,
		UserID: opts.UserID,
		Limit:  limit,
		Offset: offset,
		Sort:   nil,
	}
	result, err := b.messageRelationRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	total, err := b.messageRelationRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return result, total, nil
}

type FindMessageRelationsOpts struct {
	IDs    []string
	ChatID string
	UserID string
}
