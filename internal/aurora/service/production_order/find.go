package production_order

import (
	"context"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	model2 "mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Analysis struct {
	Status enum.ProductionOrderStatus `json:"status"`
	Count  int64                      `json:"count"`
}

func (c *productionOrderService) FindProductionOrders(ctx context.Context, opts *FindProductionOrdersOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, []*Analysis, error) {
	filter := &repository.SearchProductionOrdersOpts{
		IDs:             opts.IDs,
		CustomerID:      opts.CustomerID,
		ProductCode:     opts.ProductCode,
		ProductName:     opts.ProductName,
		Name:            opts.Name,
		PlannedDateFrom: opts.PlannedDateFrom,
		PlannedDateTo:   opts.PlannedDateTo,
		Status:          opts.Status,
		Responsible:     opts.Responsible,
		StageIDs:        opts.StageIDs,
		Limit:           limit,
		Offset:          offset,
		Sort:            sort,
	}
	productionOrders, err := c.productionOrderRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}

	total, err := c.productionOrderRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}

	results := make([]*Data, 0, len(productionOrders))
	for _, productionOrder := range productionOrders {
		//find stage
		stages, err := c.productionOrderStageRepo.Search(ctx, &repository.SearchProductionOrderStagesOpts{
			ProductionOrderID: productionOrder.ID,
			Limit:             1000,
			Offset:            0,
		})
		if err != nil {
			return nil, nil, nil, err
		}
		stageData := make([]*ProductionOrderStageData, 0, len(stages))
		// find production order stage device for each stage
		for _, stage := range stages {
			stageDevices, err := c.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
				ProductionOrderStageID: stage.ID,
				Limit:                  1000,
				Offset:                 0,
			})
			// consolidate responsibility information from user table
			users := make([]*model2.User, 0)
			for _, stageDevice := range stageDevices {
				if stageDevice.Responsible == nil {
					continue
				}
				for _, responsible := range stageDevice.Responsible {
					user, err := c.userRepo.FindByID(ctx, responsible)
					if err != nil {
						return nil, nil, nil, err
					}
					users = append(users, user)
				}
				stageDevice.ResponsibleObject = users
			}
			if err != nil {
				return nil, nil, nil, err
			}
			stageData = append(stageData, &ProductionOrderStageData{
				ProductionOrderStage:       stage,
				ProductionOrderStageDevice: stageDevices,
			})
		}

		// find custom field value
		customFieldData, err := c.customFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
			EntityType: enum.CustomFieldTypeProductionOrder,
			EntityId:   productionOrder.ID,
			Limit:      1000,
			Offset:     0,
		})

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
			ProductionOrderData:  productionOrder,
			ProductionOrderStage: stageData,
			CustomData:           customFieldMap,
		})
	}

	// analysis
	analysis := make([]*Analysis, 0)
	analysisData, err := c.productionOrderRepo.Analysis(ctx, filter)
	for _, status := range analysisData {
		analysis = append(analysis, &Analysis{
			Status: status.Status,
			Count:  status.Count,
		})
	}

	return results, total, analysis, nil
}

type FindProductionOrdersOpts struct {
	IDs             []string
	CustomerID      string
	ProductName     string
	Name            string
	ProductCode     string
	Status          enum.ProductionOrderStatus
	PlannedDateFrom time.Time
	PlannedDateTo   time.Time
	Responsible     []string
	StageIDs        []string
}
