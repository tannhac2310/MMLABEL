package production_order_stage_device

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type EditProductionOrderStageDeviceOpts struct {
	ID                  string
	DeviceID            string
	Quantity            int64
	ProcessStatus       enum.ProductionOrderStageDeviceStatus
	Status              enum.CommonStatus
	Responsible         []string
	AssignedQuantity    int64
	UserID              string
	Settings            *Settings
	Note                string
	SanPhamLoi          int64
	EstimatedStartAt    time.Time
	EstimatedCompleteAt time.Time
}
type Settings struct {
	DefectiveError string
	Description    string
}
type CreateProductionOrderStageDeviceOpts struct {
	ProductionOrderStageID string
	DeviceID               string
	Quantity               int64
	ProcessStatus          enum.ProductionOrderStageDeviceStatus
	Status                 enum.CommonStatus
	Responsible            []string
	Settings               map[string]interface{}
	Note                   string
	AssignedQuantity       int64
	EstimatedStartAt       time.Time
	EstimatedCompleteAt    time.Time
}

type FindProductionOrderStageDeviceOpts struct {
	ProductionOrderStageIDs      []string
	ProductionOrderStageStatuses []enum.ProductionOrderStageStatus
	ProductionOrderIDs           []string
	Responsible                  []string
	DeviceIDs                    []string
	ID                           string
	IDs                          []string
	ProcessStatuses              []enum.ProductionOrderStageDeviceStatus
	Limit                        int64
	Offset                       int64
	StartAt                      time.Time
	CompleteAt                   time.Time
	EstimatedStartAtFrom         time.Time
	EstimatedStartAtTo           time.Time
	StageIDs                     []string
}
type ProductionOrderStageDeviceData struct {
	*repository.ProductionOrderStageDeviceData
	Responsible         []*repository.ProductionOrderStageResponsibleData
	ProductionOrderData *repository.ProductionOrderData
}
type Service interface {
	Edit(ctx context.Context, opt *EditProductionOrderStageDeviceOpts) error
	Create(ctx context.Context, opt *CreateProductionOrderStageDeviceOpts) (string, error)
	Deletes(ctx context.Context, ids []string) error
	Find(ctx context.Context, opt *FindProductionOrderStageDeviceOpts) ([]*ProductionOrderStageDeviceData, *repository.CountResult, error)
	FindEventLog(ctx context.Context, opt *FindEventLogOpts) ([]*repository.EventLogData, error)
	FindProcessDeviceHistory(ctx context.Context, opt *FindProcessDeviceHistoryOpts, sort *repository.Sort, limit, offset int64) ([]*repository.DeviceProgressStatusHistoryData, *repository.CountResult, error)
	EditDeviceProcessHistoryIsSolved(ctx context.Context, opt *EditDeviceProcessHistoryIsSolvedOpts) error
	FindAvailabilityTime(ctx context.Context, opt *FindLostTimeOpts) (*AvailabilityTime, error)
	UpdateProcessStatus(ctx context.Context, opt *UpdateProcessStatusOpts) error
}
type productionOrderStageDeviceService struct {
	productionOrderRepo              repository.ProductionOrderRepo
	productionOrderStageDeviceRepo   repository.ProductionOrderStageDeviceRepo
	sDeviceProgressStatusHistoryRepo repository.DeviceProgressStatusHistoryRepo
	sDeviceWorkingHistoryRepo        repository.DeviceWorkingHistoryRepo
	stageResponsibleRepo             repository.ProductionOrderStageResponsibleRepo
}

func (p productionOrderStageDeviceService) Edit(ctx context.Context, opt *EditProductionOrderStageDeviceOpts) error {
	userID := opt.UserID
	table := model.ProductionOrderStageDevice{}
	tableProductProgress := model.DeviceProgressStatusHistory{}
	tableDevice := model.Device{}
	// find by id and check if it is existed
	data, err := p.productionOrderStageDeviceRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("p.productionOrderStageDeviceRepo.FindByID: %w", err)
	}
	// todo check error != notfound

	if data != nil && data.ProcessStatus != opt.ProcessStatus {
		// find lasted status of device
		lasted, err := p.sDeviceProgressStatusHistoryRepo.FindProductionOrderStageDeviceID(ctx, data.ID, data.DeviceID)
		// if lastederr != nil {
		// 	return fmt.Errorf("p.sDeviceProgressStatusHistoryRepo.FindProductionOrderStageDeviceID: %w", err)
		// }
		fmt.Println("userID===============>>>>lasted", err, lasted, data.ID, data.DeviceID)
		if lasted != nil && lasted.IsResolved == 0 && opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart && (lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed || lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusPause) {
			updaterHistory := cockroach.NewUpdater(tableProductProgress.TableName(), model.DeviceProgressStatusHistoryFieldID, lasted.ID)
			updaterHistory.Set(model.DeviceProgressStatusHistoryFieldUpdatedBy, userID)
			updaterHistory.Set(model.DeviceProgressStatusHistoryFieldUpdatedAt, time.Now())
			updaterHistory.Set(model.DeviceProgressStatusHistoryFieldIsResolved, 1)
			err := cockroach.UpdateFields(ctx, updaterHistory)
			if err != nil {
				// return fmt.Errorf("updaterHistory.cockroach.UpdateFields: %w", err)
			}
			fmt.Println("lasted.ErrorCode.String: ", lasted.ErrorCode.String)
			if lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed && lasted.ErrorCode.String == "MA1" {
				updaterDevice := cockroach.NewUpdater(tableDevice.TableName(), model.DeviceFieldID, data.DeviceID)
				updaterDevice.Set(model.DeviceFieldStatus, enum.CommonStatusActive)
				err := cockroach.UpdateFields(ctx, updaterDevice)
				if err != nil {
					// return fmt.Errorf("updaterDevice.cockroach.UpdateFields: %w", err)
				}
			}
		}
		if opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed && opt.Note == "MA1" {
			fmt.Println("opt.Note: ", opt.Note)
			updaterDevice := cockroach.NewUpdater(tableDevice.TableName(), model.DeviceFieldID, data.DeviceID)
			updaterDevice.Set(model.DeviceFieldStatus, enum.CommonStatusDamage)
			err := cockroach.UpdateFields(ctx, updaterDevice)
			if err != nil {
				// return fmt.Errorf("updaterDevice.cockroach.UpdateFields: %w", err)
			}
		}
		//  insert DeviceProgressStatusHistory
		modelData := &model.DeviceProgressStatusHistory{
			ID:                           idutil.ULIDNow(),
			ProductionOrderStageDeviceID: data.ID,
			DeviceID:                     data.DeviceID,
			ProcessStatus:                opt.ProcessStatus,
			CreatedBy:                    cockroach.String(userID),
			CreatedAt:                    time.Now(),
		}
		if opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed || opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusPause {
			modelData.IsResolved = 0
			modelData.ErrorCode = cockroach.String(opt.Note)
			modelData.ErrorReason = cockroach.String(opt.Settings.DefectiveError)
			modelData.Description = cockroach.String(opt.Settings.Description)
		}

		err = p.sDeviceProgressStatusHistoryRepo.Insert(ctx, modelData)

		if err != nil {
			return fmt.Errorf("p.sDeviceProgressStatusHistoryRepo.Insert: %w", err)
		}
	}

	updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderStageFieldID, opt.ID)

	if opt.DeviceID != "" {
		updater.Set(model.ProductionOrderStageDeviceFieldDeviceID, opt.DeviceID)
	}
	//return fmt.Errorf("This ID not exists: %s", opt.Quantity)
	if opt.Quantity > 0 {
		updater.Set(model.ProductionOrderStageDeviceFieldQuantity, opt.Quantity)
	}

	if opt.ProcessStatus != 0 {
		updater.Set(model.ProductionOrderStageDeviceFieldProcessStatus, opt.ProcessStatus)
	}

	if opt.Status != 0 {
		updater.Set(model.ProductionOrderStageDeviceFieldStatus, opt.Status)

	}

	if opt.Responsible != nil {
		updater.Set(model.ProductionOrderStageDeviceFieldResponsible, opt.Responsible)
	}

	if opt.AssignedQuantity > 0 {
		updater.Set(model.ProductionOrderStageDeviceFieldAssignedQuantity, opt.AssignedQuantity)
	}
	if opt.Note != "" {
		updater.Set(model.ProductionOrderStageDeviceFieldNote, cockroach.String(opt.Note))
	}
	settings := data.Settings
	if opt.SanPhamLoi > 0 {
		if settings == nil {
			settings = make(map[string]interface{})
		}
		settings["san_pham_loi"] = opt.SanPhamLoi
		updater.Set(model.ProductionOrderStageDeviceFieldSettings, settings)
	}

	updater.Set(model.ProductionOrderFieldEstimatedStartAt, opt.EstimatedStartAt)
	updater.Set(model.ProductionOrderFieldEstimatedCompleteAt, opt.EstimatedCompleteAt)

	updater.Set(model.ProductionOrderStageDeviceFieldUpdatedAt, time.Now())

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err = cockroach.UpdateFields(ctx, updater)
		if err != nil {
			return fmt.Errorf("update stage device %w", err)
		}
		if opt.Responsible != nil {
			err = p.stageResponsibleRepo.Delete(ctx2, opt.ID)
			if err != nil {
				return fmt.Errorf("delete stage responsible %w", err)
			}
			for _, userID := range opt.Responsible {
				err = p.stageResponsibleRepo.Insert(ctx2, &model.ProductionOrderStageResponsible{
					ID:              idutil.ULIDNow(),
					POStageDeviceID: opt.ID,
					UserID:          userID,
				})
				if err != nil {
					return fmt.Errorf("insert stage responsible %w", err)
				}
			}
		}
		return nil
	})

	if errTx != nil {
		return fmt.Errorf("po stage device edit: %w", errTx)
	}
	return nil
}

func (p productionOrderStageDeviceService) Create(ctx context.Context, opt *CreateProductionOrderStageDeviceOpts) (string, error) {
	cnt, err := p.productionOrderStageDeviceRepo.CountRows(ctx)
	if err != nil {
		return "", fmt.Errorf("p.productionOrderStageDeviceRepo.Count: %w", err)
	}
	id := fmt.Sprintf("%d", cnt+1)

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err = p.productionOrderStageDeviceRepo.Insert(ctx2, &model.ProductionOrderStageDevice{
			ID:                     id,
			ProductionOrderStageID: opt.ProductionOrderStageID,
			DeviceID:               opt.DeviceID,
			Quantity:               opt.Quantity,
			ProcessStatus:          opt.ProcessStatus,
			Status:                 opt.Status,
			Settings:               opt.Settings,
			Note:                   cockroach.String(opt.Note),
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
			Responsible:            opt.Responsible,
			AssignedQuantity:       opt.AssignedQuantity,
			EstimatedCompleteAt:    cockroach.Time(opt.EstimatedCompleteAt),
			EstimatedStartAt:       cockroach.Time(opt.EstimatedStartAt),
		})
		if err != nil {
			return err
		}

		for _, userID := range opt.Responsible {
			err = p.stageResponsibleRepo.Insert(ctx2, &model.ProductionOrderStageResponsible{
				ID:              idutil.ULIDNow(),
				POStageDeviceID: id,
				UserID:          userID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if errTx != nil {
		return "", fmt.Errorf("po stage device create: %w %s cnt.Count %s", errTx, id, cnt)
	}

	return id, nil
}

func (p productionOrderStageDeviceService) Deletes(ctx context.Context, ids []string) error {
	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		for _, id := range ids {
			err := p.stageResponsibleRepo.Delete(ctx2, id)
			if err != nil {
				return err
			}
		}
		return p.productionOrderStageDeviceRepo.SoftDeletes(ctx2, ids)
	})

	if errTx != nil {
		return fmt.Errorf("po stage device delete: %w", errTx)
	}

	return nil
}

func (p productionOrderStageDeviceService) FindProcessDeviceHistory(ctx context.Context, opt *FindProcessDeviceHistoryOpts, sort *repository.Sort, limit, offset int64) ([]*repository.DeviceProgressStatusHistoryData, *repository.CountResult, error) {
	data, err := p.sDeviceProgressStatusHistoryRepo.Search(ctx, &repository.SearchDeviceProgressStatusHistoryOpts{
		ProcessStatus: opt.ProcessStatus,
		ErrorCodes:    opt.ErrorCodes,
		CreatedFrom:   opt.CreatedFrom,
		CreatedTo:     opt.CreatedTo,
		DeviceID:      opt.DeviceID,
		DeviceIDs:     opt.DeviceIDs,
		IsResolved:    opt.IsResolved,
		Limit:         limit,
		Offset:        offset,
		Sort:          sort,
	})
	if err != nil {
		return nil, nil, err
	}

	total, err := p.sDeviceProgressStatusHistoryRepo.Count(ctx, &repository.SearchDeviceProgressStatusHistoryOpts{
		ProcessStatus: opt.ProcessStatus,
		ErrorCodes:    opt.ErrorCodes,
		CreatedFrom:   opt.CreatedFrom,
		CreatedTo:     opt.CreatedTo,
		DeviceID:      opt.DeviceID,
		IsResolved:    opt.IsResolved,
	})
	if err != nil {
		return nil, nil, err
	}

	return data, total, nil
}

type FindProcessDeviceHistoryOpts struct {
	ProcessStatus []int8
	DeviceID      string
	DeviceIDs     []string
	IsResolved    int16
	ErrorCodes    []string
	CreatedFrom   time.Time
	CreatedTo     time.Time
}
type EditDeviceProcessHistoryIsSolvedOpts struct {
	UserID string
	ID     string
}

func (p productionOrderStageDeviceService) EditDeviceProcessHistoryIsSolved(ctx context.Context, opt *EditDeviceProcessHistoryIsSolvedOpts) error {
	userID := opt.UserID
	tableProductProgress := model.DeviceProgressStatusHistory{}
	lasted, err := p.sDeviceProgressStatusHistoryRepo.FindByID(ctx, opt.ID)
	tableDevice := model.Device{}
	if err != nil {
		return fmt.Errorf("p.sDeviceProgressStatusHistoryRepo.FindByID: %w", err)
	}
	if lasted == nil {
		return fmt.Errorf("This ID not exists: %w", err)
	}
	if lasted.IsResolved == 1 {
		return fmt.Errorf("This ID is solved: %w", err)
	}
	if lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed || lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusPause {
		updaterHistory := cockroach.NewUpdater(tableProductProgress.TableName(), model.DeviceProgressStatusHistoryFieldID, lasted.ID)
		updaterHistory.Set(model.DeviceProgressStatusHistoryFieldUpdatedBy, userID)
		updaterHistory.Set(model.DeviceProgressStatusHistoryFieldUpdatedAt, time.Now())
		updaterHistory.Set(model.DeviceProgressStatusHistoryFieldIsResolved, 1)
		err := cockroach.UpdateFields(ctx, updaterHistory)
		if err != nil {
			return fmt.Errorf("updaterHistory.cockroach.UpdateFields: %w", err)
		}
		if lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed && lasted.ErrorCode.String == "MA1" {
			updaterDevice := cockroach.NewUpdater(tableDevice.TableName(), model.DeviceFieldID, lasted.DeviceID)
			updaterDevice.Set(model.DeviceFieldStatus, enum.CommonStatusActive)
			err := cockroach.UpdateFields(ctx, updaterDevice)
			if err != nil {
				// return fmt.Errorf("updaterDevice.cockroach.UpdateFields: %w", err)
			}
		}
	}
	return nil
}

func (p productionOrderStageDeviceService) Find(ctx context.Context, opt *FindProductionOrderStageDeviceOpts) ([]*ProductionOrderStageDeviceData, *repository.CountResult, error) {
	searchOpts := &repository.SearchProductionOrderStageDevicesOpts{
		ID:                           opt.ID,
		IDs:                          opt.IDs,
		Responsible:                  opt.Responsible,
		ProductionOrderStageIDs:      opt.ProductionOrderStageIDs,
		ProductionOrderIDs:           opt.ProductionOrderIDs,
		ProcessStatuses:              opt.ProcessStatuses,
		DeviceIDs:                    opt.DeviceIDs,
		ProductionOrderStageStatuses: opt.ProductionOrderStageStatuses,
		StartAt:                      opt.StartAt,
		CompleteAt:                   opt.CompleteAt,
		StageIDs:                     opt.StageIDs,
		EstimatedStartAtFrom:         opt.EstimatedStartAtFrom,
		EstimatedStartAtTo:           opt.EstimatedStartAtTo,
		Limit:                        opt.Limit,
		Offset:                       opt.Offset,
		Sort:                         nil,
	}

	productionOrderStageDevices, err := p.productionOrderStageDeviceRepo.Search(ctx, searchOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("p.productionOrderStageDeviceRepo.Search: %w", err)
	}

	total, err := p.productionOrderStageDeviceRepo.Count(ctx, searchOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("p.productionOrderStageDeviceRepo.Count: %w", err)
	}

	poIds := make([]string, 0, len(productionOrderStageDevices))
	for _, posd := range productionOrderStageDevices {
		poIds = append(poIds, posd.ProductionOrderID)
	}

	poData, err := p.productionOrderRepo.Search(ctx, &repository.SearchProductionOrdersOpts{
		IDs:    poIds,
		Limit:  int64(len(poIds)),
		Offset: 0,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("p.productionOrderRepo.Search: %w", err)
	}

	poDataMap := make(map[string]*repository.ProductionOrderData)
	for _, po := range poData {
		poDataMap[po.ID] = po
	}

	result := make([]*ProductionOrderStageDeviceData, 0, len(productionOrderStageDevices))
	for _, posd := range productionOrderStageDevices {
		responsible, err := p.stageResponsibleRepo.Search(ctx, &repository.SearchProductionOrderStageResponsibleOpts{
			POStageDeviceIDs: []string{posd.ID},
			Limit:            1000,
			Offset:           0,
		})

		if err != nil {
			return nil, nil, fmt.Errorf("p.stageResponsibleRepo.Search: %w", err)
		}

		po, ok := poDataMap[posd.ProductionOrderID]
		if !ok {
			return nil, nil, fmt.Errorf("production order not found: %s", posd.ProductionOrderID)
		}

		result = append(result, &ProductionOrderStageDeviceData{
			ProductionOrderStageDeviceData: posd,
			Responsible:                    responsible,
			ProductionOrderData:            po,
		})
	}

	return result, total, nil
}

type UpdateProcessStatusOpts struct {
	ProductionOrderStageDeviceID string
	ProcessStatus                enum.ProductionOrderStageDeviceStatus
	UserID                       string
	Settings                     *SettingsData
}

func (p productionOrderStageDeviceService) UpdateProcessStatus(ctx context.Context, opt *UpdateProcessStatusOpts) error {
	userID := opt.UserID
	//table := model.ProductionOrderStageDevice{}
	tableProductProgress := model.DeviceProgressStatusHistory{}
	tableDevice := model.Device{}

	data, err := p.productionOrderStageDeviceRepo.FindByID(ctx, opt.ProductionOrderStageDeviceID)
	if err != nil {
		return fmt.Errorf("p.productionOrderStageDeviceRepo.FindByID: %w", err)
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		device := model.ProductionOrderStageDevice{}
		updater := cockroach.NewUpdater(device.TableName(), model.ProductionOrderStageDeviceFieldID, opt.ProductionOrderStageDeviceID)
		updater.Set(model.ProductionOrderStageDeviceFieldProcessStatus, opt.ProcessStatus)
		updater.Set(model.ProductionOrderStageDeviceFieldUpdatedAt, time.Now())

		if opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart {
			updater.Set(model.ProductionOrderStageDeviceFieldStartAt, time.Now())
		}

		if opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete {
			updater.Set(model.ProductionOrderStageDeviceFieldCompleteAt, time.Now())
		}

		err := cockroach.UpdateFields(ctx, updater)
		if err != nil {
			return fmt.Errorf("update process status: %w", err)
		}

		// write log
		if data.ProcessStatus != opt.ProcessStatus {
			// find lasted status of device
			lasted, _ := p.sDeviceProgressStatusHistoryRepo.FindProductionOrderStageDeviceID(ctx, data.ID, data.DeviceID)
			fmt.Println("userID===============>>>>lasted", err, lasted, data.ID, data.DeviceID)
			if lasted != nil && lasted.IsResolved == 0 && opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart && (lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed || lasted.ProcessStatus == enum.ProductionOrderStageDeviceStatusPause) {
				updaterHistory := cockroach.NewUpdater(tableProductProgress.TableName(), model.DeviceProgressStatusHistoryFieldID, lasted.ID)
				updaterHistory.Set(model.DeviceProgressStatusHistoryFieldUpdatedBy, userID)
				updaterHistory.Set(model.DeviceProgressStatusHistoryFieldUpdatedAt, time.Now())
				updaterHistory.Set(model.DeviceProgressStatusHistoryFieldIsResolved, 1)
				err := cockroach.UpdateFields(ctx, updaterHistory)
				if err != nil {
					return fmt.Errorf("lỗi cập nhật xử lý máy: %w", err)
				}
			}
			// MA1: Máy bị lỗi.
			// Cập nhật trạng thái máy thành hỏng
			// TODO fix hardcode
			if opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed && opt.Settings.DefectiveError == "MA1" {
				updaterDevice := cockroach.NewUpdater(tableDevice.TableName(), model.DeviceFieldID, data.DeviceID)
				updaterDevice.Set(model.DeviceFieldStatus, enum.CommonStatusDamage)
				updaterDevice.Set(model.DeviceFieldUpdatedAt, time.Now())
				updaterDevice.Set(model.DeviceFieldOptionID, data.ID)
				updaterDevice.Set(model.DeviceFieldData, model.SettingsData{
					DefectiveError:               opt.Settings.DefectiveError,
					DefectiveReason:              opt.Settings.DefectiveReason,
					Description:                  opt.Settings.Description,
					ProductionOrderStageID:       data.ProductionOrderStageID,
					ProductionOrderStageDeviceID: data.ID,
				})

				err := cockroach.UpdateFields(ctx, updaterDevice)
				if err != nil {
					return fmt.Errorf("cập nhật trạng thái máy thành hỏng thất bại: %w", err)
				}
			}
			//  insert DeviceProgressStatusHistory
			modelData := &model.DeviceProgressStatusHistory{
				ID:                           idutil.ULIDNow(),
				ProductionOrderStageDeviceID: data.ID,
				DeviceID:                     data.DeviceID,
				ProcessStatus:                opt.ProcessStatus,
				CreatedBy:                    cockroach.String(userID),
				CreatedAt:                    time.Now(),
			}

			if opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed || opt.ProcessStatus == enum.ProductionOrderStageDeviceStatusPause {
				modelData.IsResolved = 0
				modelData.ErrorCode = cockroach.String(opt.Settings.DefectiveError)
				modelData.ErrorReason = cockroach.String(opt.Settings.DefectiveReason)
				modelData.Description = cockroach.String(opt.Settings.Description)
			}

			err = p.sDeviceProgressStatusHistoryRepo.Insert(ctx, modelData)

			if err != nil {
				return fmt.Errorf("p.sDeviceProgressStatusHistoryRepo.Insert: %w", err)
			}
		}

		return nil
	})

	if errTx != nil {
		return fmt.Errorf("cập nhật trạng thái lệnh làm việc thất bại: %w", errTx)
	}

	return nil
}

type SettingsData struct {
	DefectiveError  string
	DefectiveReason string
	Description     string
}
type EditSettingsOpts struct {
	ProductionOrderStageDeviceID string
	Settings                     *SettingsData
}

func NewService(
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo,
	sDeviceProgressStatusHistoryRepo repository.DeviceProgressStatusHistoryRepo,
	sDeviceWorkingHistoryRepo repository.DeviceWorkingHistoryRepo,
	stageResponsibleRepo repository.ProductionOrderStageResponsibleRepo,
	productionOrderRepo repository.ProductionOrderRepo,
) Service {
	return &productionOrderStageDeviceService{
		productionOrderStageDeviceRepo:   productionOrderStageDeviceRepo,
		sDeviceProgressStatusHistoryRepo: sDeviceProgressStatusHistoryRepo,
		sDeviceWorkingHistoryRepo:        sDeviceWorkingHistoryRepo,
		stageResponsibleRepo:             stageResponsibleRepo,
		productionOrderRepo:              productionOrderRepo,
	}

}
