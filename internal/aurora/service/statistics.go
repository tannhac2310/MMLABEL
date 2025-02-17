package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/dto"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type StatisticsService interface {
	GetStatistics(ctx context.Context, request *dto.StatisticsRequest) (*dto.StatisticsResponse, error)
}

type statisticsService struct {
	statisticsRepo repository.StatisticsRepo
}

func (s statisticsService) GetStatistics(ctx context.Context, request *dto.StatisticsRequest) (*dto.StatisticsResponse, error) {
	year := request.Year
	month := request.Month

	if month == 0 {
		month = int16(time.Now().Month())
	}
	if year == 0 {
		year = int16(time.Now().Year())
	}

	devicesErrorResponese, err := s.getTop5DevicesError(ctx, month, year)
	if err != nil {
		return nil, err
	}

	listStopTimeResponese, err := s.getListStopTime(ctx, month, year)
	if err != nil {
		return nil, err
	}

	stagesErrorResponese, err := s.getStagesError(ctx, month, year)
	if err != nil {
		return nil, err
	}

	quantityCompleteResponse, err := s.getQuantityComplete(ctx, month, year)
	if err != nil {
		return nil, err
	}

	totalDeviceWorking, err := s.statisticsRepo.CountDeviceWorking(ctx)
	if err != nil {
		return nil, err
	}

	quantityDeliveryResponse, err := s.getQuantityDelivery(ctx, month, year)
	if err != nil {
		return nil, err
	}

	productionRation, err := s.statisticsRepo.FindProductionRatio(ctx, month, year)
	if err != nil {
		return nil, err
	}

	salesRevenue, err := s.statisticsRepo.SumSalesRevenue(ctx, month, year)
	if err != nil {
		return nil, err
	}

	onTimeRatio, err := s.statisticsRepo.FindOntimeRatio(ctx, month, year)
	if err != nil {
		return nil, err
	}

	response := &dto.StatisticsResponse{
		Top5DevicesError:   devicesErrorResponese,
		ListStopTime:       listStopTimeResponese,
		ListErrorByStage:   stagesErrorResponese,
		QuantityComplete:   quantityCompleteResponse,
		TotalDeviceWorking: totalDeviceWorking,
		QuantityDelivery:   quantityDeliveryResponse,
		SalesRevenue:       salesRevenue,
		ProductionRatio:    fmt.Sprintf("%.2f%%", productionRation),
		OnTimeRatio:        fmt.Sprintf("%.2f%%", onTimeRatio),
	}

	return response, nil
}

func (s statisticsService) getTop5DevicesError(ctx context.Context, month, year int16) ([]dto.DevicesErrorResponese, error) {
	// Lấy danh sách thiết bị với lỗi
	devicesError, err := s.statisticsRepo.FindDevicesError(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("FindDevicesError: %w", err)
	}

	// Tính toán QuantityError và tỷ lệ lỗi
	for i, device := range devicesError {
		device.ErrorRate = float64(device.QuantityError) / float64(device.Quantity) * 100
		devicesError[i] = device
	}

	// Sắp xếp theo tỷ lệ lỗi (QuantityError / Quantity) theo thứ tự tăng dần
	sort.Slice(devicesError, func(i, j int) bool {
		return devicesError[i].ErrorRate < devicesError[j].ErrorRate
	})

	// Lấy top 5 thiết bị có tỷ lệ lỗi thấp nhất
	if len(devicesError) > 20 {
		devicesError = devicesError[:20]
	}

	devicesErrorResponese := make([]dto.DevicesErrorResponese, 0, len(devicesError))

	for _, device := range devicesError {
		devicesErrorResponese = append(devicesErrorResponese, dto.DevicesErrorResponese{
			DeviceID:      device.DeviceID,
			Quantity:      device.Quantity,
			QuantityError: device.QuantityError,
			ErrorRate:     fmt.Sprintf("%.2f%%", device.ErrorRate),
		})
	}
	return devicesErrorResponese, nil
}

func (s statisticsService) getListStopTime(ctx context.Context, month, year int16) ([]dto.StopTimeResponse, error) {
	// Lấy danh sách thiết bị với lịch sử tiến trình
	listDeviceHistory, err := s.statisticsRepo.FindDevicesProgressHistory(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("FindDevicesProgressHistory: %w", err)
	}

	// Tạo danh sách kết quả
	devicesStopTimeResponse := make([]dto.StopTimeResponse, 0, len(listDeviceHistory))

	// Duyệt qua từng deviceID và danh sách lịch sử tiến trình của thiết bị đó
	for deviceID, deviceHistories := range listDeviceHistory {
		var totalStopTime time.Duration

		// Duyệt qua các phần tử theo cặp (i và i+1) để tính thời gian ngừng
		for i := 0; i < len(deviceHistories)-1; i += 2 {
			// Trừ thời gian của 2 sự kiện liên tiếp để tính khoảng thời gian ngừng
			stopDuration := deviceHistories[i+1].CreatedAt.Sub(deviceHistories[i].CreatedAt)
			// Cộng dồn tổng thời gian ngừng
			totalStopTime += stopDuration
		}

		// Thêm kết quả vào danh sách
		devicesStopTimeResponse = append(devicesStopTimeResponse, dto.StopTimeResponse{
			DeviceID:      deviceID,
			TotalStopTime: totalStopTime,
		})
	}

	// Trả về kết quả
	return devicesStopTimeResponse, nil
}

func (s statisticsService) getStagesError(ctx context.Context, month, year int16) ([]dto.ErrorByStagesResponse, error) {
	// Lấy danh sách thiết bị với lỗi
	stagesError, err := s.statisticsRepo.FindDevicesErrorByStage(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("FindDevicesErrorByStage: %w", err)
	}

	stagesErrorResponese := make([]dto.ErrorByStagesResponse, 0, len(stagesError))

	for _, stage := range stagesError {
		stagesErrorResponese = append(stagesErrorResponese, dto.ErrorByStagesResponse{
			DeviceID:   stage.StageID,
			TotalError: stage.QuantityError,
		})
	}
	return stagesErrorResponese, nil
}

func (s statisticsService) getQuantityComplete(ctx context.Context, month, year int16) (*dto.QuantityCompleteResponse, error) {
	// Lấy danh sách số lượng sản xuất theo ngày
	dateQuantity, err := s.statisticsRepo.ManufacturedQuantity(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("ManufacturedQuantity: %w", err)
	}

	dayInMonth := daysInMonth(year, month)
	quantityCompleteResponse := &dto.QuantityCompleteResponse{}
	listQuantityCompleteByDate := make([]dto.QuantityCompleteByDateResponse, 0, dayInMonth)
	var totalQuantityComplete int64 = 0
	// Lặp qua từng ngày trong tháng
	for i := int16(1); i <= int16(dayInMonth); i++ {
		var quantity int64 = 0
		// Kiểm tra ngày có tồn tại trong map
		if q, ok := dateQuantity[i]; ok {
			quantity = q
			totalQuantityComplete += quantity
		}
		listQuantityCompleteByDate = append(listQuantityCompleteByDate, dto.QuantityCompleteByDateResponse{
			Day:              i,
			QuantityComplete: quantity,
		})
	}

	quantityCompleteResponse.ListQuantityCompleteByDate = listQuantityCompleteByDate
	quantityCompleteResponse.TotalQuantityComplete = totalQuantityComplete
	return quantityCompleteResponse, nil
}

func (s statisticsService) getQuantityDelivery(ctx context.Context, month, year int16) (*dto.QuantityDeliveryResponse, error) {
	// Lấy danh sách số lượng sản xuất theo ngày
	dateQuantity, err := s.statisticsRepo.QuantityDelivery(ctx, month, year)
	if err != nil {
		return nil, fmt.Errorf("ManufacturedQuantity: %w", err)
	}

	dayInMonth := daysInMonth(year, month)
	quantityCompleteResponse := &dto.QuantityDeliveryResponse{}
	listQuantityCompleteByDate := make([]dto.QuantityDeliveryByDateResponse, 0, dayInMonth)
	var totalQuantityComplete int64 = 0
	// Lặp qua từng ngày trong tháng
	for i := int16(1); i <= int16(dayInMonth); i++ {
		var quantity int64 = 0
		// Kiểm tra ngày có tồn tại trong map
		if q, ok := dateQuantity[i]; ok {
			quantity = q
			totalQuantityComplete += quantity
		}
		listQuantityCompleteByDate = append(listQuantityCompleteByDate, dto.QuantityDeliveryByDateResponse{
			Day:              i,
			QuantityDelivery: quantity,
		})
	}

	quantityCompleteResponse.ListQuantityDeliveryByDate = listQuantityCompleteByDate
	quantityCompleteResponse.TotalQuantityDelivery = totalQuantityComplete
	return quantityCompleteResponse, nil
}

func NewStatisticsService(statisticsRepo repository.StatisticsRepo) StatisticsService {
	return &statisticsService{
		statisticsRepo: statisticsRepo,
	}
}

func daysInMonth(year int16, month int16) int {
	// Xác định tháng tiếp theo
	if month == 12 {
		month = 1
		year++
	} else {
		month++
	}

	// Tạo ngày đầu tiên của tháng tiếp theo
	nextMonth := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Lùi lại 1 ngày để lấy ngày cuối cùng của tháng
	lastDayOfMonth := nextMonth.AddDate(0, 0, -1)

	return lastDayOfMonth.Day()
}
