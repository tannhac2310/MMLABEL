package ink_import

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

type EditInkImportDetailOpts struct {
	ID          string
	InkImportID string
	Quantity    float64
	ColorDetail map[string]interface{}
	Description string
	Data        map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EditInkImportOpts struct {
	ID              string
	Name            string
	Code            string
	ProductCodes    []string
	Position        string
	Location        string
	Manufacturer    string
	ColorDetail     map[string]interface{}
	Quantity        float64
	ExpirationDate  time.Time
	Description     string
	Data            map[string]interface{}
	Status          enum.InventoryCommonStatus
	InkImportDetail []*EditInkImportDetailOpts
	UpdatedBy       string
}

type CreateInkImportDetailOpts struct {
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
	Kho            string
	LoaiMuc        string
	NhaCungCap     string
	TinhTrang      string
}
type CreateInkImportOpts struct {
	Name            string
	Code            string
	Description     string
	Data            map[string]interface{}
	InkImportDetail []*CreateInkImportDetailOpts
	CreatedBy       string
}
type FindInkImportOpts struct {
	Name   string
	ID     string
	Status enum.InventoryCommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditInkImportOpts) error
	Create(ctx context.Context, opt *CreateInkImportOpts) (string, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, opt *FindInkImportOpts, sort *repository.Sort, limit, offset int64) ([]*InkImportData, *repository.CountResult, error)
}

type inkImportService struct {
	inkImportRepo       repository.InkImportRepo
	inkImportDetailRepo repository.InkImportDetailRepo
	inkRepo             repository.InkRepo
}

func (p inkImportService) Edit(ctx context.Context, opt *EditInkImportOpts) error {
	// todo implement later
	panic("not implemented")
}

func (p inkImportService) Create(ctx context.Context, opt *CreateInkImportOpts) (string, error) {

	// write code to insert to ink_import table and insert to ink_import_detail table in transaction
	importId := idutil.ULIDNow()
	now := time.Now()
	err := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// insert to ink_import
		err1 := p.inkImportRepo.Insert(c, &model.InkImport{
			ID:          importId,
			Name:        opt.Name,
			Code:        opt.Code,
			Description: cockroach.String(opt.Description),
			Status:      enum.InventoryCommonStatusStatusCompleted, // default status is completed
			Data:        opt.Data,
			CreatedAt:   now,
			UpdatedAt:   now,
			CreatedBy:   opt.CreatedBy,
		})

		if err1 != nil {
			return fmt.Errorf("error when insert ink import: %w", err1)
		}
		// write code to insert to ink_import_detail table
		for _, inkImportDetail := range opt.InkImportDetail {
			inkID := idutil.ULIDNow()
			// find ink by code, manufacturer, expiration_date
			ink, _ := p.inkRepo.Search(ctx, &repository.SearchInkOpts{
				Code:         inkImportDetail.Code,
				Manufacturer: inkImportDetail.Manufacturer,
				Expiration:   inkImportDetail.ExpirationDate,
				Limit:        1,
			})

			if len(ink) > 0 {
				inkData := ink[0].Ink
				inkID = inkData.ID
				inkData.Quantity = inkData.Quantity + inkImportDetail.Quantity
				inkData.UpdatedAt = now
				inkData.UpdatedBy = opt.CreatedBy
				err := p.inkRepo.Update(ctx, inkData)
				if err != nil {
					return fmt.Errorf("error when update ink: %w", err)
				}
			} else {
				// todo check if import status is completed, insert to ink table
				// insert into ink table
				err3 := p.inkRepo.Insert(c, &model.Ink{
					ID:             inkID,
					ImportID:       cockroach.String(importId),
					Name:           inkImportDetail.Name,
					Code:           inkImportDetail.Code,
					ProductCodes:   inkImportDetail.ProductCodes,
					Position:       inkImportDetail.Position,
					Location:       inkImportDetail.Location,
					Manufacturer:   inkImportDetail.Manufacturer,
					ColorDetail:    inkImportDetail.ColorDetail,
					Quantity:       inkImportDetail.Quantity,
					ExpirationDate: inkImportDetail.ExpirationDate,
					Description:    cockroach.String(inkImportDetail.Description),
					Data:           inkImportDetail.Data,
					Status:         enum.CommonStatusActive,
					Kho:            inkImportDetail.Kho,
					LoaiMuc:        inkImportDetail.LoaiMuc,
					NhaCungCap:     inkImportDetail.NhaCungCap,
					TinhTrang:      inkImportDetail.TinhTrang,
					CreatedBy:      opt.CreatedBy,
					UpdatedBy:      opt.CreatedBy,
					CreatedAt:      now,
					UpdatedAt:      now,
				})
				if err3 != nil {
					return fmt.Errorf("error when insert ink: %w", err3)
				}
			}
			err2 := p.inkImportDetailRepo.Insert(c, &model.InkImportDetail{
				ID:             idutil.ULIDNow(),
				InkID:          inkID,
				InkImportID:    importId,
				Name:           inkImportDetail.Name,
				Code:           inkImportDetail.Code,
				ProductCodes:   inkImportDetail.ProductCodes,
				Position:       inkImportDetail.Position,
				Location:       inkImportDetail.Location,
				Manufacturer:   inkImportDetail.Manufacturer,
				ColorDetail:    inkImportDetail.ColorDetail,
				Quantity:       inkImportDetail.Quantity,
				ExpirationDate: inkImportDetail.ExpirationDate,
				Description:    cockroach.String(inkImportDetail.Description),
				Data:           inkImportDetail.Data,
				CreatedAt:      now,
				UpdatedAt:      now,
			})
			if err2 != nil {
				return fmt.Errorf("error when insert ink import detail: %w", err2)
			}

		}

		return nil
	})

	if err != nil {
		return "", err
	}
	return importId, nil
}

func (p inkImportService) Delete(ctx context.Context, id string) error {
	return p.inkImportRepo.SoftDelete(ctx, id)
}

type InkImportData struct {
	ID              string
	Code            string
	Name            string
	Status          enum.InventoryCommonStatus
	Description     string
	Data            map[string]interface{}
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CreatedBy       string
	InkImportDetail []*InkImportDetail
}

type InkImportDetail struct {
	ID             string
	InkID          string
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
}

func (p inkImportService) Find(ctx context.Context, opt *FindInkImportOpts, sort *repository.Sort, limit, offset int64) ([]*InkImportData, *repository.CountResult, error) {
	filter := &repository.SearchInkImportOpts{
		Name:   opt.Name,
		Status: opt.Status,
		ID:     opt.ID,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	inkImports, err := p.inkImportRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*InkImportData, 0)
	// write code to get ink_import_detail
	for _, inkImport := range inkImports {
		data := &InkImportData{
			ID:              inkImport.ID,
			Code:            inkImport.Code,
			Name:            inkImport.Name,
			Status:          inkImport.Status,
			Description:     inkImport.Description.String,
			Data:            inkImport.Data,
			CreatedAt:       inkImport.CreatedAt,
			UpdatedAt:       inkImport.UpdatedAt,
			CreatedBy:       inkImport.CreatedBy,
			InkImportDetail: nil,
		}
		fmt.Println("=======================>>>>>>", inkImport.ID)
		inkImportDetails, err := p.inkImportDetailRepo.Search(ctx, &repository.SearchInkImportDetailOpts{
			InkImportID: inkImport.ID,
			Limit:       10000,
		})
		if err != nil {
			return nil, nil, err
		}
		inkImportDetailResults := make([]*InkImportDetail, 0)
		for _, inkImportDetail := range inkImportDetails {
			dataDetail := &InkImportDetail{
				ID:             inkImportDetail.ID,
				InkID:          inkImportDetail.InkID,
				Name:           inkImportDetail.Name,
				Code:           inkImportDetail.Code,
				ProductCodes:   inkImportDetail.ProductCodes,
				Position:       inkImportDetail.Position,
				Location:       inkImportDetail.Location,
				Manufacturer:   inkImportDetail.Manufacturer,
				ColorDetail:    inkImportDetail.ColorDetail,
				Quantity:       inkImportDetail.Quantity,
				ExpirationDate: inkImportDetail.ExpirationDate,
				Description:    inkImportDetail.Description.String,
				Data:           inkImportDetail.Data,
			}
			inkImportDetailResults = append(inkImportDetailResults, dataDetail)
		}
		data.InkImportDetail = inkImportDetailResults
		results = append(results, data)
	}
	total, err := p.inkImportRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return results, total, nil
}

func NewService(inkImportRepo repository.InkImportRepo, inkImportDetailRepo repository.InkImportDetailRepo, inkRepo repository.InkRepo) Service {
	return &inkImportService{
		inkImportRepo:       inkImportRepo,
		inkImportDetailRepo: inkImportDetailRepo,
		inkRepo:             inkRepo,
	}

}
