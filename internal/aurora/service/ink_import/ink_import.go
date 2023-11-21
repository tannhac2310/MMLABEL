package ink_import

import (
	"context"
	"database/sql"
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

// createInkImportDetailOpts is a struct to create inkImportDetail
type CreateInkImportDetailOpts struct {
	InkImportID string
	Quantity    float64
	ColorDetail map[string]interface{}
	Description string
	Data        map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type CreateInkImportOpts struct {
	ID              string
	Code            string
	ImportDate      time.Time
	ImportUser      string
	ImportWarehouse string
	ExportWarehouse string
	Description     sql.NullString
	Status          enum.InventoryCommonStatus
	Data            map[string]interface{}
	InkImportDetail []*CreateInkImportDetailOpts
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CreatedBy       string
}

type FindInkImportOpts struct {
	Name   string
	Code   string
	ID     string
	Status enum.InventoryCommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditInkImportOpts) error
	Create(ctx context.Context, opt *CreateInkImportOpts) (string, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, opt *FindInkImportOpts, sort *repository.Sort, limit, offset int64) ([]*repository.InkImportData, error)
}

type inkImportService struct {
	inkImportRepo       repository.InkImportRepo
	inkImportDetailRepo repository.InkImportDetailRepo
}

func (p inkImportService) Edit(ctx context.Context, opt *EditInkImportOpts) error {
	// write code to update to ink_import table and update to ink_import_detail table in transaction
	err := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// update to ink_import
		table := model.InkImport{}

		updater := cockroach.NewUpdater(table.TableName(), model.InkImportFieldID, opt.ID)

		updater.Set(model.InkImportFieldCode, opt.Code)
		updater.Set(model.InkImportFieldStatus, opt.Status)
		updater.Set(model.InkImportFieldDescription, opt.Description)
		updater.Set(model.InkImportFieldData, opt.Data)
		updater.Set(model.InkImportFieldUpdatedBy, opt.UpdatedBy)

		updater.Set(model.InkImportFieldUpdatedAt, time.Now())

		err := cockroach.UpdateFields(ctx, updater)
		if err != nil {
			return err
		}

		// write code to update to ink_import_detail table
		for _, inkImportDetail := range opt.InkImportDetail {
			// update to ink_import_detail
			table2 := model.InkImportDetail{}

			updater2 := cockroach.NewUpdater(table2.TableName(), model.InkImportDetailFieldID, inkImportDetail.ID)

			updater2.Set(model.InkImportDetailFieldInkImportID, inkImportDetail.InkImportID)
			updater2.Set(model.InkImportDetailFieldQuantity, inkImportDetail.Quantity)
			updater2.Set(model.InkImportDetailFieldColorDetail, inkImportDetail.ColorDetail)
			updater2.Set(model.InkImportDetailFieldDescription, inkImportDetail.Description)
			updater2.Set(model.InkImportDetailFieldData, inkImportDetail.Data)

			updater2.Set(model.InkImportDetailFieldUpdatedAt, time.Now())

			err2 := cockroach.UpdateFields(ctx, updater2)
			if err2 != nil {
				return err2
			}

		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p inkImportService) Create(ctx context.Context, opt *CreateInkImportOpts) (string, error) {

	// write code to insert to ink_import table and insert to ink_import_detail table in transaction
	id := idutil.ULIDNow()
	now := time.Now()
	err := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// insert to ink_import
		err1 := p.inkImportRepo.Insert(c, &model.InkImport{
			ID:              id,
			Code:            opt.Code,
			ImportDate:      opt.ImportDate,
			ImportUser:      opt.ImportUser,
			ImportWarehouse: opt.ImportWarehouse,
			ExportWarehouse: opt.ExportWarehouse,
			Description:     opt.Description,
			Status:          opt.Status,
			Data:            opt.Data,
			CreatedAt:       now,
			UpdatedAt:       now,
			CreatedBy:       opt.CreatedBy,
		})

		if err1 != nil {
			return err1
		}
		// write code to insert to ink_import_detail table
		for _, inkImportDetail := range opt.InkImportDetail {
			err2 := p.inkImportDetailRepo.Insert(c, &model.InkImportDetail{
				ID:          idutil.ULIDNow(),
				Code:        inkImportDetail.Code,
				InkImportID: inkImportDetail.InkImportID,
				Quantity:    inkImportDetail.Quantity,
				ColorDetail: inkImportDetail.ColorDetail,
				Description: cockroach.String(inkImportDetail.Description),
				Data:        inkImportDetail.Data,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			})
			if err2 != nil {
				return err2
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}
	return id, nil
}

func (p inkImportService) Delete(ctx context.Context, id string) error {
	return p.inkImportRepo.SoftDelete(ctx, id)
}

func (p inkImportService) Find(ctx context.Context, opt *FindInkImportOpts, sort *repository.Sort, limit, offset int64) ([]*repository.InkImportData, error) {
	ink_imports, err := p.inkImportRepo.Search(ctx, &repository.SearchInkImportOpts{
		Name:   opt.Name,
		Code:   opt.Code,
		Status: opt.Status,
		ID:     opt.ID,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	})
	if err != nil {
		return nil, err
	}
	return ink_imports, nil
}

func NewService(inkImportRepo repository.InkImportRepo, inkImportDetailRepo repository.InkImportDetailRepo) Service {
	return &inkImportService{
		inkImportRepo:       inkImportRepo,
		inkImportDetailRepo: inkImportDetailRepo,
	}

}
