package course

import (
	"context"
)

func (c *courseService) SoftDelete(ctx context.Context, id string) error {
	return c.courseRepo.SoftDelete(ctx, id)
}
