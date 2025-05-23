package controller

import (
	"github.com/gin-gonic/gin"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/oee"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/routeutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
)

type OEEController interface {
	CalcOEEByDevice(c *gin.Context)
	CalcOEEByAssignedWork(c *gin.Context)
}

type oeeController struct {
	oeeService oee.Service
}

func (o oeeController) CalcOEEByDevice(c *gin.Context) {
	req := &dto.FindOEERequest{}
	if err := c.ShouldBind(&req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}

	datas, mapProductionOrderStageDevice, err := o.oeeService.CalcOEEByDevice(c, &oee.CalcOEEOpts{
		ProductionOrderID:            req.Filter.ProductionOrderID,
		ProductionOrderStageDeviceID: req.Filter.ProductionOrderStageDeviceID,
		DateFrom:                     req.Filter.DateFrom,
		DateTo:                       req.Filter.DateTo,
		DeviceID:                     req.Filter.DeviceID,
		Limit:                        req.Paging.Limit,
		Offset:                       req.Paging.Offset,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	oeeList := make([]dto.OEEByDeviceResponse, 0, len(datas))

	for deviceID, data := range datas {
		availability := 1.0
		performance := 1.0
		quality := 1.0

		if data.JobRunningTime+data.Downtime > 0 {
			availability = float64(data.JobRunningTime) / float64(data.JobRunningTime+data.Downtime)
		}
		if data.JobRunningTime > 0 {
			performance = float64(data.AssignedWorkTime-data.Downtime) / float64(data.JobRunningTime)
		}
		if data.TotalQuantity > 0 {
			quality = float64(data.TotalQuantity) / float64(data.TotalQuantity+data.TotalDefective)
		}

		model := dto.OEEByDeviceResponse{
			DeviceID:            deviceID,
			JobRunningTime:      data.JobRunningTime,
			AssignedWorkTime:    data.AssignedWorkTime,
			DownTime:            data.Downtime,
			DowntimeDetails:     data.DowntimeDetails,
			Availability:        availability,
			Performance:         performance,
			Quality:             quality,
			TotalQuantity:       data.TotalQuantity,
			TotalDefective:      data.TotalDefective,
			TotalAssignQuantity: data.TotalAssignQuantity,
			OEE:                 availability * performance * quality,
		}

		assignedWork := make([]dto.AssignedWorkResponse, 0, len(data.AssignedWork))
		for _, work := range data.AssignedWork {
			assignedWork = append(assignedWork, dto.AssignedWorkResponse{
				ID:                     work.ID,
				ProductionOrderStageID: work.ProductionOrderStageID,
				StageID:                work.StageID,
				EstimatedStartAt:       work.EstimatedStartAt,
				EstimatedCompleteAt:    work.EstimatedCompleteAt,
				Quantity:               work.Quantity,
				Defective:              work.Defective,
				AssignQuantity:         work.AssignedQuantity,
			})
		}
		model.AssignedWork = assignedWork

		deviceProgressStatusHistories := make([]dto.DeviceStatusHistory, 0, len(data.DeviceProgressStatusHistories))
		for _, deviceProcessStatusHistory := range data.DeviceProgressStatusHistories {
			deviceProgressStatusHistories = append(deviceProgressStatusHistories, dto.DeviceStatusHistory{
				ID:                           deviceProcessStatusHistory.ID,
				ProductionOrderStageDeviceID: deviceProcessStatusHistory.ProductionOrderStageDeviceID,
				StageID:                      mapProductionOrderStageDevice[deviceProcessStatusHistory.ProductionOrderStageDeviceID],
				DeviceID:                     deviceProcessStatusHistory.DeviceID,
				ProcessStatus:                deviceProcessStatusHistory.ProcessStatus,
				IsResolved:                   deviceProcessStatusHistory.IsResolved,
				UpdatedAt:                    deviceProcessStatusHistory.UpdatedAt.Time,
				UpdatedBy:                    deviceProcessStatusHistory.UpdatedBy.String,
				ErrorCode:                    deviceProcessStatusHistory.ErrorCode.String,
				ErrorReason:                  deviceProcessStatusHistory.ErrorReason.String,
				Description:                  deviceProcessStatusHistory.Description.String,
				CreatedAt:                    deviceProcessStatusHistory.CreatedAt,
			})
		}
		model.DeviceProgressStatusHistories = deviceProgressStatusHistories

		oeeList = append(oeeList, model)
	}

	transportutil.SendJSONResponse(c, &dto.FindOEEByDeviceResponse{
		Total:   int64(len(oeeList)),
		OEEList: oeeList,
	})
}

func (o oeeController) CalcOEEByAssignedWork(c *gin.Context) {
	req := &dto.FindOEERequest{}
	if err := c.ShouldBind(&req); err != nil {
		transportutil.Error(c, apperror.ErrInvalidArgument.WithDebugMessage(err.Error()))
		return
	}
	datas, total, err := o.oeeService.CalcOEEByAssignedWork(c, &oee.CalcOEEOpts{
		ProductionOrderID:            req.Filter.ProductionOrderID,
		ProductionOrderStageDeviceID: req.Filter.ProductionOrderStageDeviceID,
		DateFrom:                     req.Filter.DateFrom,
		DateTo:                       req.Filter.DateTo,
		DeviceID:                     req.Filter.DeviceID,
		Limit:                        req.Paging.Limit,
		Offset:                       req.Paging.Offset,
	})
	if err != nil {
		transportutil.Error(c, err)
		return
	}

	oeeList := make([]dto.OEEByAssignedWorkResponse, 0, len(datas))

	for assignedWorkID, data := range datas {
		availability := 1.0
		performance := 1.0
		quality := 1.0

		if data.JobRunningTime+data.Downtime > 0 {
			availability = float64(data.JobRunningTime) / float64(data.JobRunningTime+data.Downtime)
		}
		if data.JobRunningTime > 0 {
			performance = float64(data.AssignedWorkTime-data.Downtime) / float64(data.JobRunningTime)
		}
		if data.TotalQuantity > 0 {
			quality = float64(data.TotalQuantity) / float64(data.TotalQuantity+data.TotalDefective)
		}

		model := dto.OEEByAssignedWorkResponse{
			AssignedWorkID:      assignedWorkID,
			ProductionOrderName: data.ProductionOrderName,
			JobRunningTime:      data.JobRunningTime,
			AssignedWorkTime:    data.AssignedWorkTime,
			DowntimeDetails:     data.DowntimeDetails,
			DownTime:            data.Downtime,
			Availability:        availability,
			Performance:         performance,
			Quality:             quality,
			TotalQuantity:       data.TotalQuantity,
			TotalAssignQuantity: data.TotalAssignQuantity,
			TotalDefective:      data.TotalDefective,
			OEE:                 availability * performance * quality,
			DeviceID:            data.DeviceID,
			MachineOperator:     data.MachineOperator,
		}
		oeeList = append(oeeList, model)
	}

	transportutil.SendJSONResponse(c, &dto.FindOEEByAssignedWorkResponse{
		Total:   total,
		OEEList: oeeList,
	})
}

func RegisterOEEController(
	r *gin.RouterGroup,
	oeeService oee.Service,
) {
	g := r.Group("oee")

	var c OEEController = &oeeController{
		oeeService: oeeService,
	}

	routeutil.AddEndpoint(
		g,
		"device",
		c.CalcOEEByDevice,
		&dto.FindOEERequest{},
		&dto.FindOEEByDeviceResponse{},
		"Calc OEE by device",
	)

	routeutil.AddEndpoint(
		g,
		"assigned-work",
		c.CalcOEEByAssignedWork,
		&dto.FindOEERequest{},
		&dto.FindOEEByAssignedWorkResponse{},
		"Calc OEE By assigned work",
	)
}
