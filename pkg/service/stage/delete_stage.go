package stage

import (
	"context"
)

func (s *stageService) SoftDelete(ctx context.Context, id string) error {
	return s.stageRepo.SoftDelete(ctx, id)
}
