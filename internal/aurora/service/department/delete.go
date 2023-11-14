package department

import (
	"context"
)

func (c *departmentService) Delete(ctx context.Context, id string) error {
	return c.departmentRepo.SoftDelete(ctx, id)
}
