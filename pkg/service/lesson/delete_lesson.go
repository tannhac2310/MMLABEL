package lesson

import (
	"context"
)

func (b *lessonService) SoftDelete(ctx context.Context, id string) error {
	return b.lessonRepo.SoftDelete(ctx, id)
}
