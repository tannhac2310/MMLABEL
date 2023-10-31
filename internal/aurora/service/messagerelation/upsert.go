package messagerelation

import (
	"context"
	"fmt"
)

func (b messageRelationService) UpsertMessageRelation(ctx context.Context, opt *UpsertMessageRelationOpts) error {

	err := b.messageRelationRepo.Upsert(ctx, opt.ChatID, opt.UserID)
	if err != nil {
		return fmt.Errorf("p.messageRelationRepo.Insert: %w", err)
	}

	return nil
}

type Opts struct {
	UserID string
	ChatID string
}
