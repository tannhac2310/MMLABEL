package dto

type FindOEERequest struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type FindOEEResponse struct {
	OEEList []OEE `json:"oeeList"`
	Total   int64 `json:"total"`
}

type OEE struct {
	DeviceID           string            `json:"deviceID"`
	ActualWorkingTime  int64             `json:"actualWorkingTime"`
	JobRunningTime     int64             `json:"jobRunningTime"`
	AssignedWorkTime   int64             `json:"assignedWorkTime"`
	DowntimeStatistics map[string]string `json:"downtimeStatistics"`
	Availability       float64           `json:"availability"`
	Performance        float64           `json:"performance"`
	Quality            float64           `json:"quality"`
	TotalQuantity      int64             `json:"totalQuantity"`
	TotalDefective     int64             `json:"totalDefective"`
}
