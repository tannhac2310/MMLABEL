package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/stage"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type StageController interface {
	CreateStage(c *gin.Context)
	EditStage(c *gin.Context)
	DeleteStage(c *gin.Context)
	FindStages(c *gin.Context)
}

type stageController struct {
	stageService stage.Service
}

func (s stageController) CreateStage(c *gin.Context) {
	req := &dto.CreateStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.stageService.CreateStage(c, &stage.CreateStageOpts{
		Name:           req.Name,
		ParentID:       req.ParentID,
		DepartmentCode: req.DepartmentCode,
		ShortName:      req.ShortName,
		Code:           req.Code,
		Sorting:        req.Sorting,
		ErrorCode:      req.ErrorCode,
		Data:           req.Data,
		Note:           req.Note,
		Status:         req.Status,
		CreatedBy:      userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateStageResponse{
		ID: id,
	})
}

func (s stageController) EditStage(c *gin.Context) {
	req := &dto.EditStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.stageService.EditStage(c, &stage.EditStageOpts{
		ID:             req.ID,
		Name:           req.Name,
		ParentID:       req.ParentID,
		DepartmentCode: req.DepartmentCode,
		ShortName:      req.ShortName,
		Code:           req.Code,
		Sorting:        req.Sorting,
		ErrorCode:      req.ErrorCode,
		Data:           req.Data,
		Note:           req.Note,
		Status:         req.Status,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditStageResponse{})
}

func (s stageController) DeleteStage(c *gin.Context) {
	req := &dto.DeleteStageRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.stageService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteStageResponse{})
}

func (s stageController) FindStages(c *gin.Context) {
	req := &dto.FindStagesRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userId := interceptor.UserIDFromCtx(c)
	if interceptor.IsAdmin(c) {
		userId = ""
	}
	stages, cnt, err := s.stageService.FindStages(c, &stage.FindStagesOpts{
		Name:   req.Filter.Name,
		Code:   req.Filter.Code,
		IDs:    req.Filter.IDs,
		Codes:  req.Filter.Codes,
		UserID: userId,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	stageResp := make([]*dto.Stage, 0, len(stages))
	for _, f := range stages {
		stageResp = append(stageResp, toStageResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindStagesResponse{
		Stages: stageResp,
		Total:  cnt.Count,
	})
}

func toStageResp(f *stage.Data) *dto.Stage {
	return &dto.Stage{
		ID:             f.ID,
		Name:           f.Name,
		ParentID:       f.ParentID.String,
		DepartmentCode: f.DepartmentCode.String,
		ShortName:      f.ShortName,
		Code:           f.Code,
		Sorting:        f.Sorting,
		ErrorCode:      f.ErrorCode.String,
		Data:           f.Data,
		Note:           f.Note.String,
		Status:         f.Status,
	}
}

func RegisterStageController(
	r *gin.RouterGroup,
	stageService stage.Service,
) {
	g := r.Group("stage")

	var c StageController = &stageController{
		stageService: stageService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateStage,
		&dto.CreateStageRequest{},
		&dto.CreateStageResponse{},
		"Create stage",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditStage,
		&dto.EditStageRequest{},
		&dto.EditStageResponse{},
		"Edit stage",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteStage,
		&dto.DeleteStageRequest{},
		&dto.DeleteStageResponse{},
		"delete stage",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindStages,
		&dto.FindStagesRequest{},
		&dto.FindStagesResponse{},
		"Find stages",
	)
}
