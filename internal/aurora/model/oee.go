package model

import (
	"time"
)

type OEE struct {
	JobRunningTime                int64
	AssignedWorkTime              int64
	AssignedWork                  []AssignWorkOEE
	Downtime                      int64
	TotalQuantity                 int64
	TotalAssignQuantity           int64
	TotalDefective                int64
	DeviceProgressStatusHistories []DeviceProgressStatusHistory
	DeviceID                      string
	DowntimeDetails               map[string]int64
	MachineOperator               []string
	ProductionOrderName           string
	ProductionOrderStageDevice    map[string]string
}

type SummaryOEE struct {
	TotalDowntime          int64
	TotalJobRunningTime    int64
	TotalAssignedWorkTime  int64
	TotalActualWorkingTime int64
}

type AssignWorkOEE struct {
	ID                     string
	ProductionOrderStageID string
	StageID                string
	EstimatedStartAt       time.Time
	EstimatedCompleteAt    time.Time
	Quantity               int64
	Defective              int64
	AssignedQuantity       int64
}
