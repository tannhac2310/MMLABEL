package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/production_order"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type ProductionOrderController interface {
	CreateProductionOrder(c *gin.Context)
	EditProductionOrder(c *gin.Context)
	DeleteProductionOrder(c *gin.Context)
	FindProductionOrders(c *gin.Context)
}

type productionOrderController struct {
	productionOrderService production_order.Service
}

func (s productionOrderController) CreateProductionOrder(c *gin.Context) {
	req := &dto.CreateProductionOrderRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.productionOrderService.CreateProductionOrder(c, &production_order.CreateProductionOrderOpts{
		ProductCode:           req.ProductCode,
		ProductName:           req.ProductName,
		CustomerID:            req.CustomerID,
		SalesID:               req.SalesID,
		QtyPaper:              req.QtyPaper,
		QtyFinished:           req.QtyFinished,
		QtyDelivered:          req.QtyDelivered,
		PlannedProductionDate: req.PlannedProductionDate,
		DeliveryDate:          req.DeliveryDate,
		DeliveryImage:         req.DeliveryImage,
		Status:                req.Status,
		Note:                  req.Note,
		CreatedBy:             userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateProductionOrderResponse{
		ID: id,
	})
}

func (s productionOrderController) EditProductionOrder(c *gin.Context) {
	req := &dto.EditProductionOrderRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionOrderService.EditProductionOrder(c, &production_order.EditProductionOrderOpts{
		ID:                    req.ID,
		QtyPaper:              req.QtyPaper,
		QtyFinished:           req.QtyFinished,
		QtyDelivered:          req.QtyDelivered,
		PlannedProductionDate: req.PlannedProductionDate,
		Status:                req.Status,
		DeliveryDate:          req.DeliveryDate,
		DeliveryImage:         req.DeliveryImage,
		Note:                  req.Note,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditProductionOrderResponse{})
}

func (s productionOrderController) DeleteProductionOrder(c *gin.Context) {
	req := &dto.DeleteProductionOrderRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionOrderService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteProductionOrderResponse{})
}

func (s productionOrderController) FindProductionOrders(c *gin.Context) {
	req := &dto.FindProductionOrdersRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	productionOrders, cnt, err := s.productionOrderService.FindProductionOrders(c, &production_order.FindProductionOrdersOpts{
		IDs:         req.Filter.IDs,
		CustomerID:  req.Filter.CustomerID,
		ProductName: req.Filter.ProductName,
		ProductCode: req.Filter.ProductCode,
		Status:      req.Filter.Status,
	}, &repository.Sort{
		Order: repository.SortOrderASC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productionOrderResp := make([]*dto.ProductionOrder, 0, len(productionOrders))
	for _, f := range productionOrders {
		productionOrderResp = append(productionOrderResp, toProductionOrderResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindProductionOrdersResponse{
		ProductionOrders: productionOrderResp,
		Total:            cnt.Count,
	})
}

func toProductionOrderResp(f *production_order.Data) *dto.ProductionOrder {
	return &dto.ProductionOrder{
		ID:                    f.ID,
		ProductCode:           f.ProductCode,
		ProductName:           f.ProductName,
		CustomerID:            f.CustomerID,
		SalesID:               f.SalesID,
		QtyPaper:              f.QtyPaper,
		QtyFinished:           f.QtyFinished,
		QtyDelivered:          f.QtyDelivered,
		PlannedProductionDate: f.PlannedProductionDate,
		DeliveryDate:          f.DeliveryDate,
		DeliveryImage:         f.DeliveryImage.String,
		Status:                f.Status,
		Note:                  f.Note.String,
		CreatedBy:             f.CreatedBy,
		CreatedAt:             f.CreatedAt,
		UpdatedAt:             f.UpdatedAt,
	}
}

func RegisterProductionOrderController(
	r *gin.RouterGroup,
	productionOrderService production_order.Service,
) {
	g := r.Group("production-order")

	var c ProductionOrderController = &productionOrderController{
		productionOrderService: productionOrderService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateProductionOrder,
		&dto.CreateProductionOrderRequest{},
		&dto.CreateProductionOrderResponse{},
		"Create productionOrder",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditProductionOrder,
		&dto.EditProductionOrderRequest{},
		&dto.EditProductionOrderResponse{},
		"Edit productionOrder",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteProductionOrder,
		&dto.DeleteProductionOrderRequest{},
		&dto.DeleteProductionOrderResponse{},
		"delete productionOrder",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindProductionOrders,
		&dto.FindProductionOrdersRequest{},
		&dto.FindProductionOrdersResponse{},
		"Find productionOrders",
	)
}
