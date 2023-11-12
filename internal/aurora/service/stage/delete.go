package stage

import (
	"context"
)

func (c *stageService) Delete(ctx context.Context, id string) error {
	return c.stageRepo.SoftDelete(ctx, id)
}
