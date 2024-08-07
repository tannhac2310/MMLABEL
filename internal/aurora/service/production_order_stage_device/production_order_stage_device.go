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
	ID               string
	DeviceID         string
	Quantity         int64
	ProcessStatus    enum.ProductionOrderStageDeviceStatus
	Status           enum.CommonStatus
	Responsible      []string
	AssignedQuantity int64
	UserID           string
	Settings         *Settings
	Note             string
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
}

type FindProductionOrderStageDeviceOpts struct {
	ProductionOrderStageID string
	ProductionOrderID      string
	ProcessStatus          enum.ProductionOrderStageDeviceStatus
	Limit                  int64
	Offset                 int64
}

type Service interface {
	Edit(ctx context.Context, opt *EditProductionOrderStageDeviceOpts) error
	Create(ctx context.Context, opt *CreateProductionOrderStageDeviceOpts) (string, error)
	Deletes(ctx context.Context, ids []string) error
	Find(ctx context.Context, opt *FindProductionOrderStageDeviceOpts) ([]*repository.ProductionOrderStageDeviceData, error)
	FindEventLog(ctx context.Context, opt *FindEventLogOpts) ([]*repository.EventLogData, error)
	FindProcessDeviceHistory(ctx context.Context, opt *FindProcessDeviceHistoryOpts, sort *repository.Sort, limit, offset int64) ([]*repository.DeviceProgressStatusHistoryData, *repository.CountResult, error)
	EditDeviceProcessHistoryIsSolved(ctx context.Context, opt *EditDeviceProcessHistoryIsSolvedOpts) error
	FindAvailabilityTime(ctx context.Context, opt *FindLostTimeOpts) (*AvailabilityTime, error)
}
type productionOrderStageDeviceService struct {
	productionOrderStageDeviceRepo   repository.ProductionOrderStageDeviceRepo
	sDeviceProgressStatusHistoryRepo repository.DeviceProgressStatusHistoryRepo
	sDeviceWorkingHistoryRepo        repository.DeviceWorkingHistoryRepo
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

	updater.Set(model.ProductionOrderStageDeviceFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return err
	}
	return nil
}

func (p productionOrderStageDeviceService) Create(ctx context.Context, opt *CreateProductionOrderStageDeviceOpts) (string, error) {
	id := idutil.ULIDNow()
	err := p.productionOrderStageDeviceRepo.Insert(ctx, &model.ProductionOrderStageDevice{
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
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p productionOrderStageDeviceService) Deletes(ctx context.Context, ids []string) error {
	return p.productionOrderStageDeviceRepo.SoftDeletes(ctx, ids)
}

func (p productionOrderStageDeviceService) FindProcessDeviceHistory(ctx context.Context, opt *FindProcessDeviceHistoryOpts, sort *repository.Sort, limit, offset int64) ([]*repository.DeviceProgressStatusHistoryData, *repository.CountResult, error) {
	data, err := p.sDeviceProgressStatusHistoryRepo.Search(ctx, &repository.SearchDeviceProgressStatusHistoryOpts{
		ProcessStatus: opt.ProcessStatus,
		ErrorCodes:    opt.ErrorCodes,
		CreatedFrom:   opt.CreatedFrom,
		CreatedTo:     opt.CreatedTo,
		DeviceID:      opt.DeviceID,
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

func (p productionOrderStageDeviceService) Find(ctx context.Context, opt *FindProductionOrderStageDeviceOpts) ([]*repository.ProductionOrderStageDeviceData, error) {
	productionOrderStageDevices, err := p.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
		ProductionOrderStageID: opt.ProductionOrderStageID,
		ProductionOrderID:      opt.ProductionOrderID,
		ProcessStatus:          opt.ProcessStatus,
		Limit:                  opt.Limit,
		Offset:                 opt.Offset,
	})
	if err != nil {
		return nil, err
	}
	return productionOrderStageDevices, nil
}

func NewService(
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo,
	sDeviceProgressStatusHistoryRepo repository.DeviceProgressStatusHistoryRepo,
	sDeviceWorkingHistoryRepo repository.DeviceWorkingHistoryRepo,
) Service {
	return &productionOrderStageDeviceService{
		productionOrderStageDeviceRepo:   productionOrderStageDeviceRepo,
		sDeviceProgressStatusHistoryRepo: sDeviceProgressStatusHistoryRepo,
		sDeviceWorkingHistoryRepo:        sDeviceWorkingHistoryRepo,
	}

}
