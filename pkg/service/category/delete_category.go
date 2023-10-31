package category

import (
	"context"
)

func (b *categoryService) SoftDelete(ctx context.Context, id string) error {
	return b.categoryRepo.SoftDelete(ctx, id)
}
