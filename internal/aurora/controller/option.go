package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/option"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type OptionController interface {
	CreateOption(c *gin.Context)
	EditOption(c *gin.Context)
	DeleteOption(c *gin.Context)
	FindOptions(c *gin.Context)
}

type optionController struct {
	optionService option.Service
}

func (s optionController) CreateOption(c *gin.Context) {
	req := &dto.CreateOptionRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.optionService.Create(c, &option.CreateOptionOpts{
		Name:      req.Name,
		Entity:    req.Entity,
		Code:      req.Code,
		Data:      req.Data,
		Status:    req.Status,
		CreatedBy: userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateOptionResponse{
		ID: id,
	})
}

func (s optionController) EditOption(c *gin.Context) {
	req := &dto.EditOptionRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.optionService.Edit(c, &option.EditOptionOpts{
		ID:       req.ID,
		Name:     req.Name,
		Code:     req.Code,
		Data:     req.Data,
		Status:   req.Status,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditOptionResponse{})
}

func (s optionController) DeleteOption(c *gin.Context) {
	req := &dto.DeleteOptionRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.optionService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteOptionResponse{})
}

func (s optionController) FindOptions(c *gin.Context) {
	req := &dto.FindOptionRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	options, cnt, err := s.optionService.Find(c, &option.FindOptionOpts{
		Name: req.Filter.Name,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	optionResp := make([]*dto.Option, 0, len(options))
	for _, f := range options {
		optionResp = append(optionResp, toOptionResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindOptionResponse{
		Options: optionResp,
		Total:   cnt.Count,
	})
}

func toOptionResp(f *option.OptionData) *dto.Option {
	return &dto.Option{
		ID:        f.ID,
		Name:      f.Name,
		Code:      f.Code,
		Data:      f.Data,
		Status:    f.Status,
		CreatedBy: f.CreatedBy,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func RegisterOptionController(
	r *gin.RouterGroup,
	optionService option.Service,
) {
	g := r.Group("option")

	var c OptionController = &optionController{
		optionService: optionService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateOption,
		&dto.CreateOptionRequest{},
		&dto.CreateOptionResponse{},
		"Create option",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditOption,
		&dto.EditOptionRequest{},
		&dto.EditOptionResponse{},
		"Edit option",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteOption,
		&dto.DeleteOptionRequest{},
		&dto.DeleteOptionResponse{},
		"delete option",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindOptions,
		&dto.FindOptionRequest{},
		&dto.FindOptionResponse{},
		"Find options",
	)
}
