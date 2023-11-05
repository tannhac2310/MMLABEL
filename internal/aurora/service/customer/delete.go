package customer

import (
	"context"
)

func (c *customerService) Delete(ctx context.Context, id string) error {
	return c.customerRepo.SoftDelete(ctx, id)
}
