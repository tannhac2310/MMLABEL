package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/production_order_stage_device"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type ProductionOrderStageDeviceController interface {
	CreateProductionOrderStageDevice(c *gin.Context)
	EditProductionOrderStageDevice(c *gin.Context)
	DeleteProductionOrderStageDevice(c *gin.Context)
}

type productionOrderStageDeviceController struct {
	productionOrderStageDeviceService production_order_stage_device.Service
}

func (s productionOrderStageDeviceController) CreateProductionOrderStageDevice(c *gin.Context) {
	req := &dto.CreateProductionOrderStageDeviceRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	id, err := s.productionOrderStageDeviceService.Create(c, &production_order_stage_device.CreateProductionOrderStageDeviceOpts{
		ProductionOrderStageID: req.ProductionOrderStageID,
		DeviceID:               req.DeviceID,
		Quantity:               req.Quantity,
		ProcessStatus:          req.ProcessStatus,
		Status:                 req.Status,
		Responsible:            req.Responsible,
	})

	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateProductionOrderStageDeviceResponse{
		ID: id,
	})
}

func (s productionOrderStageDeviceController) EditProductionOrderStageDevice(c *gin.Context) {
	req := &dto.EditProductionOrderStageDeviceRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionOrderStageDeviceService.Edit(c, &production_order_stage_device.EditProductionOrderStageDeviceOpts{
		ID:            req.ID,
		DeviceID:      req.DeviceID,
		Quantity:      req.Quantity,
		ProcessStatus: req.ProcessStatus,
		Status:        req.Status,
		Responsible:   req.Responsible,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditProductionOrderStageDeviceResponse{})
}

func (s productionOrderStageDeviceController) DeleteProductionOrderStageDevice(c *gin.Context) {
	req := &dto.DeleteProductionOrderStageDeviceRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionOrderStageDeviceService.Deletes(c, req.IDs)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteProductionOrderStageDeviceResponse{})
}

func RegisterProductionOrderStageDeviceController(
	r *gin.RouterGroup,
	productionOrderStageDeviceService production_order_stage_device.Service,
) {
	g := r.Group("production-order-stage-device")

	var c ProductionOrderStageDeviceController = &productionOrderStageDeviceController{
		productionOrderStageDeviceService: productionOrderStageDeviceService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateProductionOrderStageDevice,
		&dto.CreateProductionOrderStageDeviceRequest{},
		&dto.CreateProductionOrderStageDeviceResponse{},
		"Create productionOrderStageDevice",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditProductionOrderStageDevice,
		&dto.EditProductionOrderStageDeviceRequest{},
		&dto.EditProductionOrderStageDeviceResponse{},
		"Edit productionOrderStageDevice",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteProductionOrderStageDevice,
		&dto.DeleteProductionOrderStageDeviceRequest{},
		&dto.DeleteProductionOrderStageDeviceResponse{},
		"delete productionOrderStageDevice",
	)
}
