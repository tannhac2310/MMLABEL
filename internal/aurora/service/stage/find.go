package stage

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *stageService) FindStages(ctx context.Context, opts *FindStagesOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {

	filter := &repository.SearchStagesOpts{
		IDs:    opts.IDs,
		Name:   opts.Name,
		Code:   opts.Code,
		UserID: opts.UserID,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	stages, err := c.stageRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.stageRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(stages))
	for _, stage := range stages {
		if err != nil {
			return nil, nil, err
		}
		results = append(results, &Data{
			StageData: stage,
		})
	}
	return results, total, nil
}

type FindStagesOpts struct {
	IDs    []string
	Name   string
	Code   string
	UserID string
}
