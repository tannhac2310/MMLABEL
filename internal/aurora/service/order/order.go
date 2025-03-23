package order

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type OrderWithItems struct {
	Order OrderData
	Items []*OrderItemData
}

type OrderData struct {
	ID                 string
	Title              string
	MaDatHangMm        string
	MaHopDongKhachHang string
	MaHopDong          string
	SaleName           string
	SaleAdminName      string
	Status             enum.OrderStatus
}

type OrderItemData struct {
	ID                      string
	ProductionPlanProductID string
	ProductionPlanID        string
	ProductionQuantity      int64
	Quantity                int64
	UnitPrice               float64
	DeliveredQuantity       int64
	EstimatedDeliveryDate   time.Time
	DeliveredDate           time.Time
	Status                  string
	Attachment              map[string]any
	Note                    string
}

type CreateOrder struct {
	Order    OrderData
	Items    []*OrderItemData
	CreateBy string
}

type UpdateOrder struct {
	Order    OrderData
	Items    []*OrderItemData
	UpdateBy string
}

type OrderService interface {
	CreateOrder(ctx context.Context, orderWithItems *CreateOrder) (string, error)
	UpdateOrder(ctx context.Context, orderWithItems *UpdateOrder) error
	DeleteOrder(ctx context.Context, id string) error
	SearchOrders(ctx context.Context, opts *repository.SearchOrderOpts) ([]*OrderWithItems, *repository.CountResult, error)
	UpdateOrderStatus(ctx context.Context, id string, status enum.OrderStatus) error
}

type orderService struct {
	orderRepo     repository.OrderRepo
	orderItemRepo repository.OrderItemRepo
}

func (s *orderService) UpdateOrder(ctx context.Context, orderWithItems *UpdateOrder) error {
	errTx := cockroach.ExecInTx(ctx, func(tx context.Context) error {
		orderTable := model.Order{}
		oderUpdater := cockroach.NewUpdater(orderTable.TableName(), model.DepartmentFieldID, orderWithItems.Order.ID)

		oderUpdater.Set(model.OrderFieldTitle, orderWithItems.Order.Title)
		oderUpdater.Set(model.OrderFieldMaDatHangMm, orderWithItems.Order.MaDatHangMm)
		oderUpdater.Set(model.OrderFieldMaHopDongKhachHang, orderWithItems.Order.MaHopDongKhachHang)
		oderUpdater.Set(model.OrderFieldMaHopDong, orderWithItems.Order.MaHopDong)
		oderUpdater.Set(model.OrderFieldSaleName, orderWithItems.Order.SaleName)
		oderUpdater.Set(model.OrderFieldSaleAdminName, orderWithItems.Order.SaleAdminName)
		oderUpdater.Set(model.OrderFieldStatus, orderWithItems.Order.Status)
		oderUpdater.Set(model.OrderFieldUpdatedBy, orderWithItems.UpdateBy)
		oderUpdater.Set(model.OrderFieldUpdatedAt, time.Now())

		err := cockroach.UpdateFields(ctx, oderUpdater)
		if err != nil {
			return fmt.Errorf("cập nhật đơn hàng thất bại: %w", err)
		}

		// delete all order items
		err = s.orderItemRepo.DeleteByOrderID(tx, orderWithItems.Order.ID)

		// insert new order items
		for _, item := range orderWithItems.Items {
			orderItem := &model.OrderItem{
				ID:                      idutil.ULIDNow(),
				OrderID:                 orderWithItems.Order.ID,
				ProductionPlanProductID: item.ProductionPlanProductID,
				ProductionPlanID:        item.ProductionPlanID,
				ProductionQuantity:      item.ProductionQuantity,
				Quantity:                item.Quantity,
				UnitPrice:               item.UnitPrice,
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
	// 2. count order
	count, err := s.orderRepo.Count(ctx, opts)

	if err != nil {
		return nil, nil, fmt.Errorf("lấy số lượng đơn hàng thất bại: %w", err)
	}

	orderWithItems := make([]*OrderWithItems, 0, len(orders))
	for _, order := range orders {
		items, err := s.orderItemRepo.Search(ctx, &repository.SearchOrderItemOpts{
			OrderID: order.ID,
			Limit:   1000,
			Offset:  0,
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
				ID:                 order.ID,
				Title:              order.Title,
				MaDatHangMm:        order.MaDatHangMm,
				MaHopDongKhachHang: order.MaHopDongKhachHang,
				MaHopDong:          order.MaHopDong,
				SaleName:           order.SaleName.String,
				SaleAdminName:      order.SaleAdminName.String,
				Status:             order.Status,
			},
			Items: orderItemData,
		})
	}
	return orderWithItems, count, nil

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
		now := time.Now()
		order := &model.Order{
			ID:                 orderId,
			Title:              orderWithItems.Order.Title,
			MaDatHangMm:        orderWithItems.Order.MaDatHangMm,
			MaHopDongKhachHang: orderWithItems.Order.MaHopDongKhachHang,
			MaHopDong:          orderWithItems.Order.MaHopDong,
			SaleName:           cockroach.String(orderWithItems.Order.SaleName),
			SaleAdminName:      cockroach.String(orderWithItems.Order.SaleAdminName),
			Status:             orderWithItems.Order.Status,
			CreatedBy:          orderWithItems.CreateBy,
			UpdatedBy:          orderWithItems.CreateBy,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		err = s.orderRepo.Insert(tx, order)
		if err != nil {
			return err
		}

		for _, item := range orderWithItems.Items {
			orderItem := &model.OrderItem{
				ID:                      idutil.ULIDNow(),
				OrderID:                 orderId,
				ProductionPlanProductID: item.ProductionPlanProductID,
				ProductionPlanID:        item.ProductionPlanID,
				ProductionQuantity:      item.ProductionQuantity,
				Quantity:                item.Quantity,
				UnitPrice:               item.UnitPrice,
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

func (s *orderService) UpdateOrderStatus(ctx context.Context, id string, status enum.OrderStatus) error {
	err := s.orderRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("cập nhật trạng thái đơn hàng thất bại: %w", err)
	}
	return nil
}
