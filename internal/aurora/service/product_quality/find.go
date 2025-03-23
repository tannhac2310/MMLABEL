package product_quality

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

// FindProductQualityOpts defines the filter options for querying product quality.
type FindProductQualityOpts struct {
	IDs               []string
	ProductionOrderID string
	DefectTypes       []string
	CreatedAtFrom     time.Time
	CreatedAtTo       time.Time
	ProductSearch     string
	CustomerSearch    string
}

func (c *productQualityService) FindProductQuality(ctx context.Context, opts *FindProductQualityOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	params := &repository.SearchInspectionFormOpts{
		IDs:               opts.IDs,
		ProductionOrderID: opts.ProductionOrderID,
		DefectType:        opts.DefectTypes,
		CreatedAtFrom:     opts.CreatedAtFrom,
		CreatedAtTo:       opts.CreatedAtTo,
		ProductSearch:     opts.ProductSearch,
		CustomerSearch:    opts.CustomerSearch,
		Limit:             limit,
		Offset:            offset,
		Sort:              sort,
	}
	// Find inspection forms based on filter converted to repository options.
	forms, err := c.inspectionFormRepo.Search(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("lỗi khi tìm kiếm form: %w", err)
	}

	count, err := c.inspectionFormRepo.Count(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("lỗi khi đếm số lượng form: %w", err)
	}

	// Collect inspection form IDs.
	var formIDs []string
	for _, form := range forms {
		formIDs = append(formIDs, form.ID)
	}

	inspectionErrors, err := c.inspectionErrorRepo.Search(ctx, &repository.SearchInspectionErrorOpts{
		InspectionFormIDs: formIDs,
		Limit:             10000,
		Offset:            0,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("lỗi khi tìm kiếm lỗi: %w", err)
	}
	result := make([]*Data, 0, len(forms))
	for _, form := range forms {
		var errors []*repository.InspectionErrorData
		for _, e := range inspectionErrors {
			if e.InspectionFormID == form.ID {
				errors = append(errors, e)
			}
		}

		result = append(result, &Data{
			InspectionFormData: form,
			InspectionErrors:   errors,
		})
	}

	return result, count, nil
}
