package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/device"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type DeviceController interface {
	CreateDevice(c *gin.Context)
	EditDevice(c *gin.Context)
	DeleteDevice(c *gin.Context)
	FindDevices(c *gin.Context)
}

type deviceController struct {
	deviceService device.Service
}

func (s deviceController) CreateDevice(c *gin.Context) {
	req := &dto.CreateDeviceRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	userID := interceptor.UserIDFromCtx(c)

	id, err := s.deviceService.CreateDevice(c, &device.CreateDeviceOpts{
		Name:      req.Name,
		Code:      req.Code,
		OptionID:  req.OptionID,
		Data:      req.Data,
		Status:    req.Status,
		CreatedBy: userID,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.CreateDeviceResponse{
		ID: id,
	})
}

func (s deviceController) EditDevice(c *gin.Context) {
	req := &dto.EditDeviceRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.deviceService.EditDevice(c, &device.EditDeviceOpts{
		ID:       req.ID,
		Name:     req.Name,
		Code:     req.Code,
		OptionID: req.OptionID,
		Data:     req.Data,
		Status:   req.Status,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.EditDeviceResponse{})
}

func (s deviceController) DeleteDevice(c *gin.Context) {
	req := &dto.DeleteDeviceRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.deviceService.Delete(c, req.ID)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeleteDeviceResponse{})
}

func (s deviceController) FindDevices(c *gin.Context) {
	req := &dto.FindDevicesRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	devices, cnt, err := s.deviceService.FindDevices(c, &device.FindDevicesOpts{
		Name: req.Filter.Name,
	}, &repository.Sort{
		Order: repository.SortOrderASC,
		By:    "ID",
	}, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	deviceResp := make([]*dto.Device, 0, len(devices))
	for _, f := range devices {
		deviceResp = append(deviceResp, toDeviceResp(f))
	}

	transportutil.SendJSONResponse(c, &dto.FindDevicesResponse{
		Devices: deviceResp,
		Total:   cnt.Count,
	})
}

func toDeviceResp(f *device.Data) *dto.Device {
	return &dto.Device{
		ID:        f.ID,
		Name:      f.Name,
		Code:      f.Code,
		OptionID:  f.OptionID.String,
		Data:      f.Data,
		Status:    f.Status,
		CreatedBy: f.CreatedBy,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func RegisterDeviceController(
	r *gin.RouterGroup,
	deviceService device.Service,
) {
	g := r.Group("device")

	var c DeviceController = &deviceController{
		deviceService: deviceService,
	}

	routeutil.AddEndpoint(
		g,
		"create",
		c.CreateDevice,
		&dto.CreateDeviceRequest{},
		&dto.CreateDeviceResponse{},
		"Create device",
	)

	routeutil.AddEndpoint(
		g,
		"edit",
		c.EditDevice,
		&dto.EditDeviceRequest{},
		&dto.EditDeviceResponse{},
		"Edit device",
	)

	routeutil.AddEndpoint(
		g,
		"delete",
		c.DeleteDevice,
		&dto.DeleteDeviceRequest{},
		&dto.DeleteDeviceResponse{},
		"delete device",
	)

	routeutil.AddEndpoint(
		g,
		"find",
		c.FindDevices,
		&dto.FindDevicesRequest{},
		&dto.FindDevicesResponse{},
		"Find devices",
	)
}
