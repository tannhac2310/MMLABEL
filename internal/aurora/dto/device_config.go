package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"time"
)

type DeviceConfigFilter struct {
	Name string `json:"name"`
}

type FindDeviceConfigsRequest struct {
	Filter *DeviceConfigFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging   `json:"paging" binding:"required"`
}

type FindDeviceConfigsResponse struct {
	DeviceConfigs []*DeviceConfig `json:"deviceConfigs"`
	Total         int64           `json:"total"`
}
type DeviceConfig struct {
	ID                string                 `json:"id"`
	ProductionOrderID string                 `json:"productionOrderID"`
	DeviceID          string                 `json:"deviceID"`
	DeviceConfig      map[string]interface{} `json:"deviceConfig"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

type CreateDeviceConfigRequest struct {
	ProductionOrderID string                 `json:"productionOrderID" binding:"required"`
	DeviceID          string                 `json:"deviceID"`
	DeviceConfig      map[string]interface{} `json:"deviceConfig" binding:"required"`
}

type CreateDeviceConfigResponse struct {
	ID string `json:"id"`
}

type EditDeviceConfigRequest struct {
	ID                string                 `json:"id" binding:"required"`
	ProductionOrderID string                 `json:"productionOrderID" binding:"required"`
	DeviceID          string                 `json:"deviceID"`
	DeviceConfig      map[string]interface{} `json:"deviceConfig" binding:"required"`
}

type EditDeviceConfigResponse struct {
}

type DeleteDeviceConfigRequest struct {
	ID string `json:"id"`
}

type DeleteDeviceConfigResponse struct {
}
