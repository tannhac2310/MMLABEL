package customer

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *customerService) EditCustomer(ctx context.Context, opt *EditCustomerOpts) error {
	customer := model.Customer{
		ID:                 opt.ID,
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
		UpdatedAt:          time.Now(),
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.customerRepo.Update(ctx2, &customer)
		if err != nil {
			return fmt.Errorf("c.customerRepo.Update: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return errTx
	}

	return nil
}

type EditCustomerOpts struct {
	ID                 string
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
	Status             int16
	CreatedBy          string
}
