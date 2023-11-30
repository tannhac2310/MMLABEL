package production_order_stage_device

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type EditProductionOrderStageDeviceOpts struct {
	ID                string
	DeviceID          string
	Quantity          int64	
	ProcessStatus     enum.ProductionOrderStageDeviceStatus
	Status            enum.CommonStatus
	Responsible       []string
	NotUpdateQuantity bool
	AssignedQuantity  int64
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
}

type Service interface {
	Edit(ctx context.Context, opt *EditProductionOrderStageDeviceOpts) error
	Create(ctx context.Context, opt *CreateProductionOrderStageDeviceOpts) (string, error)
	Deletes(ctx context.Context, ids []string) error
	Find(ctx context.Context, opt *FindProductionOrderStageDeviceOpts) ([]*repository.ProductionOrderStageDeviceData, error)
	FindEventLog(ctx context.Context, opt *FindEventLogOpts) ([]*repository.EventLogData, error)
}
type productionOrderStageDeviceService struct {
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo
}

func (p productionOrderStageDeviceService) Edit(ctx context.Context, opt *EditProductionOrderStageDeviceOpts) error {
	table := model.ProductionOrderStageDevice{}

	updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderStageFieldID, opt.ID)
	updater.Set(model.ProductionOrderStageDeviceFieldDeviceID, opt.DeviceID)
	if !opt.NotUpdateQuantity {
		updater.Set(model.ProductionOrderStageDeviceFieldQuantity, opt.Quantity)
	}

	updater.Set(model.ProductionOrderStageDeviceFieldProcessStatus, opt.ProcessStatus)
	updater.Set(model.ProductionOrderStageDeviceFieldStatus, opt.Status)
	updater.Set(model.ProductionOrderStageDeviceFieldResponsible, opt.Responsible)
	if opt.AssignedQuantity > 0 {
		updater.Set(model.ProductionOrderStageDeviceFieldAssignedQuantity, opt.AssignedQuantity)
	}
	
	updater.Set(model.ProductionOrderStageDeviceFieldUpdatedAt, time.Now())

	err := cockroach.UpdateFields(ctx, updater)
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

func (p productionOrderStageDeviceService) Find(ctx context.Context, opt *FindProductionOrderStageDeviceOpts) ([]*repository.ProductionOrderStageDeviceData, error) {
	productionOrderStageDevices, err := p.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
		ProductionOrderStageID: opt.ProductionOrderStageID,
		ProductionOrderID:      opt.ProductionOrderID,
		Limit:                  10000,
	})
	if err != nil {
		return nil, err
	}
	return productionOrderStageDevices, nil
}

func NewService(productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo) Service {
	return &productionOrderStageDeviceService{
		productionOrderStageDeviceRepo: productionOrderStageDeviceRepo,
	}

}
