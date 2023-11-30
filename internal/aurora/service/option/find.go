package option

import (
	"context"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type EditOptionOpts struct {
	ID     string
	Name   string
	Code   string
	Data   map[string]interface{}
	Status enum.CommonStatus
}

type CreateOptionOpts struct {
	Name      string
	Code      string
	Entity    string
	Data      map[string]interface{}
	Status    enum.CommonStatus
	CreatedBy string
}

type FindOptionOpts struct {
	Name   string
	Entity string
	Status enum.CommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditOptionOpts) error
	Create(ctx context.Context, opt *CreateOptionOpts) (string, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, opt *FindOptionOpts, sort *repository.Sort, limit, offset int64) ([]*OptionData, *repository.CountResult, error)
}
type optionService struct {
	optionRepo repository.OptionRepo
}

type OptionData struct {
	*repository.OptionData
}

func (p optionService) Edit(ctx context.Context, opt *EditOptionOpts) error {
	table := model.Ink{}
	// find ink by id
	updater := cockroach.NewUpdater(table.TableName(), model.OptionFieldID, opt.ID)

	updater.Set(model.OptionFieldName, opt.Name)
	updater.Set(model.OptionFieldCode, opt.Code)
	updater.Set(model.OptionFieldStatus, opt.Status)
	updater.Set(model.OptionFieldData, opt.Data)

	updater.Set(model.InkFieldUpdatedAt, time.Now())

	err := cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return err
	}
	return nil
}

func (p optionService) Create(ctx context.Context, opt *CreateOptionOpts) (string, error) {
	id := idutil.ULIDNow()
	err := p.optionRepo.Insert(ctx, &model.Option{
		ID:        id,
		Name:      opt.Name,
		Entity:    opt.Entity,
		Code:      opt.Code,
		Data:      opt.Data,
		Status:    opt.Status,
		CreatedBy: opt.CreatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p optionService) Delete(ctx context.Context, id string) error {
	return p.optionRepo.SoftDelete(ctx, id)
}

func (c *optionService) Find(ctx context.Context, opts *FindOptionOpts, sort *repository.Sort, limit, offset int64) ([]*OptionData, *repository.CountResult, error) {
	filter := &repository.SearchOptionsOpts{
		Name:   opts.Name,
		Entity: opts.Entity,
		Status: opts.Status,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
	options, err := c.optionRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	total, err := c.optionRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	results := make([]*OptionData, 0, len(options))
	for _, option := range options {
		if err != nil {
			return nil, nil, err
		}
		results = append(results, &OptionData{
			OptionData: option,
		})
	}
	return results, total, nil
}

func NewService(
	optionRepo repository.OptionRepo,
) Service {
	return &optionService{
		optionRepo: optionRepo,
	}

}
