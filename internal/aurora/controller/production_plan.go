package controller

import (
	"strings"

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
	ProcessProductionOrder(c *gin.Context)
	UpdateCustomFields(c *gin.Context)
	UpdateCurrentStage(c *gin.Context)
	SummaryProductionPlan(c *gin.Context)
	UpdateWorkflow(c *gin.Context)
}

type productionPlanController struct {
	productionPlanService production_plan.Service
}

func (s productionPlanController) UpdateCurrentStage(c *gin.Context) {
	req := &dto.UpdateCurrentStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionPlanService.UpdateCurrentStage(c, req.ProductionPlanID, req.CurrentStage)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.UpdateCurrentStageResponse{})
}

func (s productionPlanController) UpdateWorkflow(c *gin.Context) {
	req := &dto.UpdateWorkflowRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionPlanService.UpdateWorkflow(c, req.ProductionPlanID, req.Workflows)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.UpdateWorkflowResponse{})
}

func (s productionPlanController) UpdateCustomFields(c *gin.Context) {
	req := &dto.UpdateCustomFieldPLValuesRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	values := make([]*production_plan.CustomField, 0)
	for _, v := range req.CustomField {
		values = append(values, &production_plan.CustomField{
			Field: v.Key,
			Value: v.Value,
		})
	}

	err = s.productionPlanService.UpdateCustomFields(c, req.ProductionPlanID, values)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.UpdateCustomFieldPLValuesResponse{})
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
		Name:         req.Name,
		ProductName:  req.ProductName,
		ProductCode:  req.ProductCode,
		QtyPaper:     req.QtyPaper,
		QtyFinished:  req.QtyFinished,
		QtyDelivered: req.QtyDelivered,
		Thumbnail:    req.Thumbnail,
		Workflow:     req.Workflow,
		Note:         req.Note,
		CustomField:  customField,
		CreatedBy:    userID,
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
		ID:           req.ID,
		Name:         req.Name,
		ProductName:  req.ProductName,
		ProductCode:  req.ProductCode,
		QtyPaper:     req.QtyPaper,
		QtyFinished:  req.QtyFinished,
		QtyDelivered: req.QtyDelivered,
		Status:       req.Status,
		Note:         req.Note,
		CustomField:  customField,
		Workflow:     req.Workflow,
		CreatedBy:    userID,
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
		IDs: req.Filter.IDs,
		//CustomerID:  req.Filter.CustomerID,
		Search:      strings.TrimSpace(req.Filter.Search),
		Name:        strings.TrimSpace(req.Filter.Name),
		ProductName: strings.TrimSpace(req.Filter.ProductName),
		ProductCode: strings.TrimSpace(req.Filter.ProductCode),
		Statuses:    req.Filter.Statuses,
		//UserID:      interceptor.UserIDFromCtx(c), // TODO add later
		Stage: req.Filter.Stage,
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productionPlanResp := make([]*dto.ProductionPlan, 0, len(productionPlans))
	for _, f := range productionPlans {
		data := &dto.ProductionPlan{
			ID:                f.ID,
			ProductionOrderID: f.ProductionOrderID.String,
			ProductName:       f.ProductName,
			ProductCode:       f.ProductCode,
			QtyPaper:          f.QtyPaper,
			QtyFinished:       f.QtyFinished,
			QtyDelivered:      f.QtyDelivered,
			Thumbnail:         f.Thumbnail.String,
			Status:            f.Status,
			Note:              f.Note.String,
			Workflow:          f.Workflow,
			CreatedBy:         f.CreatedBy,
			CreatedAt:         f.CreatedAt,
			UpdatedBy:         f.UpdatedBy,
			UpdatedAt:         f.UpdatedAt,
			Name:              f.Name,
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
		By:    "created_at",
	}
	if req.Sort != nil {
		sort = &repository.Sort{
			Order: repository.SortOrder(req.Sort.Order),
			By:    req.Sort.By,
		}
	}
	productionPlans, cnt, err := s.productionPlanService.FindProductionPlans(c, &production_plan.FindProductionPlansOpts{
		IDs: req.Filter.IDs,
		//CustomerID:  req.Filter.CustomerID,
		Search:      strings.TrimSpace(req.Filter.Search),
		Name:        strings.TrimSpace(req.Filter.Name),
		ProductName: strings.TrimSpace(req.Filter.ProductName),
		ProductCode: strings.TrimSpace(req.Filter.ProductCode),
		Statuses:    req.Filter.Statuses,
		//UserID:      interceptor.UserIDFromCtx(c),
		Stage: req.Filter.Stage,
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	productionPlanResp := make([]*dto.ProductionPlan, 0, len(productionPlans))
	for _, f := range productionPlans {
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
		userFields := make(map[string][]*dto.UserField)
		for _, uf := range f.UserFields {
			for _, v := range uf {
				userFields[v.Field] = append(userFields[v.Field], &dto.UserField{
					Key:   v.Field,
					Value: v.Value,
				})
			}
		}
		data := &dto.ProductionPlan{
			ID:                f.ID,
			ProductionOrderID: f.ProductionOrderID.String,
			ProductName:       f.ProductName,
			ProductCode:       f.ProductCode,
			QtyPaper:          f.QtyPaper,
			QtyFinished:       f.QtyFinished,
			QtyDelivered:      f.QtyDelivered,
			Thumbnail:         f.Thumbnail.String,
			Status:            f.Status,
			CurrentStage:      f.CurrentStage,
			Note:              f.Note.String,
			Workflow:          f.Workflow,
			CreatedBy:         f.CreatedBy,
			CreatedAt:         f.CreatedAt,
			UpdatedBy:         f.UpdatedBy,
			UpdatedAt:         f.UpdatedAt,
			CreatedByName:     f.CreatedByName,
			UpdatedByName:     f.UpdatedByName,
			Name:              f.Name,
			CustomData:        f.CustomData,
			CustomerData:      customerData,
			UserFields:        userFields,
		}

		productionPlanResp = append(productionPlanResp, data)
	}

	transportutil.SendJSONResponse(c, &dto.FindProductionPlansResponse{
		ProductionPlans: productionPlanResp,
		Total:           cnt.Count,
	})
}

func (s productionPlanController) ProcessProductionOrder(c *gin.Context) {
	req := &dto.ProcessProductionOrderRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	orderStages := make([]*production_plan.ProductionOrderStage, 0)
	for idx, stage := range req.Stages {
		orderStages = append(orderStages, &production_plan.ProductionOrderStage{
			StageID:             stage.StageID,
			EstimatedStartAt:    stage.EstimatedStartAt,
			EstimatedCompleteAt: stage.EstimatedCompleteAt,
			SoLuong:             stage.SoLuong,
			StartedAt:           stage.StartedAt,
			CompletedAt:         stage.CompletedAt,
			Status:              stage.Status,
			Condition:           stage.Condition,
			Note:                stage.Note,
			Data:                stage.Data,
			Sorting:             int16(len(req.Stages) - idx),
		})
	}

	id, err := s.productionPlanService.ProcessProductionOrder(c, &production_plan.ProcessProductionOrderOpts{
		ID:                  req.ID,
		LxsCode:             req.LxsCode,
		Stages:              orderStages,
		EstimatedStartAt:    req.EstimatedStartAt,
		EstimatedCompleteAt: req.EstimatedCompleteAt,
		Data:                req.Data,
		CreatedBy:           userID,
		OrderID:             req.OrderID,
		Note:                req.Note,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.ProcessProductionOrderResponse{
		ID: id,
	})
}

func (s productionPlanController) SummaryProductionPlan(c *gin.Context) {
	req := &dto.SummaryProductionPlanRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	summary, err := s.productionPlanService.SummaryProductionPlans(c, &production_plan.SummaryProductionPlanOpts{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	var total int64
	summaryResponse := make([]*dto.SummaryProductionPlanItem, 0, len(summary))
	for _, f := range summary {
		data := &dto.SummaryProductionPlanItem{
			Stage:  f.Stage,
			Status: f.Status,
			Count:  f.Count,
		}

		summaryResponse = append(summaryResponse, data)
		total += f.Count
	}

	transportutil.SendJSONResponse(c, &dto.SummaryProductionPlanResponse{
		Items: summaryResponse,
		Total: total,
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
		"Create production plan",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditProductionPlan,
		&dto.EditProductionPlanRequest{},
		&dto.EditProductionPlanResponse{},
		"Edit production plan",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteProductionPlan,
		&dto.DeleteProductionPlanRequest{},
		&dto.DeleteProductionPlanResponse{},
		"delete production plan",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindProductionPlans,
		&dto.FindProductionPlansRequest{},
		&dto.FindProductionPlansResponse{},
		"Find production plans",
	)

	routeutil.AddEndpoint(
		g,
		"find-with-no-permission",
		c.FindProductionPlansWithNoPermission,
		&dto.FindProductionPlansRequest{},
		&dto.FindProductionPlansResponse{},
		"Find production plans with no permission",
	)

	routeutil.AddEndpoint(
		g,
		"process-production-order",
		c.ProcessProductionOrder,
		&dto.ProcessProductionOrderRequest{},
		&dto.ProcessProductionOrderResponse{},
		"Process production order",
	)

	routeutil.AddEndpoint(
		g,
		"update-custom-field",
		c.UpdateCustomFields,
		&dto.UpdateCustomFieldPLValuesRequest{},
		&dto.UpdateCustomFieldPLValuesResponse{},
		"Update production plan custom fields",
	)

	routeutil.AddEndpoint(
		g,
		"update-current-stage",
		c.UpdateCurrentStage,
		&dto.UpdateCurrentStageRequest{},
		&dto.UpdateCurrentStageResponse{},
		"Update current stage",
	)

	routeutil.AddEndpoint(
		g,
		"summary",
		c.SummaryProductionPlan,
		&dto.SummaryProductionPlanRequest{},
		&dto.SummaryProductionPlanResponse{},
		"Summary production plan",
	)

	// UpdateWorkflow
	routeutil.AddEndpoint(
		g,
		"update-workflow",
		c.UpdateWorkflow,
		&dto.UpdateWorkflowRequest{},
		&dto.UpdateWorkflowResponse{},
		"Update workflow",
	)
}
