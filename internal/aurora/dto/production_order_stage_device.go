package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ProductionOrderStageDeviceFilter struct {
	ID                           string                                  `json:"id"`
	IDs                          []string                                `json:"ids"`
	ProductionOrderIDs           []string                                `json:"productionOrderIDs"`
	DeviceIDs                    []string                                `json:"deviceIDs"`
	ProductionStageIDs           []string                                `json:"productionStageIDs"`
	ProductionOrderStageStatuses []enum.ProductionOrderStageStatus       `json:"productionOrderStageStatuses"`
	Responsible                  []string                                `json:"responsible"`
	ProcessStatuses              []enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	//EstimatedStartAtFrom  time.Time                             `json:"estimatedAtFrom"`
	//EstimatedCompleteAtTo time.Time                             `json:"estimatedAtTo"`
	//StartAt               time.Time                             `json:"startAt"`
	//CompleteAt            time.Time                             `json:"completeAt"`
}

type FindProductionOrderStageDevicesRequest struct {
	Filter *ProductionOrderStageDeviceFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging                 `json:"paging" binding:"required"`
}

type FindProductionOrderStageDevicesResponse struct {
	ProductionOrderStageDevices []*ProductionOrderStageDevice `json:"productionOrderStageDevices"`
	Total                       int64                         `json:"total"`
}
type FindWorkingDevice struct {
}
type POStageDeviceResponsible struct {
	ID              string `json:"id"`
	POStageDeviceID string `json:"poStageDeviceID"`
	UserID          string `json:"userID"`
	ResponsibleName string `json:"responsibleName"`
}
type ProductionOrderStageDevice struct {
	ID                                      string                                `json:"id"`
	ProductionOrderID                       string                                `json:"productionOrderID"`
	ProductionOrderName                     string                                `json:"productionOrderName"`
	ProductionOrderData                     *ProductionOrderData                  `json:"productionOrderData"`
	ProductionOrderStatus                   enum.ProductionOrderStatus            `json:"productionOrderStatus"`
	ProductionOrderStageName                string                                `json:"productionOrderStageName"`
	ProductionOrderStageCode                string                                `json:"productionOrderStageCode"`
	ProductionOrderStageStatus              enum.ProductionOrderStageStatus       `json:"productionOrderStageStatus"`
	ProductionOrderStageID                  string                                `json:"productionOrderStageID"`
	ProductionOrderStageStartedAt           time.Time                             `json:"productionOrderStageStartedAt"`
	ProductionOrderStageCompletedAt         time.Time                             `json:"productionOrderStageCompletedAt"`
	ProductionOrderStageEstimatedStartAt    time.Time                             `json:"productionOrderStageEstimatedStartAt"`
	ProductionOrderStageEstimatedCompleteAt time.Time                             `json:"productionOrderStageEstimatedCompleteAt"`
	EstimatedStartAt                        time.Time                             `json:"estimatedStartAt"`
	EstimatedCompleteAt                     time.Time                             `json:"estimatedCompleteAt"`
	StartedAt                               time.Time                             `json:"startedAt"`
	CompleteAt                              time.Time                             `json:"completeAt"`
	DeviceID                                string                                `json:"deviceID"`
	DeviceName                              string                                `json:"deviceName"`
	Quantity                                int64                                 `json:"quantity"`
	AssignedQuantity                        int64                                 `json:"assignedQuantity"`
	ProcessStatus                           enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	Status                                  enum.CommonStatus                     `json:"status"`
	Responsible                             []*POStageDeviceResponsible           `json:"responsible"`
	Settings                                map[string]interface{}                `json:"settings"`
	Note                                    string                                `json:"note"`
}

type CreateProductionOrderStageDeviceRequest struct {
	ProductionOrderStageID string                                `json:"productionOrderStageID" binding:"required"`
	DeviceID               string                                `json:"deviceID" binding:"required"`
	Quantity               int64                                 `json:"quantity"`
	ProcessStatus          enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	Status                 enum.CommonStatus                     `json:"status" binding:"required"`
	Responsible            []string                              `json:"responsible"`
	Settings               map[string]interface{}                `json:"settings"`
	Note                   string                                `json:"note"`
	AssignedQuantity       int64                                 `json:"assignedQuantity"`
	EstimatedStartAt       time.Time                             `json:"estimatedStartAt"`
	EstimatedCompleteAt    time.Time                             `json:"estimatedCompleteAt"`
}

type CreateProductionOrderStageDeviceResponse struct {
	ID string `json:"id"`
}

type EditProductionOrderStageDeviceRequest struct {
	ID                  string                                  `json:"id" binding:"required"`
	DeviceID            string                                  `json:"deviceID"`
	Quantity            int64                                   `json:"quantity"`
	AssignedQuantity    int64                                   `json:"assignedQuantity"`
	ProcessStatus       enum.ProductionOrderStageDeviceStatus   `json:"processStatus"`
	Status              enum.CommonStatus                       `json:"status"`
	Responsible         []string                                `json:"responsible"`
	Settings            *EditProductionOrderStageDeviceSettings `json:"settings"`
	Note                string                                  `json:"note"`
	SanPhamLoi          int64                                   `json:"sanPhamLoi"`
	EstimatedStartAt    time.Time                               `json:"estimatedStartAt"`
	EstimatedCompleteAt time.Time                               `json:"estimatedCompleteAt"`
}
type EditProductionOrderStageDeviceSettings struct {
	DefectiveError string `json:"defectiveError"`
	Description    string `json:"description"`
}
type EditProductionOrderStageDeviceResponse struct {
}

type DeleteProductionOrderStageDeviceRequest struct {
	IDs []string `json:"id"`
}

type DeleteProductionOrderStageDeviceResponse struct {
}

type FindEvenLogRequest struct {
	DeviceID string `json:"name" `
	Date     string `json:"date"`
}

type FindEventLog struct {
	ID         int64     `json:"id"`
	DeviceID   string    `json:"deviceID"`
	DeviceName string    `json:"deviceName"`
	StageID    string    `json:"stageID"`
	Quantity   float64   `json:"quantity"`
	Msg        string    `json:"msg"`
	Date       string    `json:"date"`
	CreatedAt  time.Time `json:"createdAt"`
}
type FindEventLogResponse struct {
	EventLogs []*FindEventLog `json:"eventLogs"`
}

type DeviceStatusHistoryFilter struct {
	ProductionOrderStageID string    `json:"productionOrderStageID"`
	ProductionDeviceID     string    `json:"productionDeviceID"`
	CreatedFrom            time.Time `json:"createdFrom"`
	CreatedTo              time.Time `json:"createdTo"`
}

type FindDeviceStatusHistoryFilter struct {
	ProcessStatus []int8    `json:"processStatus"`
	DeviceID      string    `json:"deviceID"`
	IsResolved    int16     `json:"isResolved"`
	ErrorCodes    []string  `json:"errorCodes"`
	CreatedFrom   time.Time `json:"createdFrom"`
	CreatedTo     time.Time `json:"createdTo"`
}

type FindDeviceStatusHistoryRequest struct {
	Filter *FindDeviceStatusHistoryFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging              `json:"paging" binding:"required"`
	Sort   *commondto.Sort                `json:"sort"`
}

type FindDeviceStatusHistoryResponse struct {
	DeviceStatusHistory []*DeviceStatusHistory `json:"deviceStatusHistory"`
	Total               int64                  `json:"total"`
}
type DeviceStatusHistory struct {
	ID                           string                                `json:"id"`
	ProductionOrderStageDeviceID string                                `json:"productionOrderStageDeviceID"`
	DeviceID                     string                                `json:"deviceID"`
	StageID                      string                                `json:"stageID"`
	ProcessStatus                enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	IsResolved                   int16                                 `json:"isResolved"`
	UpdatedAt                    time.Time                             `json:"updatedAt"`
	UpdatedBy                    string                                `json:"updatedBy"`
	ErrorCode                    string                                `json:"errorCode"`
	ErrorReason                  string                                `json:"errorReason"`
	Description                  string                                `json:"description"`
	CreatedAt                    time.Time                             `json:"createdAt"`
	CreatedUserName              string                                `json:"createdUserName"`
	UpdatedUserName              string                                `json:"updatedUserName"`
}

type DeviceStatusHistoryUpdateSolved struct {
	ID string `json:"id"`
}
type DeviceStatusHistoryUpdateSolvedResponse struct {
}

// lostime
type FindAvailabilityTimeRequest struct {
	DeviceID string `json:"deviceID"`
	Date     string `json:"date"`
}

type FindAvailabilityTimeResponse struct {
	LossTime    int64 `json:"lossTime"`
	WorkingTime int64 `json:"workingTime"`
}
