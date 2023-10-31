package syllabus

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (b *syllabusService) FindSyllabusByID(ctx context.Context, id string) (*Data, error) {
	syllabuses, _, err := b.FindSyllabuses(ctx, &FindSyllabusesOpts{
		IDs: []string{id},
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "id",
	}, 1, 0)

	if err != nil {
		return nil, err
	}
	if len(syllabuses) != 1 {
		return nil, fmt.Errorf("syllabuses.Search:FindSyllabusByID not found")
	}

	return syllabuses[0], nil
}
