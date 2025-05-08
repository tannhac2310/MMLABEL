package production_order

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type Analysis struct {
	Status enum.ProductionOrderStatus `json:"status"`
	Count  int64                      `json:"count"`
}

func (c *productionOrderService) FindAnalysis(ctx context.Context, opts *FindProductionOrdersOpts, sort *repository.Sort, limit, offset int64) ([]*Analysis, *repository.CountResult, error) {
	// find permission stage for user
	permissions, err := c.roleService.FindRolePermissionsByUser(ctx, opts.UserID)

	if err != nil {
		return nil, nil, err
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
	// analysis
	analysis := make([]*Analysis, 0)
	analysisData, err := c.productionOrderRepo.Analysis(ctx, filter)
	for _, status := range analysisData {
		analysis = append(analysis, &Analysis{
			Status: status.Status,
			Count:  status.Count,
		})
	}
	if err != nil {
		return nil, nil, err
	}

	// count (if required)
	total, err := c.productionOrderRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return analysis, total, err
}
