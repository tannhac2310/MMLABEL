package customer

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *customerService) CreateCustomer(ctx context.Context, opt *CreateCustomerOpts) (string, error) {
	now := time.Now()

	customer := &model.Customer{
		ID:                 idutil.ULIDNow(),
		Name:               opt.Name,
		Tax:                cockroach.String(opt.Tax),
		Code:               opt.Code,
		Country:            opt.Country,
		Province:           opt.Province,
		Address:            opt.Address,
		PhoneNumber:        opt.PhoneNumber,
		Fax:                cockroach.String(opt.Fax),
		CompanyWebsite:     cockroach.String(opt.CompanyWebsite),
		CompanyPhone:       cockroach.String(opt.CompanyPhone),
		ContactPersonName:  opt.ContactPersonName,
		ContactPersonEmail: opt.ContactPersonEmail,
		ContactPersonPhone: opt.ContactPersonPhone,
		ContactPersonRole:  opt.ContactPersonRole,
		Note:               cockroach.String(opt.Note),
		Status:             opt.Status,
		CreatedBy:          opt.CreatedBy,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.customerRepo.Insert(ctx2, customer)
		if err != nil {
			return fmt.Errorf("c.customerRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}

	return customer.ID, nil
}

type CreateCustomerOpts struct {
	Name               string
	Tax                string
	Code               string
	Country            string
	Province           string
	Address            string
	PhoneNumber        string
	Fax                string
	CompanyWebsite     string
	CompanyPhone       string
	ContactPersonName  string
	ContactPersonEmail string
	ContactPersonPhone string
	ContactPersonRole  string
	Note               string
	Status             enum.CustomerStatus
	CreatedBy          string
}
