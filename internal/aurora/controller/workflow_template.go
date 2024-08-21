package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/workflow_template"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type WorkflowTemplateController interface {
	CreateWorkflowTemplate(c *gin.Context)
	EditWorkflowTemplate(c *gin.Context)
	DeleteWorkflowTemplate(c *gin.Context)
	FindWorkflowTemplates(c *gin.Context)
}

type workflowTemplateController struct {
	workflowTemplateService workflow_template.Service
}

func (s workflowTemplateController) CreateWorkflowTemplate(c *gin.Context) {
	req := &dto.CreateWorkflowTemplateRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.workflowTemplateService.CreateWorkflowTemplate(c, &workflow_template.CreateWorkflowTemplateOpts{
		Name:       req.Name,
		ConfigData: req.ConfigData,
		Status:     req.Status,
		CreatedBy:  userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateWorkflowTemplateResponse{
		ID: id,
	})
}

func (s workflowTemplateController) EditWorkflowTemplate(c *gin.Context) {
	req := &dto.EditWorkflowTemplateRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	err = s.workflowTemplateService.EditWorkflowTemplate(c, &workflow_template.EditWorkflowTemplateOpts{
		ID:         req.ID,
		Name:       req.Name,
		ConfigData: req.ConfigData,
		Status:     req.Status,
		UpdatedBy:  userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditWorkflowTemplateResponse{})
}

func (s workflowTemplateController) DeleteWorkflowTemplate(c *gin.Context) {
	req := &dto.DeleteWorkflowTemplateRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.workflowTemplateService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteWorkflowTemplateResponse{})
}

func (s workflowTemplateController) FindWorkflowTemplates(c *gin.Context) {
	req := &dto.FindWorkflowTemplatesRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	workflowTemplates, cnt, err := s.workflowTemplateService.FindWorkflowTemplates(c, &workflow_template.FindWorkflowTemplatesOpts{
		Name: req.Filter.Name,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	workflowTemplateResp := make([]*dto.WorkflowTemplate, 0, len(workflowTemplates))
	for _, f := range workflowTemplates {
		workflowTemplateResp = append(workflowTemplateResp, toWorkflowTemplateResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindWorkflowTemplatesResponse{
		WorkflowTemplates: workflowTemplateResp,
		Total:             cnt.Count,
	})
}

func toWorkflowTemplateResp(f *workflow_template.Data) *dto.WorkflowTemplate {
	return &dto.WorkflowTemplate{
		ID:     f.ID,
		Name:   f.Name,
		Config: f.Config,
		Status: f.Status,
	}
}

func RegisterWorkflowTemplateController(
	r *gin.RouterGroup,
	workflowTemplateService workflow_template.Service,
) {
	g := r.Group("workflow-template")

	var c WorkflowTemplateController = &workflowTemplateController{
		workflowTemplateService: workflowTemplateService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateWorkflowTemplate,
		&dto.CreateWorkflowTemplateRequest{},
		&dto.CreateWorkflowTemplateResponse{},
		"Create workflowTemplate",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditWorkflowTemplate,
		&dto.EditWorkflowTemplateRequest{},
		&dto.EditWorkflowTemplateResponse{},
		"Edit workflowTemplate",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteWorkflowTemplate,
		&dto.DeleteWorkflowTemplateRequest{},
		&dto.DeleteWorkflowTemplateResponse{},
		"delete workflowTemplate",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindWorkflowTemplates,
		&dto.FindWorkflowTemplatesRequest{},
		&dto.FindWorkflowTemplatesResponse{},
		"Find workflowTemplates",
	)
}
