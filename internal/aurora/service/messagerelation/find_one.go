package messagerelation

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (b *messageRelationService) FindOne(ctx context.Context, opts *FindMessageRelationsOpts) (*repository.MessageRelationData, error) {
	filter := &repository.SearchMessageRelationOpts{
		IDs:    opts.IDs,
		ChatID: opts.ChatID,
		UserID: opts.UserID,
		Limit:  1,
		Offset: 0,
	}
	return b.messageRelationRepo.SearchOne(ctx, filter)
}
