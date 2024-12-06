package ink

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type EditInkOpts struct {
	ID             string
	Name           string
	Code           string
	ProductCodes   []string
	Position       string
	Location       string
	Quantity       float64
	Manufacturer   string
	ColorDetail    map[string]interface{}
	ExpirationDate string // DD-MM-YYYY
	Description    string
	Data           map[string]interface{}
	Status         enum.CommonStatus
	Kho            string
	LoaiMuc        string
	NhaCungCap     string
	TinhTrang      string
	UpdatedBy      string
}

type CreateInkOpts struct {
	Name           string
	Code           string
	ProductCodes   []string
	Position       string
	Location       string
	Manufacturer   string
	ColorDetail    map[string]interface{}
	Quantity       float64
	ExpirationDate string // DD-MM-YYYY
	Description    string
	Data           map[string]interface{}
	Status         enum.CommonStatus
	Kho            string
	LoaiMuc        string
	NhaCungCap     string
	TinhTrang      string
	CreatedBy      string
}

type FindInkOpts struct {
	Name   string
	ID     string
	Code   string
	NotIDs []string
	Status enum.CommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditInkOpts) error
	Create(ctx context.Context, opt *CreateInkOpts) (string, error)
	Delete(ctx context.Context, id string) error
	MixInk(ctx context.Context, opt *CreateInkMixingOpts) (string, error)
	EditInkMixing(ctx context.Context, opt *EditInkMixingOpts) error
	FindInkMixing(ctx context.Context, opt *FindInkMixingOpts) ([]*InkMixingData, *repository.CountResult, error)
	Find(ctx context.Context, opt *FindInkOpts, sort *repository.Sort, limit, offset int64) ([]*InkData, *repository.CountResult, error)
	CalculateInkQuantity(ctx context.Context, inkID string) (float64, float64, float64, error)
}

type inkService struct {
	inkRepo                         repository.InkRepo
	inkReturnRepo                   repository.InkReturnRepo
	inkReturnDetailRepo             repository.InkReturnDetailRepo
	inkExportDetailRepo             repository.InkExportDetailRepo
	inkImportDetailRepo             repository.InkImportDetailRepo
	historyRepo                     repository.HistoryRepo
	inkMixingRepo                   repository.InkMixingRepo
	inkMixingDetailRepo             repository.InkMixingDetailRepo
	productionOrderDeviceConfigRepo repository.ProductionOrderDeviceConfigRepo
}

func (p inkService) CalculateInkQuantity(ctx context.Context, inkID string) (float64, float64, float64, error) {

	inkData, err := p.inkRepo.FindByID(ctx, inkID)
	if err != nil {
		return 0, 0, 0, err
	}
	// get ink export quantity
	exportQuantity, err := p.calculateInkExportQuantity(ctx, inkID)
	if err != nil {
		return 0, 0, 0, err
	}
	// get ink import quantity
	importQuantity, err := p.calculateInkImportQuantity(ctx, inkID)
	if err != nil {
		return 0, 0, 0, err
	}
	// get ink return quantity
	returnQuantity, err := p.calculateInkReturnQuantity(ctx, inkID)
	if err != nil {
		return 0, 0, 0, err
	}
	// calculate ink quantity in stock
	quantityInStock := inkData.Quantity + importQuantity - exportQuantity - returnQuantity
	return quantityInStock, importQuantity, exportQuantity, nil
}

func (p inkService) Edit(ctx context.Context, opt *EditInkOpts) error {
	table := model.Ink{}
	// find ink by id
	ink, err := p.inkRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return err
	}
	fmt.Println(ink)
	// write history
	// old value
	b, _ := json.Marshal(ink)
	oldValue := map[string]interface{}{}
	b, _ = json.Marshal(opt)
	newValue := map[string]interface{}{}
	_ = json.Unmarshal(b, &newValue)
	err = json.Unmarshal(b, &oldValue)

	err = p.historyRepo.Insert(ctx, &model.History{
		ID:        time.Now().UnixNano(),
		Table:     table.TableName(),
		RowID:     opt.ID,
		OldValue:  oldValue,
		NewValue:  newValue,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("p.historyRepo.Insert: %w", err)
	}

	updater := cockroach.NewUpdater(table.TableName(), model.InkFieldID, opt.ID)

	updater.Set(model.InkFieldName, opt.Name)
	updater.Set(model.InkFieldCode, opt.Code)
	updater.Set(model.InkFieldProductCodes, opt.ProductCodes)
	updater.Set(model.InkFieldPosition, opt.Position)
	updater.Set(model.InkFieldLocation, opt.Location)
	updater.Set(model.InkFieldQuantity, opt.Quantity)
	updater.Set(model.InkFieldManufacturer, opt.Manufacturer)
	updater.Set(model.InkFieldColorDetail, opt.ColorDetail)
	updater.Set(model.InkFieldExpirationDate, opt.ExpirationDate)
	updater.Set(model.InkFieldStatus, opt.Status)
	updater.Set(model.InkFieldDescription, opt.Description)
	updater.Set(model.InkFieldData, opt.Data)
	updater.Set(model.InkFieldUpdatedBy, opt.UpdatedBy)
	//kho
	//loai_muc
	//nha_cung_cap
	//tinh_trang
	updater.Set(model.InkFieldKho, opt.Kho)
	updater.Set(model.InkFieldLoaiMuc, opt.LoaiMuc)
	updater.Set(model.InkFieldNhaCungCap, opt.NhaCungCap)
	updater.Set(model.InkFieldTinhTrang, opt.TinhTrang)

	updater.Set(model.InkFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return err
	}
	return nil
}

func (p inkService) Create(ctx context.Context, opt *CreateInkOpts) (string, error) {
	id := idutil.ULIDNow()
	err := p.inkRepo.Insert(ctx, &model.Ink{
		ID:             id,
		Name:           opt.Name,
		Code:           opt.Code,
		ProductCodes:   opt.ProductCodes,
		Position:       opt.Position,
		Location:       opt.Location,
		Manufacturer:   opt.Manufacturer,
		ColorDetail:    opt.ColorDetail,
		Quantity:       opt.Quantity,
		ExpirationDate: opt.ExpirationDate,
		Description:    cockroach.String(opt.Description),
		Kho:            opt.Kho,
		LoaiMuc:        opt.LoaiMuc,
		NhaCungCap:     opt.NhaCungCap,
		TinhTrang:      opt.TinhTrang,
		Data:           opt.Data,
		Status:         opt.Status,
		CreatedBy:      opt.CreatedBy,
		UpdatedBy:      opt.CreatedBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p inkService) Delete(ctx context.Context, id string) error {
	return p.inkRepo.SoftDelete(ctx, id)
}

type InkData struct {
	*repository.InkData
	MixingData                      *InkMixingData
	ProductionOrderDeviceConfigData []*repository.ProductionOrderDeviceConfigData
}

func (p inkService) Find(ctx context.Context, opt *FindInkOpts, sort *repository.Sort, limit, offset int64) ([]*InkData, *repository.CountResult, error) {
	filter := &repository.SearchInkOpts{
		Name:   opt.Name,
		Status: opt.Status,
		ID:     opt.ID,
		NotIDs: opt.NotIDs,
		Code:   opt.Code,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	inks, err := p.inkRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	mixingIDs := make([]string, 0)
	inkIds := make([]string, 0)
	for _, ink := range inks {
		if ink.MixingID.String != "" {
			mixingIDs = append(mixingIDs, ink.MixingID.String)
		}
		inkIds = append(inkIds, ink.ID)
	}

	// find in production_order_device_config by color
	productionOrderDeviceConfigs, err := p.productionOrderDeviceConfigRepo.Search(ctx, &repository.SearchProductionOrderDeviceConfigOpts{
		InkIDs: inkIds,
		Limit:  10000,
	})
	if err != nil {
		return nil, nil, err
	}

	// map production_order_device_config by ink_id
	productionOrderDeviceConfigMap := make(map[string][]*repository.ProductionOrderDeviceConfigData)
	for _, f := range productionOrderDeviceConfigs {
		productionOrderDeviceConfigMap[f.InkID.String] = append(productionOrderDeviceConfigMap[f.InkID.String], f)
	}
	// get ink mixing detail
	inkMixing, _, err := p.FindInkMixing(ctx, &FindInkMixingOpts{
		IDs:    mixingIDs,
		Limit:  int64(len(mixingIDs)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("inkMixingDetailsinkMixingDetails============================>", inkMixing)
	// map ink mixing detail
	inkMixingMap := make(map[string]*InkMixingData)
	for _, f := range inkMixing {
		inkMixingMap[f.ID] = f
	}

	results := make([]*InkData, 0)
	for _, ink := range inks {
		results = append(results, &InkData{
			InkData:                         ink,
			MixingData:                      inkMixingMap[ink.MixingID.String],
			ProductionOrderDeviceConfigData: productionOrderDeviceConfigMap[ink.ID],
		})
	}

	total, err := p.inkRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return results, total, nil
}

func (p inkService) calculateInkExportQuantity(ctx context.Context, inkID string) (float64, error) {
	// get ink export detail
	inkExportDetails, err := p.inkExportDetailRepo.Search(ctx, &repository.SearchInkExportDetailOpts{
		InkID: inkID,
		Limit: 10000,
	})
	if err != nil {
		return 0, err
	}

	// we dont care performance here, so we can use loop to calculate ink export quantity
	// calculate ink export quantity
	var exportQuantity float64
	for _, inkExportDetail := range inkExportDetails {
		exportQuantity += inkExportDetail.Quantity
	}
	return exportQuantity, nil
}

func (p inkService) calculateInkImportQuantity(ctx context.Context, inkID string) (float64, error) {
	// get ink import detail
	inkImportDetails, err := p.inkImportDetailRepo.Search(ctx, &repository.SearchInkImportDetailOpts{
		ID:    inkID, // when importing, I write ink_import_detail.ID = ink.ID
		Limit: 10000,
	})
	if err != nil {
		return 0, err
	}
	// we dont care performance here, so we can use loop to calculate ink import quantity
	// calculate ink import quantity
	var importQuantity float64
	for _, inkImportDetail := range inkImportDetails {
		importQuantity += inkImportDetail.Quantity
	}
	return importQuantity, nil

}

func (p inkService) calculateInkReturnQuantity(ctx context.Context, inkID string) (float64, error) {
	// get ink return detail
	inkReturnDetails, err := p.inkReturnDetailRepo.Search(ctx, &repository.SearchInkReturnDetailOpts{
		InkID: inkID,
		Limit: 10000,
	})
	if err != nil {
		return 0, err
	}
	// we dont care performance here, so we can use loop to calculate ink return quantity
	// calculate ink return quantity
	var returnQuantity float64
	for _, inkReturnDetail := range inkReturnDetails {
		returnQuantity += inkReturnDetail.Quantity
	}
	return returnQuantity, nil
}

func NewService(
	inkRepo repository.InkRepo,
	inkReturnRepo repository.InkReturnRepo,
	inkReturnDetailRepo repository.InkReturnDetailRepo,
	inkExportDetailRepo repository.InkExportDetailRepo,
	inkImportDetailRepo repository.InkImportDetailRepo,
	historyRepo repository.HistoryRepo,
	inkMixingRepo repository.InkMixingRepo,
	inkMixingDetailRepo repository.InkMixingDetailRepo,
	productionOrderDeviceConfigRepo repository.ProductionOrderDeviceConfigRepo,
) Service {
	return &inkService{
		inkReturnRepo:                   inkReturnRepo,
		inkReturnDetailRepo:             inkReturnDetailRepo,
		inkExportDetailRepo:             inkExportDetailRepo,
		inkImportDetailRepo:             inkImportDetailRepo,
		inkRepo:                         inkRepo,
		historyRepo:                     historyRepo,
		inkMixingRepo:                   inkMixingRepo,
		inkMixingDetailRepo:             inkMixingDetailRepo,
		productionOrderDeviceConfigRepo: productionOrderDeviceConfigRepo,
	}

}
