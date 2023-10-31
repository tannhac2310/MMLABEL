package stagestudent

import (
	"context"
)

func (b *stageStudentService) SoftDelete(ctx context.Context, id string) error {
	return b.stageStudentRepo.SoftDelete(ctx, id)
}
