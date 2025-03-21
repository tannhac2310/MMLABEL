package model

type OEE struct {
	ActualWorkingTime             int64
	JobRunningTime                int64
	AssignedWorkTime              int64
	AssignedWork                  []ProductionOrderStageDevice
	DowntimeStatistics            map[string]string
	Downtime                      int64
	TotalQuantity                 int64
	TotalAssignQuantity           int64
	TotalDefective                int64
	DeviceProgressStatusHistories []DeviceProgressStatusHistory
	DeviceID                      string
	DowntimeDetails               map[string]int64
	MachineOperator               []string
}
