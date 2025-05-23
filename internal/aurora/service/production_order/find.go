package production_order

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	model2 "mmlabel.gitlab.com/mm-printing-backend/pkg/model"
	repository2 "mmlabel.gitlab.com/mm-printing-backend/pkg/repository"

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
	// Tập hợp giao stages và opts.StageIDs
	if len(opts.StageIDs) > 0 {
		intersect := make([]string, 0)
		for _, stage := range stages { // loop user can access stages
			for _, p := range opts.StageIDs { // loop stage in opts.StageIDs
				if stage == p { // if stage in opts.StageIDs is in user can access stages
					intersect = append(intersect, stage)
				}
			}
		}
		stages = intersect
	}

	fmt.Println("FindProductionOrders stages", stages)

	filter := &repository.SearchProductionOrdersOpts{
		IDs:                             opts.IDs,
		ProductionPlanIDs:               opts.ProductionPlanIDs,
		CustomerID:                      opts.CustomerID,
		ProductCode:                     opts.ProductCode,
		ProductName:                     opts.ProductName,
		Name:                            opts.Name,
		EstimatedStartAtFrom:            opts.EstimatedStartAtFrom,
		EstimatedStartAtTo:              opts.EstimatedStartAtTo,
		EstimatedCompletedFrom:          opts.EstimatedCompleteAtFrom,
		EstimatedCompletedTo:            opts.EstimatedCompleteAtTo,
		Status:                          opts.Status,
		Statuses:                        opts.Statuses,
		OrderStageStatus:                opts.OrderStageStatus,
		OrderStageEstimatedStartFrom:    opts.OrderStageEstimatedStartFrom,
		OrderStageEstimatedStartTo:      opts.OrderStageEstimatedStartTo,
		OrderStageEstimatedCompleteFrom: opts.OrderStageEstimatedCompleteFrom,
		OrderStageEstimatedCompleteTo:   opts.OrderStageEstimatedCompleteTo,
		Responsible:                     opts.Responsible,
		StageIDs:                        stages,
		StageInLine:                     opts.StageInLine, // search lsx mà theo công đoạn StageInLine đang sản xuất: production_start
		UserID:                          opts.UserID,
		DeviceID:                        opts.DeviceID,
		Limit:                           limit,
		Offset:                          offset,
		Sort:                            sort,
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
	customerIds := make([]string, 0, len(productionOrders))
	userIds := make([]string, 0)
	productionPlanIDs := make([]string, 0)
	orderIDs := make([]string, 0)

	for _, productionOrder := range productionOrders {
		customerIds = append(customerIds, productionOrder.CustomerID)
		userIds = append(userIds, productionOrder.CreatedBy)
		if productionOrder.ProductionPlanID.String != "" {
			productionPlanIDs = append(productionPlanIDs, productionOrder.ProductionPlanID.String)
		}

		if productionOrder.OrderID.String != "" {
			orderIDs = append(orderIDs, productionOrder.OrderID.String)
		}
	}
	// find customer name
	customerData, err := c.customerRepo.Search(ctx, &repository.SearchCustomerOpts{
		IDs:    customerIds,
		Limit:  int64(len(customerIds)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("c.customerRepo.FindByIDs: %w", err)
	}

	customerMap := make(map[string]*repository.CustomerData)
	for _, customer := range customerData {
		customerMap[customer.ID] = customer
	}

	// find users
	userData, err := c.userRepo.Search(ctx, &repository2.SearchUsersOpts{
		IDs:    userIds,
		Limit:  int64(len(userIds)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("c.userRepo.FindByIDs: %w", err)
	}

	userMap := make(map[string]*model2.User)
	for _, user := range userData {
		userMap[user.ID] = user.User
	}

	// find production plan
	productionPlanData, err := c.productionPlanRepo.Search(ctx, &repository.SearchProductionPlanOpts{
		IDs:    productionPlanIDs,
		Limit:  int64(len(productionPlanIDs)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("c.productionPlanRepo.FindByIDs: %w", err)
	}

	productionPlanMap := make(map[string]*repository.ProductionPlanData)
	for _, productionPlan := range productionPlanData {
		productionPlanMap[productionPlan.ID] = productionPlan
	}

	// find order
	orderData, err := c.orderRepo.Search(ctx, &repository.SearchOrderOpts{
		IDs:    orderIDs,
		Limit:  int64(len(orderIDs)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("c.orderRepo.FindByIDs: %w", err)
	}

	orderMap := make(map[string]*repository.OrderData)
	for _, order := range orderData {
		orderMap[order.ID] = order
	}

	// find stage for each production order
	allStages, err := c.stageRepo.Search(ctx, &repository.SearchStagesOpts{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	stageMap := make(map[string]string)
	for _, stage := range allStages {
		stageMap[stage.ID] = stage.Name
	}
	for _, productionOrder := range productionOrders {
		if _, ok := idMap[productionOrder.ID]; ok {
			continue
		}
		idMap[productionOrder.ID] = true
		//find stage
		wf, err := c.productionOrderStageRepo.Search(ctx, &repository.SearchProductionOrderStagesOpts{
			ProductionOrderID: productionOrder.ID,
			Limit:             1000,
			Offset:            0,
		})
		if err != nil {
			return nil, nil, nil, err
		}
		stageData := make([]*ProductionOrderStageData, 0, len(wf))
		stageDevicesOfPO, err := c.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
			ProductionOrderIDs: []string{productionOrder.ID},
			Limit:              1000,
			Offset:             0,
		})
		mapStageDevices := make(map[string][]*repository.ProductionOrderStageDeviceData)
		for _, stageDevice := range stageDevicesOfPO {
			mapStageDevices[stageDevice.ProductionOrderStageID] = append(mapStageDevices[stageDevice.ProductionOrderStageID], stageDevice)
		}
		// find production order stage device for each stage
		for _, stage := range wf {
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
				StageName:                  stageMap[stage.StageID],
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

		//poCustomFields := c.GetCustomField()
		customFieldMap := make(map[string]string)
		//for _, customField := range poCustomFields {
		//	customFieldMap[customField] = ""
		for _, datum := range customFieldData {
			//if datum.Field == customField {
			if _, ok := customFieldMap[datum.Field]; ok {
				customFieldMap[datum.Field] = customFieldMap[datum.Field] + "," + datum.Value // nếu trùng field thì nối chuỗi. todo: this is not good. need to improve
			} else {
				customFieldMap[datum.Field] = datum.Value
			}
			//break
			//}
		}
		//}

		results = append(results, &Data{
			ProductionOrderData:  productionOrder,
			ProductionOrderStage: stageData,
			CustomData:           customFieldMap,
			CustomerData:         customerMap[productionOrder.CustomerID],
			CreatedByName:        userMap[productionOrder.CreatedBy].Name,
			ProductionPlanData:   productionPlanMap[productionOrder.ProductionPlanID.String],
			OrderData:            orderMap[productionOrder.OrderID.String],
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

func (c *productionOrderService) FindProductionOrdersWithNoPermission(ctx context.Context, opts *FindProductionOrdersOpts, sort *repository.Sort, limit, offset int64) ([]*DataWithNoPermission, *repository.CountResult, error) {
	filter := &repository.SearchProductionOrdersOpts{
		IDs:                  opts.IDs,
		CustomerID:           opts.CustomerID,
		ProductCode:          opts.ProductCode,
		ProductName:          opts.ProductName,
		ProductionPlanIDs:    opts.ProductionPlanIDs,
		Name:                 opts.Name,
		EstimatedStartAtFrom: opts.EstimatedStartAtFrom,
		EstimatedStartAtTo:   opts.EstimatedStartAtTo,
		Status:               opts.Status,
		Statuses:             opts.Statuses,
		Responsible:          opts.Responsible,
		StageIDs:             opts.StageIDs,
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
		return nil, nil, err
	}

	total, err := c.productionOrderRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*DataWithNoPermission, 0, len(productionOrders))

	for _, productionOrder := range productionOrders {
		results = append(results, &DataWithNoPermission{
			ProductionOrderData: productionOrder,
		})
	}

	return results, total, nil
}

type FindProductionOrdersOpts struct {
	IDs                             []string
	CustomerID                      string
	ProductName                     string
	Name                            string
	ProductCode                     string
	Status                          enum.ProductionOrderStatus
	Statuses                        []enum.ProductionOrderStatus
	ProductionPlanIDs               []string
	EstimatedStartAtFrom            time.Time
	EstimatedStartAtTo              time.Time
	EstimatedCompleteAtFrom         time.Time
	EstimatedCompleteAtTo           time.Time
	OrderStageStatus                enum.ProductionOrderStageStatus
	OrderStageEstimatedStartFrom    time.Time
	OrderStageEstimatedStartTo      time.Time
	OrderStageEstimatedCompleteFrom time.Time
	OrderStageEstimatedCompleteTo   time.Time
	Responsible                     []string
	StageIDs                        []string
	StageInLine                     string // search lsx mà theo công đoạn StageInLine đang sản xuất: production_start
	DeviceID                        string
	UserID                          string
}
