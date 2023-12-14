package product_quality

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *productQualityService) FindProductQuality(ctx context.Context, opts *FindProductQualityOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, []*ProductQualityAnalysis, error) {
	filter := &repository.SearchProductQualitysOpts{
		ProductionOrderID: opts.ProductionOrderID,
		DefectType:        opts.DefectType,
		DeviceID:          opts.DeviceID,
		DefectCode:        opts.DefectCode,
		UserID:            opts.UserID,
		CreatedAtFrom:     opts.CreatedAtFrom,
		CreatedAtTo:       opts.CreatedAtTo,
		Limit:             limit,
		Offset:            offset,
		Sort:              sort,
	}
	productQuality, err := c.productQualityRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}

	total, err := c.productQualityRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}

	results := make([]*Data, 0, len(productQuality))
	for _, productQuality := range productQuality {
		if err != nil {
			return nil, nil, nil, err
		}
		results = append(results, &Data{
			ProductQualityData: productQuality,
		})
	}

	analysis, err := c.productQualityRepo.Analysis(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}
	analysisData := make([]*ProductQualityAnalysis, 0, len(analysis))
	for _, a := range analysis {
		analysisData = append(analysisData, &ProductQualityAnalysis{
			DefectType: a.DefectType,
			Count:      a.Count,
		})
	}
	return results, total, analysisData, nil
}

type FindProductQualityOpts struct {
	ProductionOrderID string
	DeviceID          string
	DefectType        string
	DefectCode        string
	CreatedAtFrom     time.Time
	CreatedAtTo       time.Time
	UserID            string
}
