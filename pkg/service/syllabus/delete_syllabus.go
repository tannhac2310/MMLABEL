package syllabus

import (
	"context"
)

func (b *syllabusService) SoftDelete(ctx context.Context, id string) error {
	return b.syllabusRepo.SoftDelete(ctx, id)
}
