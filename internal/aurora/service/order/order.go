package order

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type OrderWithItems struct {
	Order OrderData
	Items []*OrderItemData
}

type OrderData struct {
	ID                  string
	Title               string
	Code                string
	SaleName            string
	SaleAdminName       string
	ProductCode         string
	ProductName         string
	CustomerID          string
	CustomerProductCode string
	CustomerProductName string
	Status              string
}

type OrderItemData struct {
	ID                      string
	ProductionPlanProductID string
	ProductionPlanID        string
	ProductionQuantity      int64
	Quantity                int64
	UnitPrice               float64
	TotalAmount             float64
	DeliveredQuantity       int64
	EstimatedDeliveryDate   time.Time
	DeliveredDate           time.Time
	Status                  string
	Attachment              map[string]string
	Note                    string
}

type CreateOrder struct {
	Order OrderData
	Items []*OrderItemData
}

type UpdateOrder struct {
	Order OrderData
	Items []*OrderItemData
}

type OrderService interface {
	CreateOrder(ctx context.Context, orderWithItems *CreateOrder) (string, error)
	UpdateOrder(ctx context.Context, orderWithItems *UpdateOrder) error
	DeleteOrder(ctx context.Context, id string) error
	SearchOrders(ctx context.Context, opts *repository.SearchOrderOpts) ([]*OrderWithItems, *repository.CountResult, error)
}

type orderService struct {
	orderRepo     repository.OrderRepo
	orderItemRepo repository.OrderItemRepo
}

func (s *orderService) UpdateOrder(ctx context.Context, orderWithItems *UpdateOrder) error {
	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		order := &model.Order{
			ID:                  orderWithItems.Order.ID,
			Title:               orderWithItems.Order.Title,
			Code:                orderWithItems.Order.Code,
			SaleName:            orderWithItems.Order.SaleName,
			SaleAdminName:       orderWithItems.Order.SaleAdminName,
			ProductCode:         orderWithItems.Order.ProductCode,
			ProductName:         orderWithItems.Order.ProductName,
			CustomerID:          orderWithItems.Order.CustomerID,
			CustomerProductCode: orderWithItems.Order.CustomerProductCode,
			CustomerProductName: orderWithItems.Order.CustomerProductName,
			Status:              orderWithItems.Order.Status,
		}

		err := s.orderRepo.Update(tx, order)
		if err != nil {
			return err
		}

		// delete all order items
		err = s.orderItemRepo.DeleteByOrderID(tx, order.ID)

		// insert new order items
		for _, item := range orderWithItems.Items {
			orderItem := &model.OrderItem{
				ID:                      item.ID,
				OrderID:                 orderWithItems.Order.ID,
				ProductionPlanProductID: item.ProductionPlanProductID,
				ProductionPlanID:        item.ProductionPlanID,
				ProductionQuantity:      item.ProductionQuantity,
				Quantity:                item.Quantity,
				UnitPrice:               item.UnitPrice,
				TotalAmount:             item.TotalAmount,
				DeliveredQuantity:       item.DeliveredQuantity,
				EstimatedDeliveryDate:   cockroach.Time(item.EstimatedDeliveryDate),
				DeliveredDate:           cockroach.Time(item.DeliveredDate),
				Status:                  item.Status,
				Attachment:              item.Attachment,
				Note:                    item.Note,
			}
			err := s.orderItemRepo.Insert(tx, orderItem)
			if err != nil {
				return err
			}
		}
		return nil

	})

	if errTx != nil {
		return fmt.Errorf("cập nhật đơn hàng thất bại: %w", errTx)
	}

	return nil
}

func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		err := s.orderRepo.SoftDelete(tx, id)
		if err != nil {
			return err
		}

		err = s.orderItemRepo.DeleteByOrderID(tx, id)
		if err != nil {
			return err
		}

		return nil
	})

	if errTx != nil {
		return fmt.Errorf("xóa đơn hàng thất bại: %w", errTx)
	}

	return nil
}

func (s *orderService) SearchOrders(ctx context.Context, opts *repository.SearchOrderOpts) ([]*OrderWithItems, *repository.CountResult, error) {
	orders, err := s.orderRepo.Search(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("lấy danh sách đơn hàng thất bại: %w", err)
	}

	if len(orders) > 0 {
		orderWithItems := make([]*OrderWithItems, 0, len(orders))
		for _, order := range orders {
			items, err := s.orderItemRepo.Search(ctx, &repository.SearchOrderItemOpts{
				OrderID: order.ID,
			})
			if err != nil {
				return nil, nil, fmt.Errorf("lấy danh sách sản phẩm trong đơn hàng thất bại: %w", err)
			}
			orderItemData := make([]*OrderItemData, 0, len(items))
			for _, item := range items {
				orderItemData = append(orderItemData, &OrderItemData{
					ID:                      item.ID,
					ProductionPlanProductID: item.ProductionPlanProductID,
					ProductionPlanID:        item.ProductionPlanID,
					ProductionQuantity:      item.ProductionQuantity,
					Quantity:                item.Quantity,
					UnitPrice:               item.UnitPrice,
					TotalAmount:             item.TotalAmount,
					DeliveredQuantity:       item.DeliveredQuantity,
					EstimatedDeliveryDate:   item.EstimatedDeliveryDate.Time,
					DeliveredDate:           item.DeliveredDate.Time,
					Status:                  item.Status,
					Attachment:              item.Attachment,
					Note:                    item.Note,
				})
			}
			orderWithItems = append(orderWithItems, &OrderWithItems{
				Order: OrderData{
					ID:                  order.ID,
					Title:               order.Title,
					Code:                order.Code,
					SaleName:            order.SaleName,
					SaleAdminName:       order.SaleAdminName,
					ProductCode:         order.ProductCode,
					ProductName:         order.ProductName,
					CustomerID:          order.CustomerID,
					CustomerProductCode: order.CustomerProductCode,
					CustomerProductName: order.CustomerProductName,
					Status:              order.Status,
				},
				Items: orderItemData,
			})
		}
		return orderWithItems, nil, nil
	}

	return nil, nil, nil

}

func NewOrderService(orderRepo repository.OrderRepo, orderItemRepo repository.OrderItemRepo) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, orderWithItems *CreateOrder) (string, error) {
	orderId := idutil.ULIDNow()
	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		cntRow, err := s.orderRepo.CntRows(tx)
		if err != nil {
			return err
		}
		orderId = fmt.Sprintf("order-%d", cntRow+1)

		order := &model.Order{
			ID:                  orderId,
			Title:               orderWithItems.Order.Title,
			Code:                orderWithItems.Order.Code,
			SaleName:            orderWithItems.Order.SaleName,
			SaleAdminName:       orderWithItems.Order.SaleAdminName,
			ProductCode:         orderWithItems.Order.ProductCode,
			ProductName:         orderWithItems.Order.ProductName,
			CustomerID:          orderWithItems.Order.CustomerID,
			CustomerProductCode: orderWithItems.Order.CustomerProductCode,
			CustomerProductName: orderWithItems.Order.CustomerProductName,
			Status:              orderWithItems.Order.Status,
		}

		err = s.orderRepo.Insert(tx, order)
		if err != nil {
			return err
		}

		for _, item := range orderWithItems.Items {
			orderItem := &model.OrderItem{
				ID:                      item.ID,
				OrderID:                 orderId,
				ProductionPlanProductID: item.ProductionPlanProductID,
				ProductionPlanID:        item.ProductionPlanID,
				ProductionQuantity:      item.ProductionQuantity,
				Quantity:                item.Quantity,
				UnitPrice:               item.UnitPrice,
				TotalAmount:             item.TotalAmount,
				DeliveredQuantity:       item.DeliveredQuantity,
				EstimatedDeliveryDate:   cockroach.Time(item.EstimatedDeliveryDate),
				DeliveredDate:           cockroach.Time(item.DeliveredDate),
				Status:                  item.Status,
				Attachment:              item.Attachment,
				Note:                    item.Note,
			}
			err := s.orderItemRepo.Insert(tx, orderItem)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTx != nil {
		return "", fmt.Errorf("tạo đơn hàng thất bại: %w", errTx)
	}

	return orderId, nil
}
