package controller

import (
	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/interceptor"

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
	FindEventLog(c *gin.Context)
	FindProcessDeviceHistory(c *gin.Context)
	UpdateProcessDeviceHistoryIsSolved(c *gin.Context)
	FindAvailabilityTime(c *gin.Context)
	FindWorkingDevice(c *gin.Context)
}

type productionOrderStageDeviceController struct {
	productionOrderStageDeviceService production_order_stage_device.Service
}

func (s productionOrderStageDeviceController) FindWorkingDevice(c *gin.Context) {
	data, err := s.productionOrderStageDeviceService.Find(c, &production_order_stage_device.FindProductionOrderStageDeviceOpts{
		ProcessStatus: enum.ProductionOrderStageDeviceStatusStart,
		Limit:         10000,
		Offset:        0,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	result := make([]*dto.ProductionOrderStageDevice, 0, len(data))

	for _, d := range data {
		result = append(result, &dto.ProductionOrderStageDevice{
			ID:                     d.ID,
			ProductionOrderStageID: d.ProductionOrderStageID,
			DeviceID:               d.DeviceID,
			Quantity:               d.Quantity,
			ProcessStatus:          d.ProcessStatus,
			Status:                 d.Status,
			Responsible:            d.Responsible,
			Settings:               d.Settings,
			Note:                   d.Note.String,
		})
	}

	transportutil.SendJSONResponse(c, &dto.FindProductionOrderStageDevicesResponse{
		ProductionOrderStageDevices: result,
		Total:                       int64(len(result)),
	})
}
func (s productionOrderStageDeviceController) FindProcessDeviceHistory(c *gin.Context) {
	req := &dto.FindDeviceStatusHistoryRequest{}
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

	deviceProcessStatusHistoryData, total, err := s.productionOrderStageDeviceService.FindProcessDeviceHistory(c, &production_order_stage_device.FindProcessDeviceHistoryOpts{
		ProcessStatus: req.Filter.ProcessStatus,
		DeviceID:      req.Filter.DeviceID,
		IsResolved:    req.Filter.IsResolved,
		ErrorCodes:    req.Filter.ErrorCodes,
		CreatedFrom:   req.Filter.CreatedFrom,
		CreatedTo:     req.Filter.CreatedTo,
	}, sort, req.Paging.Limit, req.Paging.Offset)
	if err != nil {
		transportutil.Error(c, err)
		return
	}
	deviceProcessStatusHistoryResponses := make([]*dto.DeviceStatusHistory, 0, len(deviceProcessStatusHistoryData))
	for _, deviceProcessStatusHistory := range deviceProcessStatusHistoryData {
		deviceProcessStatusHistoryResponses = append(deviceProcessStatusHistoryResponses, &dto.DeviceStatusHistory{
			ID:                           deviceProcessStatusHistory.ID,
			ProductionOrderStageDeviceID: deviceProcessStatusHistory.ProductionOrderStageDeviceID,
			DeviceID:                     deviceProcessStatusHistory.DeviceID,
			StageID:                      deviceProcessStatusHistory.StageID.String,
			ProcessStatus:                deviceProcessStatusHistory.ProcessStatus,
			IsResolved:                   deviceProcessStatusHistory.IsResolved,
			UpdatedAt:                    deviceProcessStatusHistory.UpdatedAt.Time,
			UpdatedBy:                    deviceProcessStatusHistory.UpdatedBy.String,
			ErrorCode:                    deviceProcessStatusHistory.ErrorCode.String,
			ErrorReason:                  deviceProcessStatusHistory.ErrorReason.String,
			Description:                  deviceProcessStatusHistory.Description.String,
			CreatedAt:                    deviceProcessStatusHistory.CreatedAt,
			CreatedUserName:              deviceProcessStatusHistory.CreatedUserName.String,
			UpdatedUserName:              deviceProcessStatusHistory.UpdatedUserName.String,
		})
	}

	transportutil.SendJSONResponse(c, &dto.FindDeviceStatusHistoryResponse{
		DeviceStatusHistory: deviceProcessStatusHistoryResponses,
		Total:               total.Count,
	})
}
func (s productionOrderStageDeviceController) FindAvailabilityTime(c *gin.Context) {
	req := &dto.FindAvailabilityTimeRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	availabilityTime, err := s.productionOrderStageDeviceService.FindAvailabilityTime(c, &production_order_stage_device.FindLostTimeOpts{
		DeviceID: req.DeviceID,
		Date:     req.Date,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.FindAvailabilityTimeResponse{
		LossTime:    availabilityTime.LossTime,
		WorkingTime: availabilityTime.AvailabilityTime,
	})
}
func (s productionOrderStageDeviceController) FindEventLog(c *gin.Context) {
	req := &dto.FindEvenLogRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	eventLogs, err := s.productionOrderStageDeviceService.FindEventLog(c, &production_order_stage_device.FindEventLogOpts{
		DeviceID: req.DeviceID,
		Date:     req.Date,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	eventLogResponses := make([]*dto.FindEventLog, 0, len(eventLogs))
	for _, eventLog := range eventLogs {
		eventLogResponses = append(eventLogResponses, &dto.FindEventLog{
			ID:        eventLog.ID,
			DeviceID:  eventLog.DeviceID,
			StageID:   eventLog.StageID.String,
			Quantity:  eventLog.Quantity,
			Msg:       eventLog.Msg.String,
			Date:      eventLog.Date.String,
			CreatedAt: eventLog.CreatedAt,
		})
	}

	transportutil.SendJSONResponse(c, &dto.FindEventLogResponse{
		EventLogs: eventLogResponses,
	})
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
		AssignedQuantity:       req.AssignedQuantity,
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
	settings := &production_order_stage_device.Settings{}
	if req.Settings != nil {
		settings = &production_order_stage_device.Settings{
			DefectiveError: req.Settings.DefectiveError,
			Description:    req.Settings.Description,
		}
	}

	err = s.productionOrderStageDeviceService.Edit(c, &production_order_stage_device.EditProductionOrderStageDeviceOpts{
		ID:               req.ID,
		DeviceID:         req.DeviceID,
		Quantity:         req.Quantity,
		ProcessStatus:    req.ProcessStatus,
		Status:           req.Status,
		Responsible:      req.Responsible,
		AssignedQuantity: req.AssignedQuantity,
		Settings:         settings,
		Note:             req.Note,
		SanPhamLoi:       req.SanPhamLoi,
		UserID:           interceptor.UserIDFromCtx(c),
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
func (s productionOrderStageDeviceController) UpdateProcessDeviceHistoryIsSolved(c *gin.Context) {
	req := &dto.DeviceStatusHistoryUpdateSolved{}
	err := c.ShouldBind(req)
	if err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	err = s.productionOrderStageDeviceService.EditDeviceProcessHistoryIsSolved(c, &production_order_stage_device.EditDeviceProcessHistoryIsSolvedOpts{
		ID:     req.ID,
		UserID: interceptor.UserIDFromCtx(c),
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	transportutil.SendJSONResponse(c, &dto.DeviceStatusHistoryUpdateSolvedResponse{})
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
		"find-event-log",
		c.FindEventLog,
		&dto.FindEvenLogRequest{},
		&dto.FindEventLogResponse{},
		"Find event log",
	)
	routeutil.AddEndpoint(
		g,
		"find-device-status-history",
		c.FindProcessDeviceHistory,
		&dto.FindDeviceStatusHistoryRequest{},
		&dto.FindDeviceStatusHistoryResponse{},
		"Lịch sử thay đổi trạng thái của thiết bị",
	)
	routeutil.AddEndpoint(
		g,
		"update-device-status-history-solved",
		c.UpdateProcessDeviceHistoryIsSolved,
		&dto.DeviceStatusHistoryUpdateSolved{},
		&dto.DeviceStatusHistoryUpdateSolvedResponse{},
		"Cập nhật thay đổi trạng thái của thiết bị",
	)
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
	routeutil.AddEndpoint(
		g,
		"availability-time",
		c.FindAvailabilityTime,
		&dto.FindAvailabilityTimeRequest{},
		&dto.FindAvailabilityTimeResponse{},
		"FindAvailabilityTime",
	)
	routeutil.AddEndpoint(
		g,
		"find-working-device",
		c.FindWorkingDevice,
		&dto.FindWorkingDevice{},
		&dto.FindProductionOrderStageDevicesResponse{},
		"Lay danh sach thiet bi dang lam viec",
	)
}
