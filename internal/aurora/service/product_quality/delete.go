package product_quality

import (
	"context"
)

func (c *productQualityService) Delete(ctx context.Context, id string) error {
	return c.productQualityRepo.SoftDelete(ctx, id)
}
