package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type ProductionOrderStageDeviceFilter struct {
	ProductionOrderStageID string `json:"productionOrderStageID"`
	ProductionOrderID      string `json:"productionOrderID"`
}

type FindProductionOrderStageDevicesRequest struct {
	Filter *ProductionOrderStageDeviceFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging                 `json:"paging" binding:"required"`
}

type FindProductionOrderStageDevicesResponse struct {
	ProductionOrderStageDevices []*ProductionOrderStageDevice `json:"productionOrderStageDevices"`
	Total                       int64                         `json:"total"`
}
type ProductionOrderStageDevice struct {
	ID                     string                                `json:"id"`
	ProductionOrderStageID string                                `json:"productionOrderStageID"`
	DeviceID               string                                `json:"deviceID"`
	Quantity               int64                                 `json:"quantity"`
	ProcessStatus          enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	Status                 enum.CommonStatus                     `json:"status"`
	Responsible            []string                              `json:"responsible"`
	Settings               map[string]interface{}                `json:"settings"`
	Note                   string                                `json:"note"`
}

type CreateProductionOrderStageDeviceRequest struct {
	ProductionOrderStageID string                                `json:"productionOrderStageID"`
	DeviceID               string                                `json:"deviceID"`
	Quantity               int64                                 `json:"quantity"`
	ProcessStatus          enum.ProductionOrderStageDeviceStatus `json:"processStatus"`
	Status                 enum.CommonStatus                     `json:"status"`
	Responsible            []string                              `json:"responsible"`
	Settings               map[string]interface{}                `json:"settings"`
	Note                   string                                `json:"note"`
}

type CreateProductionOrderStageDeviceResponse struct {
	ID string `json:"id"`
}

type EditProductionOrderStageDeviceRequest struct {
	ID            string                                `json:"id" binding:"required"`
	DeviceID      string                                `json:"deviceID" binding:"required"`
	Quantity      int64                                 `json:"quantity"`
	ProcessStatus enum.ProductionOrderStageDeviceStatus `json:"processStatus" binding:"required"`
	Status        enum.CommonStatus                     `json:"status" binding:"required"`
	Responsible   []string                              `json:"responsible"`
	Settings      map[string]interface{}                `json:"settings"`
	Note          string                                `json:"note"`
}

type EditProductionOrderStageDeviceResponse struct {
}

type DeleteProductionOrderStageDeviceRequest struct {
	IDs []string `json:"id"`
}

type DeleteProductionOrderStageDeviceResponse struct {
}
