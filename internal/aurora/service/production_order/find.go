package production_order

import (
	"context"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (c *productionOrderService) FindProductionOrders(ctx context.Context, opts *FindProductionOrdersOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	filter := &repository.SearchProductionOrdersOpts{
		IDs:         opts.IDs,
		CustomerID:  opts.CustomerID,
		ProductCode: opts.ProductCode,
		ProductName: opts.ProductName,
		Name:        opts.Name,
		Status:      opts.Status,
		Limit:       limit,
		Offset:      offset,
		Sort:        sort,
	}
	productionOrders, err := c.productionOrderRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.productionOrderRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
		}
		stageData := make([]*ProductionOrderStageData, 0, len(stages))
		// find production order stage device for each stage
		for _, stage := range stages {
			stageDevices, err := c.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
				ProductionOrderStageID: stage.ID,
				Limit:                  1000,
				Offset:                 0,
			})
			if err != nil {
				return nil, nil, err
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
	return results, total, nil
}

type FindProductionOrdersOpts struct {
	IDs         []string
	CustomerID  string
	ProductName string
	Name        string
	ProductCode string
	Status      enum.ProductionOrderStatus
}
