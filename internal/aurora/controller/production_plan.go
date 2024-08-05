package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/production_plan"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/generic"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type ProductionPlanController interface {
	CreateProductionPlan(c *gin.Context)
	EditProductionPlan(c *gin.Context)
	DeleteProductionPlan(c *gin.Context)
	FindProductionPlans(c *gin.Context)
	FindProductionPlansWithNoPermission(c *gin.Context)
}

type productionPlanController struct {
	productionPlanService production_plan.Service
}

func (s productionPlanController) CreateProductionPlan(c *gin.Context) {
	req := &dto.CreateProductionPlanRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	customField := make([]*production_plan.CustomField, 0)
	for _, field := range req.CustomField {
		customField = append(customField, &production_plan.CustomField{
			Field: field.Key,
			Value: field.Value,
		})
	}
	id, err := s.productionPlanService.CreateProductionPlan(c, &production_plan.CreateProductionPlanOpts{
		Name:        req.Name,
		CustomerID:  req.CustomerID,
		SalesID:     req.SalesID,
		Status:      req.Status,
		Note:        req.Note,
		CustomField: customField,
		CreatedBy:   userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateProductionPlanResponse{
		ID: id,
	})
}

func (s productionPlanController) EditProductionPlan(c *gin.Context) {
	req := &dto.EditProductionPlanRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	customField := generic.Map(req.CustomField, func(f *dto.CustomField) *production_plan.CustomField {
		return &production_plan.CustomField{
			Field: f.Key,
			Value: f.Value,
		}
	})

	err = s.productionPlanService.EditProductionPlan(c, &production_plan.EditProductionPlanOpts{
		ID:          req.ID,
		Name:        req.Name,
		CustomerID:  req.CustomerID,
		SalesID:     req.SalesID,
		Status:      req.Status,
		Note:        req.Note,
		CustomField: customField,
		CreatedBy:   userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditProductionPlanRequest{})
}

func (s productionPlanController) DeleteProductionPlan(c *gin.Context) {
	req := &dto.DeleteProductionPlanRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionPlanService.DeleteProductionPlan(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteProductionPlanResponse{})
}

func (s productionPlanController) FindProductionPlansWithNoPermission(c *gin.Context) {
	req := &dto.FindProductionPlansRequest{}
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
	productionPlans, cnt, err := s.productionPlanService.FindProductionPlansWithNoPermission(c, &production_plan.FindProductionPlansOpts{
		IDs:        req.Filter.IDs,
		CustomerID: req.Filter.CustomerID,
		Name:       req.Filter.Name,
		Statuses:   req.Filter.Statuses,
		UserID:     interceptor.UserIDFromCtx(c),
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productionPlanResp := make([]*dto.ProductionPlan, 0, len(productionPlans))
	for _, f := range productionPlans {
		data := &dto.ProductionPlan{
			ID:         f.ID,
			CustomerID: f.CustomerID,
			SalesID:    f.SalesID,
			Thumbnail:  f.Thumbnail.String,
			Status:     f.Status,
			Note:       f.Note.String,
			CreatedBy:  f.CreatedBy,
			CreatedAt:  f.CreatedAt,
			UpdatedBy:  f.UpdatedBy,
			UpdatedAt:  f.UpdatedAt,
			DeletedAt:  f.DeletedAt.Time,
			Name:       f.Name,
		}
		productionPlanResp = append(productionPlanResp, data)
	}

	transportutil.SendJSONResponse(c, &dto.FindProductionPlansResponse{
		ProductionPlans: productionPlanResp,
		Total:           cnt.Count,
	})
}

func (s productionPlanController) FindProductionPlans(c *gin.Context) {
	req := &dto.FindProductionPlansRequest{}
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
	productionPlans, cnt, err := s.productionPlanService.FindProductionPlans(c, &production_plan.FindProductionPlansOpts{
		IDs:        req.Filter.IDs,
		CustomerID: req.Filter.CustomerID,
		Name:       req.Filter.Name,
		Statuses:   req.Filter.Statuses,
		UserID:     interceptor.UserIDFromCtx(c),
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productionPlanResp := make([]*dto.ProductionPlan, 0, len(productionPlans))
	for _, f := range productionPlans {
		data := &dto.ProductionPlan{
			ID:         f.ID,
			CustomerID: f.CustomerID,
			SalesID:    f.SalesID,
			Thumbnail:  f.Thumbnail.String,
			Status:     f.Status,
			Note:       f.Note.String,
			CreatedBy:  f.CreatedBy,
			CreatedAt:  f.CreatedAt,
			UpdatedBy:  f.UpdatedBy,
			UpdatedAt:  f.UpdatedAt,
			DeletedAt:  f.DeletedAt.Time,
			Name:       f.Name,
			CustomData: f.CustomData,
		}

		productionPlanResp = append(productionPlanResp, data)
	}

	transportutil.SendJSONResponse(c, &dto.FindProductionPlansResponse{
		ProductionPlans: productionPlanResp,
		Total:           cnt.Count,
	})
}

func RegisterProductionPlanController(
	r *gin.RouterGroup,
	productionPlanService production_plan.Service,
) {
	g := r.Group("production-plan")

	var c ProductionPlanController = &productionPlanController{
		productionPlanService: productionPlanService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateProductionPlan,
		&dto.CreateProductionPlanRequest{},
		&dto.CreateProductionPlanResponse{},
		"Create productionPlan",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditProductionPlan,
		&dto.EditProductionPlanRequest{},
		&dto.EditProductionPlanResponse{},
		"Edit productionPlan",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteProductionPlan,
		&dto.DeleteProductionPlanRequest{},
		&dto.DeleteProductionPlanResponse{},
		"delete productionPlan",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindProductionPlans,
		&dto.FindProductionPlansRequest{},
		&dto.FindProductionPlansResponse{},
		"Find productionPlans",
	)

	routeutil.AddEndpoint(
		g,
		"find-with-no-permission",
		c.FindProductionPlansWithNoPermission,
		&dto.FindProductionPlansRequest{},
		&dto.FindProductionPlansResponse{},
		"Find productionPlans",
	)
}
