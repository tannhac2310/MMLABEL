package ink

import (
	"context"
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
	Manufacturer   string
	ColorDetail    map[string]interface{}
	Quantity       float64
	ExpirationDate time.Time
	Description    string
	Data           map[string]interface{}
	Status         enum.InventoryCommonStatus
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
	ExpirationDate time.Time
	Description    string
	Data           map[string]interface{}
	Status         enum.InventoryCommonStatus
	CreatedBy      string
}

type FindInkOpts struct {
	Name   string
	Code   string
	ID     string
	Status enum.InventoryCommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditInkOpts) error
	Create(ctx context.Context, opt *CreateInkOpts) (string, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, opt *FindInkOpts, sort *repository.Sort, limit, offset int64) ([]*repository.InkData, error)
}

type inkService struct {
	inkRepo repository.InkRepo
}

func (p inkService) Edit(ctx context.Context, opt *EditInkOpts) error {
	table := model.Ink{}

	updater := cockroach.NewUpdater(table.TableName(), model.InkFieldID, opt.ID)

	updater.Set(model.InkFieldCode, opt.Code)
	updater.Set(model.InkFieldProductCodes, opt.ProductCodes)
	updater.Set(model.InkFieldStatus, opt.Status)
	updater.Set(model.InkFieldColorDetail, opt.ColorDetail)
	updater.Set(model.InkFieldDescription, opt.Description)
	updater.Set(model.InkFieldData, opt.Data)
	updater.Set(model.InkFieldUpdatedBy, opt.UpdatedBy)

	updater.Set(model.InkFieldUpdatedAt, time.Now())

	err := cockroach.UpdateFields(ctx, updater)
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

func (p inkService) Find(ctx context.Context, opt *FindInkOpts, sort *repository.Sort, limit, offset int64) ([]*repository.InkData, error) {
	inks, err := p.inkRepo.Search(ctx, &repository.SearchInkOpts{
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
	return inks, nil
}

func NewService(inkRepo repository.InkRepo) Service {
	return &inkService{
		inkRepo: inkRepo,
	}

}
