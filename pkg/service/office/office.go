package office

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

type Service interface {
	CreateOffice(ctx context.Context, opt *CreateOfficeOpts) (string, error)
	EditOffice(ctx context.Context, opt *EditOfficeOpts) error
	FindOffices(ctx context.Context, opts *FindOfficesOpts, limit, offset int64) ([]*repository.OfficeData, *repository.CountResult, error)
	SoftDelete(ctx context.Context, id string) error
	FindOfficeByID(ctx context.Context, id string) (*repository.OfficeData, error)
}

type officeService struct {
	officeRepo repository.OfficeRepo
}

func NewService(
	officeRepo repository.OfficeRepo,
) Service {
	return &officeService{
		officeRepo: officeRepo,
	}
}
