package office

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"
)

func (b officeService) CreateOffice(ctx context.Context, opt *CreateOfficeOpts) (string, error) {
	now := time.Now()

	office := &model.Office{
		ID:          idutil.ULIDNow(),
		Name:        opt.Name,
		Phone:       opt.Phone,
		Address:     opt.Address,
		ProvinceID:  opt.ProvinceID,
		DistrictID:  opt.DistrictID,
		Status:      opt.Status,
		PhotoURL:    opt.PhotoURL,
		Description: opt.Description,
		CreatedBy:   opt.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := b.officeRepo.Insert(ctx, office)
	if err != nil {
		return "", fmt.Errorf("p.officeRepo.Insert: %w", err)
	}

	return office.ID, nil
}

type CreateOfficeOpts struct {
	Name        string
	Phone       string
	Address     string
	ProvinceID  int64
	DistrictID  int64
	Status      enum.CommonStatus
	PhotoURL    string
	Description string
	CreatedBy   string
}
