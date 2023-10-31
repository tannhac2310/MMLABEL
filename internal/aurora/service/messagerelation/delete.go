package messagerelation

import (
	"context"
)

func (b *messageRelationService) SoftDelete(ctx context.Context, id string) error {
	return b.messageRelationRepo.SoftDelete(ctx, id)
}
