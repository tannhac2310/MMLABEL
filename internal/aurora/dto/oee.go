package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"time"
)

type FindOEEFilter struct {
	DateFrom                     string `json:"dateFrom"`
	DateTo                       string `json:"dateTo"`
	ProductionOrderStageDeviceID string `json:"productionOrderStageDeviceID"`
	ProductionOrderID            string `json:"productionOrderStageID"`
	DeviceID                     string `json:"deviceID"`
}
type FindOEERequest struct {
	Filter *FindOEEFilter    `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindOEEByDeviceResponse struct {
	OEEList []OEEByDeviceResponse `json:"oeeList"`
	Total   int64                 `json:"total"`
}

type FindOEEByAssignedWorkResponse struct {
	OEEList []OEEByAssignedWorkResponse `json:"oeeList"`
	Total   int64                       `json:"total"`
}

type OEEByDeviceResponse struct {
	DeviceID                      string                 `json:"deviceID"`
	ActualWorkingTime             int64                  `json:"actualWorkingTime"`
	JobRunningTime                int64                  `json:"jobRunningTime"`
	DownTime                      int64                  `json:"downTime"`
	AssignedWorkTime              int64                  `json:"assignedWorkTime"`
	AssignedWork                  []AssignedWorkResponse `json:"assignedWork"`
	DowntimeDetails               map[string]int64       `json:"downtimeDetails"`
	Availability                  float64                `json:"availability"`
	Performance                   float64                `json:"performance"`
	Quality                       float64                `json:"quality"`
	TotalQuantity                 int64                  `json:"totalQuantity"`
	TotalDefective                int64                  `json:"totalDefective"`
	OEE                           float64                `json:"oee"`
	DeviceProgressStatusHistories []DeviceStatusHistory  `json:"deviceStatusHistory"`
}

type AssignedWorkResponse struct {
	ID                     string    `json:"id"`
	ProductionOrderStageID string    `json:"productionOrderID"`
	StageID                string    `json:"stageID"`
	EstimatedStartAt       time.Time `json:"estimatedStartAt"`
	EstimatedCompleteAt    time.Time `json:"estimatedCompleteAt"`
	Quantity               int64     `json:"quantity"`
	Defective              int64     `json:"defective"`
}

type OEEByAssignedWorkResponse struct {
	AssignedWorkID      string           `json:"assignedWorkID"`
	ProductionOrderName string           `json:"ProductionOrderID"`
	DeviceID            string           `json:"deviceID"`
	ActualWorkingTime   int64            `json:"actualWorkingTime"`
	JobRunningTime      int64            `json:"jobRunningTime"`
	DownTime            int64            `json:"downTime"`
	AssignedWorkTime    int64            `json:"assignedWorkTime"`
	Availability        float64          `json:"availability"`
	Performance         float64          `json:"performance"`
	Quality             float64          `json:"quality"`
	TotalQuantity       int64            `json:"totalQuantity"`
	TotalAssignQuantity int64            `json:"totalAssignQuantity"`
	TotalDefective      int64            `json:"totalDefective"`
	OEE                 float64          `json:"oee"`
	DowntimeDetails     map[string]int64 `json:"downtimeDetails"`
	MachineOperator     []string         `json:"machineOperator"`
}
