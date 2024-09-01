package production_plan

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type SummaryProductionPlanOpts struct {
	StartDate time.Time
	EndDate   time.Time
}

func (c *productionPlanService) SummaryProductionPlans(ctx context.Context, opts *SummaryProductionPlanOpts) ([]*SummaryData, error) {
	filter := &repository.SummaryProductionPlanOpts{
		StartDate: opts.StartDate,
		EndDate:   opts.EndDate,
	}
	summaryDatas, err := c.productionPlanRepo.Summary(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := make([]*SummaryData, 0, len(summaryDatas))
	for _, data := range summaryDatas {
		results = append(results, &SummaryData{data})
	}

	return results, nil
}
