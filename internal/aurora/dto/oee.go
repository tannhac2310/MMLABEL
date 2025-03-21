package dto

import (
	"time"
)

type FindOEERequest struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type FindOEEByDeviceResponse struct {
	OEEList []OEEByDevice `json:"oeeList"`
	Total   int64         `json:"total"`
}

type FindOEEByAssignedWorkResponse struct {
	OEEList []OEEByAssignedWork `json:"oeeList"`
	Total   int64               `json:"total"`
}

type OEEByDevice struct {
	DeviceID           string            `json:"deviceID"`
	ActualWorkingTime  int64             `json:"actualWorkingTime"`
	JobRunningTime     int64             `json:"jobRunningTime"`
	DownTime           int64             `json:"downTime"`
	AssignedWorkTime   int64             `json:"assignedWorkTime"`
	AssignedWork       []AssignedWork    `json:"assignedWork"`
	DownTimeStatistics map[string]string `json:"downTimeStatistics"`
	Availability       float64           `json:"availability"`
	Performance        float64           `json:"performance"`
	Quality            float64           `json:"quality"`
	TotalQuantity      int64             `json:"totalQuantity"`
	TotalDefective     int64             `json:"totalDefective"`
	OEE                float64           `json:"oee"`
}

type AssignedWork struct {
	ID                     string    `json:"id"`
	ProductionOrderStageID string    `json:"productionOrderID"`
	StageID                string    `json:"stageID"`
	EstimatedStartAt       time.Time `json:"estimatedStartAt"`
	EstimatedCompleteAt    time.Time `json:"estimatedCompleteAt"`
	Quantity               int64     `json:"quantity"`
	Defective              int64     `json:"defective"`
}

type OEEByAssignedWork struct {
	AssignedWorkID      string  `json:"assignedWorkID"`
	DeviceID            string  `json:"deviceID"`
	ActualWorkingTime   int64   `json:"actualWorkingTime"`
	JobRunningTime      int64   `json:"jobRunningTime"`
	DownTime            int64   `json:"downTime"`
	AssignedWorkTime    int64   `json:"assignedWorkTime"`
	Availability        float64 `json:"availability"`
	Performance         float64 `json:"performance"`
	Quality             float64 `json:"quality"`
	TotalQuantity       int64   `json:"totalQuantity"`
	TotalAssignQuantity int64   `json:"totalAssignQuantity"`
	TotalDefective      int64   `json:"totalDefective"`
	OEE                 float64 `json:"oee"`
	//DeviceProgressStatusHistories []DeviceStatusHistory `json:"deviceStatusHistory"`
	DowntimeDetails map[string]int64 `json:"downtimeDetails"`
	MachineOperator []string         `json:"machineOperator"`
}
