package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/order"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type OrderController interface {
	CreateOrder(ctx *gin.Context)
	UpdateOrder(ctx *gin.Context)
	DeleteOrder(ctx *gin.Context)
	FindOrder(ctx *gin.Context)
}

type orderController struct {
	orderSvc order.OrderService
}

func (o orderController) CreateOrder(ctx *gin.Context) {
	var req dto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	orderParams := order.OrderData{
		Title:              req.Order.Title,
		MaDatHangMm:        req.Order.MaDatHangMm,
		MaHopDongKhachHang: req.Order.MaHopDongKhachHang,
		MaHopDong:          req.Order.MaHopDong,
		SaleName:           req.Order.SaleName,
		SaleAdminName:      req.Order.SaleAdminName,
		Status:             req.Order.Status,
	}

	orderItems := make([]*order.OrderItemData, 0, len(req.Items))
	for _, item := range req.Items {
		orderItems = append(orderItems, &order.OrderItemData{
			ProductionPlanProductID: item.ProductionPlanProductID,
			ProductionPlanID:        item.ProductionPlanID,
			ProductionQuantity:      item.ProductionQuantity,
			Quantity:                item.Quantity,
			UnitPrice:               item.UnitPrice,
			DeliveredQuantity:       item.DeliveredQuantity,
			EstimatedDeliveryDate:   item.EstimatedDeliveryDate,
			DeliveredDate:           item.DeliveredDate,
			Status:                  item.Status,
			Attachment:              item.Attachment,
			Note:                    item.Note,
		})
	}
	id, err := o.orderSvc.CreateOrder(ctx, &order.CreateOrder{
		Order: orderParams,
		Items: orderItems,
	})
	if err != nil {
		transportutil.Error(ctx, err)
		return
	}

	transportutil.SendJSONResponse(ctx, dto.CreateOrderResponse{ID: id})
}

func (o orderController) UpdateOrder(ctx *gin.Context) {
	var req dto.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	oderData := order.OrderData{
		ID:                 req.Order.ID,
		Title:              req.Order.Title,
		MaDatHangMm:        req.Order.MaDatHangMm,
		MaHopDongKhachHang: req.Order.MaHopDongKhachHang,
		MaHopDong:          req.Order.MaHopDong,
		SaleName:           req.Order.SaleName,
		SaleAdminName:      req.Order.SaleAdminName,
		Status:             req.Order.Status,
	}

	orderItems := make([]*order.OrderItemData, 0, len(req.Items))
	for _, item := range req.Items {
		orderItems = append(orderItems, &order.OrderItemData{
			ID:                      item.ID,
			ProductionPlanProductID: item.ProductionPlanProductID,
			ProductionPlanID:        item.ProductionPlanID,
			ProductionQuantity:      item.ProductionQuantity,
			Quantity:                item.Quantity,
			UnitPrice:               item.UnitPrice,
			DeliveredQuantity:       item.DeliveredQuantity,
			EstimatedDeliveryDate:   item.EstimatedDeliveryDate,
			DeliveredDate:           item.DeliveredDate,
			Status:                  item.Status,
			Attachment:              item.Attachment,
			Note:                    item.Note,
		})
	}

	err := o.orderSvc.UpdateOrder(ctx, &order.UpdateOrder{
		Order: oderData,
		Items: orderItems,
	})

	if err != nil {
		transportutil.Error(ctx, err)
		return
	}

	transportutil.SendJSONResponse(ctx, dto.UpdateOrderResponse{ID: req.Order.ID})
}

func (o orderController) DeleteOrder(ctx *gin.Context) {
	var req dto.DeleteOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err := o.orderSvc.DeleteOrder(ctx, req.ID)
	if err != nil {
		transportutil.Error(ctx, err)
		return
	}

	transportutil.SendJSONResponse(ctx, dto.DeleteOrderResponse{ID: req.ID})
}

func (o orderController) FindOrder(ctx *gin.Context) {
	var req dto.SearchOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		transportutil.Error(ctx, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	orders, total, err := o.orderSvc.SearchOrders(ctx, &repository.SearchOrderOpts{
		IDs:                     req.Filter.IDs,
		Search:                  req.Filter.Search,
		ProductionPlanID:        req.Filter.ProductionPlanID,
		ProductionPlanProductID: req.Filter.ProductionPlanProductID,
		Status:                  req.Filter.Status,
		Limit:                   req.Paging.Limit,
		Offset:                  req.Paging.Offset,
		Sort:                    nil,
	})
	if err != nil {
		transportutil.Error(ctx, err)
		return
	}
	orderData := make([]*dto.OrderWithItems, 0, len(orders))
	for _, orderWithItems := range orders {
		items := make([]*dto.OrderItemData, 0, len(orderWithItems.Items))
		for _, item := range orderWithItems.Items {
			items = append(items, &dto.OrderItemData{
				ID:                      item.ID,
				ProductionPlanProductID: item.ProductionPlanProductID,
				ProductionPlanID:        item.ProductionPlanID,
				ProductionQuantity:      item.ProductionQuantity,
				Quantity:                item.Quantity,
				UnitPrice:               item.UnitPrice,
				DeliveredQuantity:       item.DeliveredQuantity,
				EstimatedDeliveryDate:   item.EstimatedDeliveryDate,
				DeliveredDate:           item.DeliveredDate,
				Status:                  item.Status,
				Attachment:              item.Attachment,
				Note:                    item.Note,
			})
		}
		orderData = append(orderData, &dto.OrderWithItems{
			Order: dto.OrderData{
				ID:                 orderWithItems.Order.ID,
				Title:              orderWithItems.Order.Title,
				MaDatHangMm:        orderWithItems.Order.MaDatHangMm,
				MaHopDongKhachHang: orderWithItems.Order.MaHopDongKhachHang,
				MaHopDong:          orderWithItems.Order.MaHopDong,
				SaleName:           orderWithItems.Order.SaleName,
				SaleAdminName:      orderWithItems.Order.SaleAdminName,
				Status:             orderWithItems.Order.Status,
			},
			Items: items,
			ID:    orderWithItems.Order.ID,
			Title: orderWithItems.Order.Title,
		})
	}

	transportutil.SendJSONResponse(ctx, dto.SearchOrderResponse{
		Orders: orderData,
		Total:  total.Count,
	})
}

func RegisterOrderController(
	r *gin.RouterGroup,
	orderSvc order.OrderService,
) {
	g := r.Group("order")
	var c OrderController = &orderController{
		orderSvc: orderSvc,
	}

	routeutil.AddEndpoint(
		g, "create",
		c.CreateOrder,
		&dto.CreateOrderRequest{},
		&dto.CreateOrderResponse{},
		"create order",
	)

	routeutil.AddEndpoint(
		g, "update",
		c.UpdateOrder,
		&dto.UpdateOrderRequest{},
		&dto.UpdateOrderResponse{},
		"update order",
	)

	routeutil.AddEndpoint(
		g, "delete",
		c.DeleteOrder,
		&dto.DeleteOrderRequest{},
		&dto.DeleteOrderResponse{},
		"delete order",
	)

	routeutil.AddEndpoint(
		g, "find",
		c.FindOrder,
		&dto.SearchOrderRequest{},
		&dto.SearchOrderResponse{},
		"find order",
	)
}
