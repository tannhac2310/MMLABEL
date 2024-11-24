package order

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

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
	CreateOrder(ctx context.Context, orderWithItems *CreateOrder) error
	UpdateOrder(ctx context.Context, orderWithItems *UpdateOrder) error
	GetOrderById(ctx context.Context, id string) (*model.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	SearchOrders(ctx context.Context, opts *repository.SearchOrderOpts) ([]*model.Order, error)
	CountOrders(ctx context.Context, opts *repository.SearchOrderOpts) (*repository.CountResult, error)
}

type orderService struct {
	orderRepo     repository.OrderRepo
	orderItemRepo repository.OrderItemRepo
}

func (s *orderService) UpdateOrder(ctx context.Context, orderWithItems *UpdateOrder) error {
	//TODO implement me
	panic("implement me")
}

func (s *orderService) GetOrderById(ctx context.Context, id string) (*model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *orderService) SearchOrders(ctx context.Context, opts *repository.SearchOrderOpts) ([]*model.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (s *orderService) CountOrders(ctx context.Context, opts *repository.SearchOrderOpts) (*repository.CountResult, error) {
	//TODO implement me
	panic("implement me")
}

func NewOrderService(orderRepo repository.OrderRepo, orderItemRepo repository.OrderItemRepo) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, orderWithItems *CreateOrder) error {
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
		err := s.orderRepo.Insert(tx, order)
		if err != nil {
			return err
		}

		for _, item := range orderWithItems.Items {
			orderItem := &model.OrderItem{
				ID: item.ID,

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
		return fmt.Errorf("tạo đơn hàng thất bại: %w", errTx)
	}

	return nil
}
