package production_order

import (
	"context"
	"fmt"
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
	// find permission stage for user
	permissions, err := c.roleService.FindRolePermissionsByUser(ctx, opts.UserID)

	if err != nil {
		return nil, nil, nil, err
	}
	// find stage for user
	stages := make([]string, 0)
	stages = append(stages, "-1")
	for _, p := range permissions {
		if p.EntityType == enum.PermissionEntityTypeStage {
			stages = append(stages, p.EntityID)
		}
	}
	// intercect stage with opts.StageIDs
	if len(opts.StageIDs) > 0 {
		intersect := make([]string, 0)
		for _, stage := range opts.StageIDs {
			for _, p := range stages {
				if stage == p {
					intersect = append(intersect, stage)
				}
			}
		}
		stages = intersect
	}

	fmt.Println("FindProductionOrders stages", stages)

	filter := &repository.SearchProductionOrdersOpts{
		IDs:                  opts.IDs,
		CustomerID:           opts.CustomerID,
		ProductCode:          opts.ProductCode,
		ProductName:          opts.ProductName,
		Name:                 opts.Name,
		EstimatedStartAtFrom: opts.EstimatedStartAtFrom,
		EstimatedStartAtTo:   opts.EstimatedStartAtTo,
		Status:               opts.Status,
		Statuses:             opts.Statuses,
		Responsible:          opts.Responsible,
		StageIDs:             stages,
		UserID:               opts.UserID,
		OrderStageStatus:     opts.OrderStageStatus,
		StageInLine:          opts.StageInLine, // search lsx mà theo công đoạn StageInLine đang sản xuất: production_start
		DeviceID:             opts.DeviceID,
		Limit:                limit,
		Offset:               offset,
		Sort:                 sort,
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
	idMap := make(map[string]bool)

	for _, productionOrder := range productionOrders {
		if _, ok := idMap[productionOrder.ID]; ok {
			continue
		}
		idMap[productionOrder.ID] = true
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
		stageDevicesOfPO, err := c.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
			ProductionOrderID: productionOrder.ID,
			Limit:             1000,
			Offset:            0,
		})
		mapStageDevices := make(map[string][]*repository.ProductionOrderStageDeviceData)
		for _, stageDevice := range stageDevicesOfPO {
			mapStageDevices[stageDevice.ProductionOrderStageID] = append(mapStageDevices[stageDevice.ProductionOrderStageID], stageDevice)
		}
		// find production order stage device for each stage
		for _, stage := range stages {
			stageDevices := mapStageDevices[stage.ID]
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
	IDs                  []string
	CustomerID           string
	ProductName          string
	Name                 string
	ProductCode          string
	Status               enum.ProductionOrderStatus
	Statuses             []enum.ProductionOrderStatus
	EstimatedStartAtFrom time.Time
	EstimatedStartAtTo   time.Time
	OrderStageStatus     enum.ProductionOrderStageStatus
	Responsible          []string
	StageIDs             []string
	StageInLine          string // search lsx mà theo công đoạn StageInLine đang sản xuất: production_start
	DeviceID             string
	UserID               string
}
