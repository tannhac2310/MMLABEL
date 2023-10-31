package attendance

import (
	"context"
)

func (b *attendanceService) SoftDelete(ctx context.Context, id string) error {
	return b.attendanceRepo.SoftDelete(ctx, id)
}
