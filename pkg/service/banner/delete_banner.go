package banner

import (
	"context"
)

func (b *bannerService) SoftDelete(ctx context.Context, id string) error {
	return b.bannerRepo.SoftDelete(ctx, id)
}
