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
	FindProductionOrdersWithNoPermission(c *gin.Context)
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

	customField := make([]*production_order.CustomField, 0)
	for _, field := range req.CustomField {
		customField = append(customField, &production_order.CustomField{
			Field: field.Key,
			Value: field.Value,
		})
	}
	id, err := s.productionOrderService.CreateProductionOrder(c, &production_order.CreateProductionOrderOpts{
		Name:                 req.Name,
		ProductCode:          req.ProductCode,
		ProductName:          req.ProductName,
		CustomerID:           req.CustomerID,
		SalesID:              req.SalesID,
		QtyPaper:             req.QtyPaper,
		QtyFinished:          req.QtyFinished,
		QtyDelivered:         req.QtyDelivered,
		EstimatedStartAt:     req.EstimatedStartAt,
		EstimatedCompleteAt:  req.EstimatedCompleteAt,
		DeliveryDate:         req.DeliveryDate,
		DeliveryImage:        req.DeliveryImage,
		Status:               req.Status,
		Note:                 req.Note,
		ProductionOrderStage: orderStage,
		CustomField:          customField,
		CreatedBy:            userID,
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
	// todo implement later
	//productionOderStage := make([]*production_order.ProductionOrderStage, 0)
	//for idx, stage := range req.ProductionOrderStages {
	//	productionOderStage = append(productionOderStage, &production_order.ProductionOrderStage{
	//		ID:                  stage.ID,
	//		StageID:             stage.StageID,
	//		EstimatedStartAt:    stage.EstimatedStartAt,
	//		EstimatedCompleteAt: stage.EstimatedCompleteAt,
	//		StartedAt:           stage.StartedAt,
	//		CompletedAt:         stage.CompletedAt,
	//		Status:              stage.Status,
	//		Condition:           stage.Condition,
	//		Note:                stage.Note,
	//		Data:                stage.Data,
	//		Sorting:             int16(len(req.ProductionOrderStages) - idx),
	//	})
	//}
	err = s.productionOrderService.EditProductionOrder(c, &production_order.EditProductionOrderOpts{
		ID:                  req.ID,
		Name:                req.Name,
		QtyPaper:            req.QtyPaper,
		QtyFinished:         req.QtyFinished,
		QtyDelivered:        req.QtyDelivered,
		EstimatedStartAt:    req.EstimatedStartAt,
		EstimatedCompleteAt: req.EstimatedCompleteAt,
		Status:              req.Status,
		DeliveryDate:        req.DeliveryDate,
		DeliveryImage:       req.DeliveryImage,
		Note:                req.Note,
		//ProductionOrderStage: productionOderStage,
		//CustomData: nil,
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

func (s productionOrderController) FindProductionOrdersWithNoPermission(c *gin.Context) {
	req := &dto.FindProductionOrdersRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	sort := &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}
	if req.Sort != nil {
		sort = &repository.Sort{
			Order: repository.SortOrder(req.Sort.Order),
			By:    req.Sort.By,
		}
	}
	productionOrders, cnt, err := s.productionOrderService.FindProductionOrdersWithNoPermission(c, &production_order.FindProductionOrdersOpts{
		IDs:                  req.Filter.IDs,
		CustomerID:           req.Filter.CustomerID,
		Name:                 req.Filter.Name,
		Status:               req.Filter.Status,
		Statuses:             req.Filter.Statuses,
		EstimatedStartAtFrom: req.Filter.EstimatedStartAtFrom,
		EstimatedStartAtTo:   req.Filter.EstimatedStartAtTo,
		OrderStageStatus:     req.Filter.OrderStageStatus,
		Responsible:          req.Filter.Responsible,
		StageIDs:             req.Filter.StageIDs,
		StageInLine:          req.Filter.StageInLine, // search lsx mà theo công đoạn StageInLine đang sản xuất: production_start
		DeviceID:             req.Filter.DeviceID,
		UserID:               interceptor.UserIDFromCtx(c),
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	//FindProductionOrdersWithNoPermissionResponse
	productionOrderResp := make([]*dto.ProductionOrder2, 0, len(productionOrders))
	for _, f := range productionOrders {
		data := &dto.ProductionOrder2{
			ID:                  f.ID,
			Name:                f.Name,
			ProductCode:         f.ProductCode,
			ProductName:         f.ProductName,
			CustomerID:          f.CustomerID,
			SalesID:             f.SalesID,
			QtyPaper:            f.QtyPaper,
			QtyFinished:         f.QtyFinished,
			QtyDelivered:        f.QtyDelivered,
			EstimatedStartAt:    f.EstimatedStartAt.Time,
			EstimatedCompleteAt: f.EstimatedCompleteAt.Time,
			DeliveryDate:        f.DeliveryDate,
			DeliveryImage:       f.DeliveryImage.String,
			Status:              f.Status,
			Note:                f.Note.String,
		}
		productionOrderResp = append(productionOrderResp, data)
	}

	transportutil.SendJSONResponse(c, &dto.FindProductionOrdersWithNoPermissionResponse{
		ProductionOrders: productionOrderResp,
		Total:            cnt.Count,
	})
}

func (s productionOrderController) FindProductionOrders(c *gin.Context) {
	req := &dto.FindProductionOrdersRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	sort := &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}
	if req.Sort != nil {
		sort = &repository.Sort{
			Order: repository.SortOrder(req.Sort.Order),
			By:    req.Sort.By,
		}
	}
	productionOrders, cnt, analysis, err := s.productionOrderService.FindProductionOrders(c, &production_order.FindProductionOrdersOpts{
		IDs:                             req.Filter.IDs,
		CustomerID:                      req.Filter.CustomerID,
		Name:                            req.Filter.Name,
		Status:                          req.Filter.Status,
		Statuses:                        req.Filter.Statuses,
		EstimatedStartAtFrom:            req.Filter.EstimatedStartAtFrom,
		EstimatedStartAtTo:              req.Filter.EstimatedStartAtTo,
		EstimatedCompleteAtFrom:         req.Filter.EstimatedCompleteAtFrom,
		EstimatedCompleteAtTo:           req.Filter.EstimatedCompleteAtTo,
		OrderStageStatus:                req.Filter.OrderStageStatus,
		OrderStageEstimatedStartFrom:    req.Filter.OrderStageEstimatedStartFrom,
		OrderStageEstimatedStartTo:      req.Filter.OrderStageEstimatedStartTo,
		OrderStageEstimatedCompleteFrom: req.Filter.OrderStageEstimatedCompleteFrom,
		OrderStageEstimatedCompleteTo:   req.Filter.OrderStageEstimatedCompleteTo,
		Responsible:                     req.Filter.Responsible,
		StageIDs:                        req.Filter.StageIDs,
		StageInLine:                     req.Filter.StageInLine, // search lsx mà theo công đoạn StageInLine đang sản xuất: production_start
		DeviceID:                        req.Filter.DeviceID,
		UserID:                          interceptor.UserIDFromCtx(c),
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productionOrderResp := make([]*dto.ProductionOrder, 0, len(productionOrders))
	for _, f := range productionOrders {
		productionOrderResp = append(productionOrderResp, toProductionOrderResp(f))
	}
	// analysis
	analysisResp := make([]*dto.Analysis, 0, len(analysis))
	for _, a := range analysis {
		analysisResp = append(analysisResp, &dto.Analysis{
			Status: a.Status,
			Count:  a.Count,
		})
	}

	transportutil.SendJSONResponse(c, &dto.FindProductionOrdersResponse{
		ProductionOrders: productionOrderResp,
		Total:            cnt.Count,
		Analysis:         analysisResp,
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
				DeviceData:             device.DeviceData,
				ResponsibleObject:      responsibleObject,
				Quantity:               device.Quantity,
				ProcessStatus:          device.ProcessStatus,
				Status:                 device.Status,
				Responsible:            device.Responsible,
				Settings:               device.Settings,
				Note:                   device.Note.String,
				EstimatedCompleteAt:    device.EstimatedCompleteAt.Time,
				AssignedQuantity:       device.AssignedQuantity,
			})
		}

		orderStage = append(orderStage, &dto.OrderStage{
			ID:                     item.ID,
			ProductionOrderID:      item.ProductionOrderID,
			StageID:                item.StageID,
			StageName:              item.StageName,
			StartedAt:              item.StartedAt.Time,
			CompletedAt:            item.CompletedAt.Time,
			Status:                 item.Status,
			Condition:              item.Condition.String,
			Note:                   item.Note.String,
			Data:                   item.Data,
			CreatedAt:              item.CreatedAt,
			UpdatedAt:              item.UpdatedAt,
			WaitingAt:              item.WaitingAt.Time,
			ReceptionAt:            item.ReceptionAt.Time,
			ProductionStartAt:      item.ProductionStartAt.Time,
			ProductionCompletionAt: item.ProductionCompletionAt.Time,
			ProductDeliveryAt:      item.ProductDeliveryAt.Time,
			EstimatedStartAt:       item.EstimatedStartAt.Time,
			EstimatedCompleteAt:    item.EstimatedCompleteAt.Time,
			Sorting:                item.Sorting,
			OrderStageDevices:      productionOrderStageDevice,
		})
	}
	customerData := &dto.Customer{}
	if f.CustomerData != nil {
		customerData = &dto.Customer{
			ID:                 f.CustomerData.ID,
			Name:               f.CustomerData.Name,
			Tax:                f.CustomerData.Tax.String,
			Code:               f.CustomerData.Code,
			Country:            f.CustomerData.Country,
			Province:           f.CustomerData.Province,
			Address:            f.CustomerData.Address,
			Fax:                f.CustomerData.Fax.String,
			CompanyWebsite:     f.CustomerData.CompanyWebsite.String,
			CompanyPhone:       f.CustomerData.CompanyPhone.String,
			ContactPersonName:  f.CustomerData.ContactPersonName.String,
			ContactPersonEmail: f.CustomerData.ContactPersonEmail.String,
			ContactPersonPhone: f.CustomerData.ContactPersonPhone.String,
			ContactPersonRole:  f.CustomerData.ContactPersonRole.String,
			Note:               f.CustomerData.Note.String,
			Status:             f.CustomerData.Status,
		}
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
		EstimatedStartAt:      f.EstimatedStartAt.Time,
		EstimatedCompleteAt:   f.EstimatedCompleteAt.Time,
		DeliveryDate:          f.DeliveryDate,
		DeliveryImage:         f.DeliveryImage.String,
		Status:                f.Status,
		Note:                  f.Note.String,
		ProductionOrderStages: orderStage,
		CustomData:            f.CustomData,
		CustomerData:          customerData,
		CreatedBy:             f.CreatedBy,
		CreatedAt:             f.CreatedAt,
		UpdatedAt:             f.UpdatedAt,
		Version:               f.Version,
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
		"find-with-no-permission",
		c.FindProductionOrdersWithNoPermission,
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
