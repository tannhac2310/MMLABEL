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
	AcceptAndChangeNextStage(c *gin.Context)
}

type productionOrderController struct {
	productionOrderService production_order.Service
}

func (s productionOrderController) AcceptAndChangeNextStage(c *gin.Context) {
	req := &dto.AcceptAndChangeNextStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	err = s.productionOrderService.AcceptAndChangeNextStage(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.AcceptAndChangeNextStageResponse{})
}
func (s productionOrderController) CreateProductionOrder(c *gin.Context) {
	req := &dto.CreateProductionOrderRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)
	orderStage := make([]*production_order.ProductionOrderStage, 0)

	for idx, stage := range req.ProductionOrderStages {
		orderStage = append(orderStage, &production_order.ProductionOrderStage{
			StageID:             stage.StageID,
			EstimatedStartAt:    stage.EstimatedStartAt,
			EstimatedCompleteAt: stage.EstimatedCompleteAt,
			StartedAt:           stage.StartedAt,
			CompletedAt:         stage.CompletedAt,
			Status:              stage.Status,
			Condition:           stage.Condition,
			Note:                stage.Note,
			Data:                stage.Data,
			Sorting:             int16(len(req.ProductionOrderStages) - idx),
		})
	}

	id, err := s.productionOrderService.CreateProductionOrder(c, &production_order.CreateProductionOrderOpts{
		Name:                  req.Name,
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
		ProductionOrderStage:  orderStage,
		CustomData:            nil,
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
	// write code to edit production order and production order stage
	productionOderStage := make([]*production_order.ProductionOrderStage, 0)
	for idx, stage := range req.ProductionOrderStages {
		productionOderStage = append(productionOderStage, &production_order.ProductionOrderStage{
			ID:                  stage.ID,
			StageID:             stage.StageID,
			EstimatedStartAt:    stage.EstimatedStartAt,
			EstimatedCompleteAt: stage.EstimatedCompleteAt,
			StartedAt:           stage.StartedAt,
			CompletedAt:         stage.CompletedAt,
			Status:              stage.Status,
			Condition:           stage.Condition,
			Note:                stage.Note,
			Data:                stage.Data,
			Sorting:             int16(len(req.ProductionOrderStages) - idx),
		})
	}
	err = s.productionOrderService.EditProductionOrder(c, &production_order.EditProductionOrderOpts{
		ID:                    req.ID,
		Name:                  req.Name,
		QtyPaper:              req.QtyPaper,
		QtyFinished:           req.QtyFinished,
		QtyDelivered:          req.QtyDelivered,
		PlannedProductionDate: req.PlannedProductionDate,
		Status:                req.Status,
		DeliveryDate:          req.DeliveryDate,
		DeliveryImage:         req.DeliveryImage,
		Note:                  req.Note,
		ProductionOrderStage:  productionOderStage,
		//CustomData:            req.CustomData,
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

	err = s.productionOrderService.DeleteProductionOrder(c, req.ID)
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
	orderStage := make([]*dto.OrderStage, 0)

	for _, item := range f.ProductionOrderStage {
		productionOrderStageDevice := make([]*dto.OrderStageDevice, 0)

		for _, device := range item.ProductionOrderStageDevice {
			responsibleObject := make([]*dto.User, 0)
			for _, r := range device.ResponsibleObject {
				responsibleObject = append(responsibleObject, &dto.User{
					ID:      r.ID,
					Name:    r.Name,
					Avatar:  r.Avatar,
					Address: r.Address,
				})
			}

			productionOrderStageDevice = append(productionOrderStageDevice, &dto.OrderStageDevice{
				ID:                     device.ID,
				ProductionOrderStageID: device.ProductionOrderStageID,
				DeviceID:               device.DeviceID,
				DeviceName:             device.DeviceName,
				ResponsibleObject:      responsibleObject,
				Quantity:               device.Quantity,
				ProcessStatus:          device.ProcessStatus,
				Status:                 device.Status,
				Responsible:            device.Responsible,
				Settings:               device.Settings,
				Note:                   device.Note.String,
			})
		}
		orderStage = append(orderStage, &dto.OrderStage{
			ID:                  item.ID,
			StageID:             item.StageID,
			EstimatedStartAt:    item.EstimatedStartAt.Time,
			EstimatedCompleteAt: item.EstimatedCompleteAt.Time,
			StartedAt:           item.StartedAt.Time,
			CompletedAt:         item.CompletedAt.Time,
			Status:              item.Status,
			Condition:           item.Condition.String,
			Note:                item.Note.String,
			Data:                item.Data,
			OrderStageDevices:   productionOrderStageDevice,
			Sorting:             item.Sorting,
		})
	}
	return &dto.ProductionOrder{
		ID:                    f.ID,
		Name:                  f.Name,
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
		ProductionOrderStages: orderStage,
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
	routeutil.AddEndpoint(
		g,
		"accept-and-change-next-stage",
		c.AcceptAndChangeNextStage,
		&dto.AcceptAndChangeNextStageRequest{},
		&dto.AcceptAndChangeNextStageResponse{},
		"Accept and change next stage",
	)
}
