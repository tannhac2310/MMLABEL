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
		if err != nil {
			return nil, nil, err
		}
		results = append(results, &Data{
			ProductionOrderData: productionOrder,
		})
	}
	return results, total, nil
}

type FindProductionOrdersOpts struct {
	IDs         []string
	CustomerID  string
	ProductName string
	ProductCode string
	Status      enum.ProductionOrderStatus
}
