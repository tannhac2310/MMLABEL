package office

import (
	"context"
)

func (b *officeService) SoftDelete(ctx context.Context, id string) error {
	return b.officeRepo.SoftDelete(ctx, id)
}
