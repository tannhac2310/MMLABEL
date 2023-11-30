package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/device_config"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type DeviceConfigController interface {
	CreateDeviceConfig(c *gin.Context)
	EditDeviceConfig(c *gin.Context)
	DeleteDeviceConfig(c *gin.Context)
	FindDeviceConfigs(c *gin.Context)
}

type deviceConfigController struct {
	deviceConfigService device_config.Service
}

func (s deviceConfigController) CreateDeviceConfig(c *gin.Context) {
	req := &dto.CreateDeviceConfigRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	userId := interceptor.UserIDFromCtx(c)
	id, err := s.deviceConfigService.CreateDeviceConfig(c, &device_config.CreateDeviceConfigOpts{
		ProductionOrderID: req.ProductionOrderID,
		DeviceID:          req.DeviceID,
		DeviceConfig:      req.DeviceConfig,
		Color:             req.Color,
		Description:       req.Description,
		Search:            req.Search,
		CreatedBy:         userId,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateDeviceConfigResponse{
		ID: id,
	})
}

func (s deviceConfigController) EditDeviceConfig(c *gin.Context) {
	req := &dto.EditDeviceConfigRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.deviceConfigService.EditDeviceConfig(c, &device_config.EditDeviceConfigOpts{
		ID:                req.ID,
		ProductionOrderID: req.ProductionOrderID,
		DeviceID:          req.DeviceID,
		DeviceConfig:      req.DeviceConfig,
		Color:             req.Color,
		Description:       req.Description,
		Search:            req.Search,
		UpdatedBy:         interceptor.UserIDFromCtx(c),
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditDeviceConfigResponse{})
}

func (s deviceConfigController) DeleteDeviceConfig(c *gin.Context) {
	req := &dto.DeleteDeviceConfigRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.deviceConfigService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteDeviceConfigResponse{})
}

func (s deviceConfigController) FindDeviceConfigs(c *gin.Context) {
	req := &dto.FindDeviceConfigsRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	deviceConfigs, cnt, err := s.deviceConfigService.FindDeviceConfigs(c, &device_config.FindDeviceConfigsOpts{
		Search: req.Filter.Search,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	deviceConfigResp := make([]*dto.DeviceConfig, 0, len(deviceConfigs))
	for _, f := range deviceConfigs {
		deviceConfigResp = append(deviceConfigResp, toDeviceConfigResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindDeviceConfigsResponse{
		DeviceConfigs: deviceConfigResp,
		Total:         cnt.Count,
	})
}

func toDeviceConfigResp(f *device_config.Data) *dto.DeviceConfig {
	return &dto.DeviceConfig{
		ID:                  f.ID,
		ProductionOrderID:   f.ProductionOrderID,
		ProductionOrderName: f.ProductionOrderName,
		DeviceID:            f.DeviceID.String,
		DeviceConfig:        f.DeviceConfig,
		Color:               f.Color.String,
		Description:         f.Description.String,
		CreatedBy:           f.CreatedBy,
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
	}
}

func RegisterDeviceConfigController(
	r *gin.RouterGroup,
	deviceConfigService device_config.Service,
) {
	g := r.Group("device-config")

	var c DeviceConfigController = &deviceConfigController{
		deviceConfigService: deviceConfigService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateDeviceConfig,
		&dto.CreateDeviceConfigRequest{},
		&dto.CreateDeviceConfigResponse{},
		"Create deviceConfig",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditDeviceConfig,
		&dto.EditDeviceConfigRequest{},
		&dto.EditDeviceConfigResponse{},
		"Edit deviceConfig",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteDeviceConfig,
		&dto.DeleteDeviceConfigRequest{},
		&dto.DeleteDeviceConfigResponse{},
		"delete deviceConfig",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindDeviceConfigs,
		&dto.FindDeviceConfigsRequest{},
		&dto.FindDeviceConfigsResponse{},
		"Find deviceConfigs",
	)
}
