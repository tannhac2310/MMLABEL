package order

import (
	"context"
)

func (b *orderService) SoftDelete(ctx context.Context, id string) error {
	return b.orderRepo.SoftDelete(ctx, id)
}
