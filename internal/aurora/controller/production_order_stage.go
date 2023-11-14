package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/production_order"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type ProductionOrderStageController interface {
	CreateProductionOrderStage(c *gin.Context)
	EditProductionOrderStage(c *gin.Context)
	DeleteProductionOrderStage(c *gin.Context)
}

type productionOrderStageController struct {
	productionOrderStageService production_order.Service
}

func (s productionOrderStageController) CreateProductionOrderStage(c *gin.Context) {
	req := &dto.CreateProductionOrderStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	id, err := s.productionOrderStageService.CreateProductionOrderStage(c, req.ProductionOrderID, &production_order.ProductionOrderStage{
		StageID:             req.StageID,
		EstimatedStartAt:    req.EstimatedStartAt,
		EstimatedCompleteAt: req.EstimatedCompleteAt,
		StartedAt:           req.StartedAt,
		CompletedAt:         req.CompletedAt,
		Status:              req.Status,
		Condition:           req.Condition,
		Note:                req.Note,
		Data:                req.Data,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateProductionOrderStageResponse{
		ID: id,
	})
}

func (s productionOrderStageController) EditProductionOrderStage(c *gin.Context) {
	req := &dto.EditProductionOrderStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionOrderStageService.EditProductionOrderStage(c, req.ProductionOrderID, &production_order.ProductionOrderStage{
		ID:                  req.ID,
		StageID:             req.StageID,
		EstimatedStartAt:    req.EstimatedStartAt,
		EstimatedCompleteAt: req.EstimatedCompleteAt,
		StartedAt:           req.StartedAt,
		CompletedAt:         req.CompletedAt,
		Status:              req.Status,
		Condition:           req.Condition,
		Note:                req.Note,
		Data:                req.Data,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditProductionOrderStageResponse{})
}

func (s productionOrderStageController) DeleteProductionOrderStage(c *gin.Context) {
	req := &dto.DeleteProductionOrderStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionOrderStageService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteProductionOrderStageResponse{})
}

func RegisterProductionOrderStageController(
	r *gin.RouterGroup,
	productionOrderStageService production_order.Service,
) {
	g := r.Group("production-order-stage")

	var c ProductionOrderStageController = &productionOrderStageController{
		productionOrderStageService: productionOrderStageService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateProductionOrderStage,
		&dto.CreateProductionOrderStageRequest{},
		&dto.CreateProductionOrderStageResponse{},
		"Create productionOrderStage",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditProductionOrderStage,
		&dto.EditProductionOrderStageRequest{},
		&dto.EditProductionOrderStageResponse{},
		"Edit productionOrderStage",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteProductionOrderStage,
		&dto.DeleteProductionOrderStageRequest{},
		&dto.DeleteProductionOrderStageResponse{},
		"delete productionOrderStage",
	)
}
