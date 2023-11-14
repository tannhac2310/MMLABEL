package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

type DeviceFilter struct {
	Name string `json:"name"`
}

type FindDevicesRequest struct {
	Filter *DeviceFilter     `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindDevicesResponse struct {
	Devices []*Device `json:"devices"`
	Total   int64     `json:"total"`
}
type Device struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Code      string                 `json:"code"`
	OptionID  string                 `json:"optionID"`
	Data      map[string]interface{} `json:"data"`
	Status    enum.CommonStatus      `json:"status"`
	CreatedBy string                 `json:"createdBy"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type CreateDeviceRequest struct {
	Name     string                 `json:"name"`
	Code     string                 `json:"code"`
	OptionID string                 `json:"optionID"`
	Data     map[string]interface{} `json:"data"`
	Status   enum.CommonStatus      `json:"status"`
}

type CreateDeviceResponse struct {
	ID string `json:"id"`
}

type EditDeviceRequest struct {
	ID       string                 `json:"id" binding:"required"`
	Name     string                 `json:"name"`
	Code     string                 `json:"code"`
	OptionID string                 `json:"optionID"`
	Data     map[string]interface{} `json:"data"`
	Status   enum.CommonStatus      `json:"status"`
}

type EditDeviceResponse struct {
}

type DeleteDeviceRequest struct {
	ID string `json:"id"`
}

type DeleteDeviceResponse struct {
}
