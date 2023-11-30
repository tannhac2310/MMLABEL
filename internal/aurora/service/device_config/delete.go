package device_config

import (
	"context"
)

func (c *deviceConfigService) Delete(ctx context.Context, id string) error {
	return c.deviceConfigRepo.SoftDelete(ctx, id)
}
