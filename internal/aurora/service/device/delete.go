package device

import (
	"context"
)

func (c *deviceService) Delete(ctx context.Context, id string) error {
	return c.deviceRepo.SoftDelete(ctx, id)
}
