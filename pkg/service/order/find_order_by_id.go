package order

import (
	"context"
	"fmt"
)

func (b *orderService) FindOrderByID(ctx context.Context, id string) (*Data, error) {
	orders, _, err := b.FindOrders(ctx, &FindOrdersOpts{
		IDs: []string{id},
	}, 1, 0)

	if err != nil {
		return nil, err
	}
	if len(orders) != 1 {
		return nil, fmt.Errorf("order.Search:FindOrderByID not found")
	}

	return orders[0], nil
}
