package dto

import "time"

type StatisticsRequest struct {
	Month int16 `json:"month"`
	Year  int16 `json:"year"`
}

type DevicesErrorResponese struct {
	DeviceID      string `json:"device_id"`
	Quantity      int64  `json:"quantity"`
	QuantityError int64  `json:"quantity_error"`
	ErrorRate     string `json:"error_rate"`
}

type StopTimeResponse struct {
	DeviceID      string        `json:"device_id"`
	TotalStopTime time.Duration `json:"total_stop_time"`
}

type ErrorByStagesResponse struct {
	DeviceID   string `json:"device_id"`
	TotalError int64  `json:"quantity_error"`
}

type QuantityCompleteResponse struct {
	ListQuantityCompleteByDate []QuantityCompleteByDateResponse `json:"list_quantity_complete_by_date"`
	TotalQuantityComplete      int64                            `json:"total_quantity_complete"`
}

type QuantityCompleteByDateResponse struct {
	Day              int16 `json:"day"`
	QuantityComplete int64 `json:"quantity_complete"`
}

type QuantityDeliveryResponse struct {
	ListQuantityDeliveryByDate []QuantityDeliveryByDateResponse `json:"list_quantity_complete_by_date"`
	TotalQuantityDelivery      int64                            `json:"total_quantity_complete"`
}

type QuantityDeliveryByDateResponse struct {
	Day              int16 `json:"day"`
	QuantityDelivery int64 `json:"quantity_delivery"`
}

type StatisticsResponse struct {
	Top5DevicesError   []DevicesErrorResponese   `json:"top5_devices_error"`
	ListStopTime       []StopTimeResponse        `json:"list_stop_time"`
	ListErrorByStage   []ErrorByStagesResponse   `json:"list_error_by_stage"`
	QuantityComplete   *QuantityCompleteResponse `json:"quantity_complete"`
	QuantityDelivery   *QuantityDeliveryResponse `json:"quantity_delivery"`
	TotalDeviceWorking int64                     `json:"total_device_working"`
	SalesRevenue       int64                     `json:"sales_revenue"`
	ProductionRatio    string                    `json:"production_ratio"`
	OnTimeRatio        string                    `json:"on_time_ratio"`
}
