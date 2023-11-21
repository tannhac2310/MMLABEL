package ink_export

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

type EditInkExportDetailOpts struct {
}

type EditInkExportOpts struct {
	ID string
	// implement later
}

type CreateInkExportDetailOpts struct {
	InkID       string
	Quantity    float64
	Description string
	Data        map[string]interface{}
}
type CreateInkExportOpts struct {
	Name              string
	Code              string
	ProductionOrderID string
	ExportDate        string
	Description       string
	Data              map[string]interface{}
	CreatedBy         string
	InkExportDetail   []*CreateInkExportDetailOpts
}
type FindInkExportOpts struct {
	Name   string
	Code   string
	ID     string
	Status enum.InventoryCommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditInkExportOpts) error
	Create(ctx context.Context, opt *CreateInkExportOpts) (string, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, opt *FindInkExportOpts, sort *repository.Sort, limit, offset int64) ([]*InkExportData, *repository.CountResult, error)
}

type inkExportService struct {
	inkExportRepo       repository.InkExportRepo
	inkExportDetailRepo repository.InkExportDetailRepo
	inkRepo             repository.InkRepo
}

func (p inkExportService) Edit(ctx context.Context, opt *EditInkExportOpts) error {
	// todo implement later
	panic("not implemented")
}

func (p inkExportService) Create(ctx context.Context, opt *CreateInkExportOpts) (string, error) {

	// write code to insert to ink_export table and insert to ink_export_detail table in transaction
	exportId := idutil.ULIDNow()
	now := time.Now()
	err := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// insert to ink_export
		err1 := p.inkExportRepo.Insert(c, &model.InkExport{
			ID:                exportId,
			Name:              opt.Name,
			Code:              opt.Code,
			ProductionOrderID: opt.ProductionOrderID,
			ExportDate:        cockroach.Time(now),
			Description:       cockroach.String(opt.Description),
			Status:            enum.InventoryCommonStatusStatusCompleted, // default status is completed
			Data:              opt.Data,
			CreatedBy:         opt.CreatedBy,
			UpdatedBy:         opt.CreatedBy,
			CreatedAt:         now,
			UpdatedAt:         now,
		})

		if err1 != nil {
			return fmt.Errorf("error when insert ink import: %w", err1)
		}
		// write code to insert to ink_export_detail table
		for _, inkExportDetail := range opt.InkExportDetail {
			err2 := p.inkExportDetailRepo.Insert(c, &model.InkExportDetail{
				ID:          idutil.ULIDNow(),
				InkExportID: exportId,
				InkID:       inkExportDetail.InkID,
				Quantity:    inkExportDetail.Quantity,
				Description: cockroach.String(inkExportDetail.Description),
				Data:        inkExportDetail.Data,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err2 != nil {
				return fmt.Errorf("error when insert ink import detail: %w", err2)
			}
			// todo check if import status is completed, insert to ink table
			// check if inkExportDetail.InkID exist in ink.ID
			inkData, err := p.inkRepo.FindByID(c, inkExportDetail.InkID)
			if err != nil {
				return fmt.Errorf("error when find ink import detail by id: %w", err)
			}
			// update ink quantity
			inkData.Quantity = inkData.Quantity - inkExportDetail.Quantity
			if inkData.Quantity < 0 {
				return fmt.Errorf("không đủ số lượng để xuất khoa: trong kho còn %v", inkData.Quantity)
			}
			err3 := p.inkRepo.Update(c, inkData.Ink)
			if err3 != nil {
				return fmt.Errorf("error when insert ink: %w", err3)
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}
	return exportId, nil
}

func (p inkExportService) Delete(ctx context.Context, id string) error {
	return p.inkExportRepo.SoftDelete(ctx, id)
}

type InkExportData struct {
	ID                string
	ProductionOrderID string
	Code              string
	Name              string
	Status            enum.InventoryCommonStatus
	Description       string
	Data              map[string]interface{}
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CreatedBy         string
	InkExportDetail   []*InkExportDetail
}

type InkExportDetail struct {
	ID          string
	InkExportID string
	InkID       string
	InkData     *repository.InkData
	Quantity    float64
	ColorDetail map[string]interface{}
	Description string
	Data        map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p inkExportService) Find(ctx context.Context, opt *FindInkExportOpts, sort *repository.Sort, limit, offset int64) ([]*InkExportData, *repository.CountResult, error) {
	filter := &repository.SearchInkExportOpts{
		Name:   opt.Name,
		Status: opt.Status,
		ID:     opt.ID,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	inkExports, err := p.inkExportRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*InkExportData, 0)
	// write code to get ink_export_detail
	for _, inkExport := range inkExports {
		data := &InkExportData{
			ID:                inkExport.ID,
			ProductionOrderID: inkExport.ProductionOrderID,
			Code:              inkExport.Code,
			Name:              inkExport.Name,
			Status:            inkExport.Status,
			Description:       inkExport.Description.String,
			Data:              inkExport.Data,
			CreatedAt:         inkExport.CreatedAt,
			UpdatedAt:         inkExport.UpdatedAt,
			CreatedBy:         inkExport.CreatedBy,
			InkExportDetail:   nil,
		}
		inkExportDetails, err := p.inkExportDetailRepo.Search(ctx, &repository.SearchInkExportDetailOpts{
			InkExportID: inkExport.ID,
			Limit:       10000,
		})
		if err != nil {
			return nil, nil, err
		}
		inkExportDetailResults := make([]*InkExportDetail, 0)
		for _, inkExportDetail := range inkExportDetails {
			inkData, _ := p.inkRepo.FindByID(ctx, inkExportDetail.InkID)
			dataDetail := &InkExportDetail{
				ID:          inkExportDetail.ID,
				InkExportID: inkExportDetail.InkExportID,
				InkID:       inkExportDetail.InkID,
				InkData:     inkData,
				Quantity:    inkData.Quantity,
				ColorDetail: inkData.ColorDetail,
				Description: inkExportDetail.Description.String,
				Data:        inkExportDetail.Data,
				CreatedAt:   inkExportDetail.CreatedAt,
				UpdatedAt:   inkExportDetail.UpdatedAt,
			}
			inkExportDetailResults = append(inkExportDetailResults, dataDetail)
		}
		data.InkExportDetail = inkExportDetailResults
		results = append(results, data)
	}
	total, err := p.inkExportRepo.Count(ctx, filter)
	return results, total, nil
}

func NewService(inkExportRepo repository.InkExportRepo, inkExportDetailRepo repository.InkExportDetailRepo, inkRepo repository.InkRepo) Service {
	return &inkExportService{
		inkExportRepo:       inkExportRepo,
		inkExportDetailRepo: inkExportDetailRepo,
		inkRepo:             inkRepo,
	}

}
