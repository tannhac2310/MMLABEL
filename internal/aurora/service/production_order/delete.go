package production_order

import (
	"context"
)

func (c *productionOrderService) Delete(ctx context.Context, id string) error {
	return c.productionOrderRepo.SoftDelete(ctx, id)
}
