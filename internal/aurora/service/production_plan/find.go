package production_plan

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type Analysis struct {
	Status enum.ProductionOrderStatus `json:"status"`
	Count  int64                      `json:"count"`
}

type FindProductionPlansOpts struct {
	IDs        []string
	CustomerID string
	Name       string
	Statuses   []enum.ProductionPlanStatus
	UserID     string
	Stage      int
}

func (c *productionPlanService) FindProductionPlans(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchProductionPlanOpts{
		IDs:        opts.IDs,
		CustomerID: opts.CustomerID,
		Name:       opts.Name,
		Statuses:   opts.Statuses,
		UserID:     opts.UserID,
		Limit:      limit,
		Offset:     offset,
		Sort:       sort,
	}
	productionPlans, err := c.productionPlanRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.productionPlanRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*Data, 0, len(productionPlans))
	for _, productionPlan := range productionPlans {
		// find custom field value
		customFieldData, err := c.customFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
			EntityType: enum.CustomFieldTypeProductionPlan,
			EntityId:   productionPlan.ID,
			Limit:      1000,
			Offset:     0,
		})
		if err != nil {
			return nil, nil, err
		}

		poCustomFields := c.GetCustomField()
		customFieldMap := make(map[string]string)
		for _, customField := range poCustomFields {
			customFieldMap[customField] = ""
			for _, datum := range customFieldData {
				if datum.Field == customField {
					customFieldMap[customField] = datum.Value
					break
				}
			}
		}

		results = append(results, &Data{
			ProductionPlanData: productionPlan,
			CustomData:         customFieldMap,
		})
	}

	return results, total, nil
}

func (c *productionPlanService) FindProductionPlansWithNoPermission(ctx context.Context, opts *FindProductionPlansOpts, sort *repository.Sort, limit, offset int64) ([]*DataWithNoPermission, *repository.CountResult, error) {
	filter := &repository.SearchProductionPlanOpts{
		IDs:        opts.IDs,
		CustomerID: opts.CustomerID,
		Name:       opts.Name,
		Statuses:   opts.Statuses,
		UserID:     opts.UserID,
		Limit:      limit,
		Offset:     offset,
		Sort:       sort,
	}
	productionPlans, err := c.productionPlanRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.productionPlanRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*DataWithNoPermission, 0, len(productionPlans))
	for _, productionPlan := range productionPlans {
		results = append(results, &DataWithNoPermission{
			ProductionPlanData: productionPlan,
		})
	}

	return results, total, nil
}
