package stage

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (s *stageService) FindStageByID(ctx context.Context, id string) (*Stage, error) {
	stages, total, err := s.FindStages(ctx, &FindStagesOpts{
		IDs: []string{id},
	},
		&repository.Sort{
			Order: repository.SortOrderDESC,
			By:    "id",
		}, 1, 0)

	if err != nil {
		return nil, err
	}
	if total.Count != 1 {
		return nil, fmt.Errorf("stage.Search:FindStageByID not found")
	}

	return stages[0], nil
}
