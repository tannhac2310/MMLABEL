package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
)

type DeviceConfigFilter struct {
	IDs               []string `json:"IDs"`
	Search            string   `json:"search"`
	ProductionOrderID string   `json:"productionOrderID"`
	ProductionPlanID  string   `json:"productionPlanID"`
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
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	ProductionOrderID   string                 `json:"productionOrderID"`
	ProductionPlanID    string                 `json:"productionPlanID"`
	ProductionOrderName string                 `json:"productionOrderName"`
	DeviceID            string                 `json:"deviceID"`
	DeviceCode          string                 `json:"deviceCode"`
	DeviceName          string                 `json:"deviceName"`
	DeviceConfig        map[string]interface{} `json:"deviceConfig"`
	Color               string                 `json:"color"`
	Description         string                 `json:"description"`
	CreatedBy           string                 `json:"createdBy"`
	CreatedAt           time.Time              `json:"createdAt"`
	UpdatedAt           time.Time              `json:"updatedAt"`
}

type CreateDeviceConfigRequest struct {
	ProductionOrderID string                 `json:"productionOrderID"`
	ProductionPlanID  string                 `json:"productionPlanID"`
	DeviceID          string                 `json:"deviceID"`
	DeviceConfig      map[string]interface{} `json:"deviceConfig" binding:"required"`
	Color             string                 `json:"color"`
	Description       string                 `json:"description"`
	Search            string                 `json:"search"`
}

type CreateDeviceConfigResponse struct {
	ID string `json:"id"`
}

type EditDeviceConfigRequest struct {
	ID                string                 `json:"id" binding:"required"`
	ProductionOrderID string                 `json:"productionOrderID"`
	ProductionPlanID  string                 `json:"productionPlanID"`
	DeviceID          string                 `json:"deviceID"`
	DeviceConfig      map[string]interface{} `json:"deviceConfig" binding:"required"`
	Color             string                 `json:"color"`
	Description       string                 `json:"description"`
	Search            string                 `json:"search"`
}

type EditDeviceConfigResponse struct {
}

type DeleteDeviceConfigRequest struct {
	ID string `json:"id"`
}

type DeleteDeviceConfigResponse struct {
}
