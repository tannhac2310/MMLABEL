package comment

import (
	"context"
)

func (b *commentService) SoftDelete(ctx context.Context, id string) error {
	return b.commentRepo.SoftDelete(ctx, id)
}
