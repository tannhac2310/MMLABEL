package repository

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type StatisticsRepo interface {
	FindDevicesError(ctx context.Context, month, year int16) ([]model.DevicesError, error)
	FindDevicesProgressHistory(ctx context.Context, month, year int16) (map[string][]model.DevicesProgressHistory, error)
	FindDevicesErrorByStage(ctx context.Context, month, year int16) ([]model.DevicesErrorByStage, error)
	ManufacturedQuantity(ctx context.Context, month, year int16) (map[int16]int64, error)
	QuantityDelivery(ctx context.Context, month, year int16) (map[int16]int64, error)
	CountDeviceWorking(ctx context.Context) (int64, error)
	FindProductionRatio(ctx context.Context, month, year int16) (float64, error)
	FindOntimeRatio(ctx context.Context, month, year int16) (float64, error)
}

type statisticsRepo struct {
}

func (s statisticsRepo) FindProductionRatio(ctx context.Context, month, year int16) (float64, error) {
	// Biến để lưu tỷ lệ sản xuất
	var ratio float64

	// Câu truy vấn SQL
	query := `
        SELECT 
            (SUM(quantity) - SUM(value::numeric)) / NULLIF(SUM(qty_delivered), 0) AS production_ratio
        FROM 
            production_order_stage_devices posd
        JOIN 
            production_order_stages pos ON posd.production_order_stage_id = pos.id
        JOIN 
            production_orders po ON pos.production_order_id = po.id
        JOIN 
            LATERAL jsonb_each_text(posd.settings) AS settings(key, value) ON TRUE
        WHERE 
            jsonb_typeof(posd.settings) = 'object' 
            AND EXTRACT(MONTH FROM delivery_date) = $1 
            AND EXTRACT(YEAR FROM delivery_date) = $2
            AND settings.key = 'san_pham_loi';
    `

	// Thực hiện truy vấn và lưu kết quả vào biến `ratio`
	err := cockroach.QueryRow(ctx, query, month, year).Scan(&ratio)
	if err != nil {
		return 0, fmt.Errorf("cockroach.QueryRow: %w", err)
	}

	// Trả về tỷ lệ sản xuất tính theo phần trăm
	return ratio * 100, nil
}
func (s statisticsRepo) FindOntimeRatio(ctx context.Context, month, year int16) (float64, error) {
	// Biến để lưu tỷ lệ sản xuất
	var ratio float64

	// Câu truy vấn SQL
	query := `
        SELECT 
			SUM(CASE 
					WHEN delivery_date <= estimated_complete_at THEN 1 
					ELSE 0 
				END) / COUNT(*) AS on_time_ratio
		FROM 
			public.production_orders
		WHERE 
			EXTRACT(MONTH FROM delivery_date) = $1 
			AND EXTRACT(YEAR FROM delivery_date) = $2
			AND estimated_complete_at IS NOT NULL;

    `

	// Thực hiện truy vấn và lưu kết quả vào biến `ratio`
	err := cockroach.QueryRow(ctx, query, month, year).Scan(&ratio)
	if err != nil {
		return 0, fmt.Errorf("cockroach.QueryRow: %w", err)
	}

	// Trả về tỷ lệ sản xuất tính theo phần trăm
	return ratio * 100, nil
}
func (s statisticsRepo) CountDeviceWorking(ctx context.Context) (int64, error) {
	// Biến để lưu số lượng kết quả
	var count int64

	// Câu truy vấn SQL
	query := `
		SELECT COUNT(*) 
		FROM public.production_order_stage_devices
		WHERE process_status = $1;
	`

	// Thực hiện truy vấn và lưu kết quả vào biến `count`
	err := cockroach.QueryRow(ctx, query, enum.ProductionOrderStageDeviceStatusStart).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("cockroach.QueryRow: %w", err)
	}

	// Trả về số lượng thiết bị
	return count, nil
}

func (s statisticsRepo) FindDevicesError(ctx context.Context, month, year int16) ([]model.DevicesError, error) {
	// Khai báo slice để lưu kết quả
	result := make([]model.DevicesError, 0)

	// Câu truy vấn SQL với tham số tháng và năm
	query := `
		SELECT 
			device_id, 
			sum(value::numeric) AS quantity_error,
			sum(quantity) AS quantity
		FROM 
			production_order_stage_devices pos
		JOIN 
			LATERAL jsonb_each_text(pos.settings) AS settings(key, value) ON TRUE
		WHERE 
			jsonb_typeof(pos.settings) = 'object' 
			AND EXTRACT(MONTH FROM complete_at) = $1 
            AND EXTRACT(YEAR FROM complete_at) = $2
            AND key = 'san_pham_loi'
		GROUP BY 
			device_id;
	`

	// Thực hiện truy vấn với tham số tháng và năm
	err := cockroach.Select(ctx, query, month, year).ScanAll(&result)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	// Trả về kết quả
	return result, nil
}

func (s statisticsRepo) ManufacturedQuantity(ctx context.Context, month, year int16) (map[int16]int64, error) {
	// Khai báo slice để lưu kết quả
	result := make([]model.QuantityByDate, 0)

	// Câu truy vấn SQL với tham số tháng và năm
	query := `
		SELECT 
			EXTRACT(DAY FROM complete_at) AS day,
			sum(quantity) - sum(value::numeric) AS quantity_complete
		FROM 
			production_order_stage_devices pos
		JOIN 
			LATERAL jsonb_each_text(pos.settings) AS settings(key, value) ON TRUE
		WHERE 
			jsonb_typeof(pos.settings) = 'object' 
			AND EXTRACT(MONTH FROM complete_at) = $1 
			AND EXTRACT(YEAR FROM complete_at) = $2
			AND settings.key = 'san_pham_loi'
		GROUP BY 
			day;
	`

	// Thực hiện truy vấn với tham số tháng và năm
	err := cockroach.Select(ctx, query, month, year).ScanAll(&result)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	resultMap := make(map[int16]int64)
	for _, item := range result {
		resultMap[item.Day] = item.QuantityComplete
	}

	// Trả về kết quả
	return resultMap, nil
}

func (s statisticsRepo) QuantityDelivery(ctx context.Context, month, year int16) (map[int16]int64, error) {
	// Khai báo slice để lưu kết quả
	result := make([]model.QuantityDeliveryByDate, 0)

	// Câu truy vấn SQL với tham số tháng và năm
	query := `
		SELECT 
			EXTRACT(DAY FROM delivery_date) AS day,
			sum(qty_delivered) AS quantity_delivery
		FROM 
			production_orders
		WHERE 
			EXTRACT(MONTH FROM delivery_date) = $1 
			AND EXTRACT(YEAR FROM delivery_date) = $2
		GROUP BY 
			day
		order by day;
	`

	// Thực hiện truy vấn với tham số tháng và năm
	err := cockroach.Select(ctx, query, month, year).ScanAll(&result)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	resultMap := make(map[int16]int64)
	for _, item := range result {
		resultMap[item.Day] = item.QuantityDelivery
	}

	// Trả về kết quả
	return resultMap, nil
}

func (s statisticsRepo) FindDevicesProgressHistory(ctx context.Context, month, year int16) (map[string][]model.DevicesProgressHistory, error) {
	// Khai báo slice để lưu kết quả
	result := make([]model.DevicesProgressHistory, 0)

	// Câu truy vấn SQL với tham số tháng và năm
	query := `
		SELECT device_id, process_status, created_at 
		FROM device_progress_status_history
		WHERE EXTRACT(MONTH FROM created_at) = $1 AND EXTRACT(YEAR FROM created_at) = $2
		ORDER BY created_at ASC
	`

	// Thực hiện truy vấn với tham số tháng và năm
	err := cockroach.Select(ctx, query, month, year).ScanAll(&result)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}
	// Chuyển đổi kết quả thành map với device_id là key
	deviceMap := make(map[string][]model.DevicesProgressHistory)

	for _, device := range result {
		deviceMap[device.DeviceID] = append(deviceMap[device.DeviceID], device)
	}

	// Duyệt qua từng deviceID và loại bỏ các thiết bị không hợp lệ
	for deviceID, devices := range deviceMap {
		// Nếu số lượng thiết bị ít hơn 2, xóa khỏi map
		if len(devices) < 2 {
			delete(deviceMap, deviceID)
			continue
		}

		// Kiểm tra trạng thái đầu tiên, nếu là "Start", bỏ qua thiết bị đầu tiên
		if devices[0].ProcessStatus == enum.ProductionOrderStageDeviceStatusStart {
			devices = devices[1:]
		}

		// Kiểm tra trạng thái cuối cùng, nếu là "Failed", bỏ qua thiết bị cuối cùng
		if devices[len(devices)-1].ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
			devices = devices[:len(devices)-1]
		}

		// Nếu sau khi cắt bớt mà số lượng thiết bị còn lại ít hơn 2, xóa khỏi map
		if len(devices) < 2 {
			delete(deviceMap, deviceID)
		} else {
			// Cập nhật lại map với danh sách thiết bị đã được cắt bớt
			deviceMap[deviceID] = devices
		}
	}

	// Trả về kết quả
	return deviceMap, nil
}

func (s statisticsRepo) FindDevicesErrorByStage(ctx context.Context, month, year int16) ([]model.DevicesErrorByStage, error) {
	// Khai báo slice để lưu kết quả
	result := make([]model.DevicesErrorByStage, 0)

	// Câu truy vấn SQL với tham số tháng và năm
	query := `
		SELECT 
			stages.stage_id,
			sum(value::numeric) AS quantity_error
		FROM 
			production_order_stage_devices pos
		JOIN 
			production_order_stages stages ON pos.production_order_stage_id = stages.id
		JOIN 
			LATERAL jsonb_each_text(pos.settings) AS settings(key, value) ON TRUE
		WHERE 
			jsonb_typeof(pos.settings) = 'object' 
			AND EXTRACT(MONTH FROM complete_at) = $1 
            AND EXTRACT(YEAR FROM complete_at) = $2
            AND key = 'san_pham_loi'
		GROUP BY 
			stages.stage_id;
	`

	// Thực hiện truy vấn với tham số tháng và năm
	err := cockroach.Select(ctx, query, month, year).ScanAll(&result)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return result, nil
}

func NewStatisticsRepo() StatisticsRepo {
	return &statisticsRepo{}
}
