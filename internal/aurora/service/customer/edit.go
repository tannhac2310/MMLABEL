package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (c *customerService) EditCustomer(ctx context.Context, opt *EditCustomerOpts) error {
	b, _ := json.Marshal(opt)
	searchContent := fmt.Sprintf("%s %s %s", opt.Name, opt.CompanyEmail, string(b))

	customer := model.Customer{
		ID:                 opt.ID,
		Name:               opt.Name,
		Tax:                cockroach.String(opt.Tax),
		Code:               opt.Code,
		Country:            opt.Country,
		Province:           opt.Province,
		Address:            opt.Address,
		Fax:                cockroach.String(opt.Fax),
		CompanyWebsite:     cockroach.String(opt.CompanyWebsite),
		CompanyPhone:       cockroach.String(opt.CompanyPhone),
		CompanyEmail:       cockroach.String(opt.CompanyEmail),
		ContactPersonName:  cockroach.String(opt.ContactPersonName),
		ContactPersonEmail: cockroach.String(opt.ContactPersonEmail),
		ContactPersonPhone: cockroach.String(opt.ContactPersonPhone),
		ContactPersonRole:  cockroach.String(opt.ContactPersonRole),
		Data:               opt.Data,
		SearchContent:      cockroach.String(searchContent),
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
	Fax                string
	CompanyWebsite     string
	CompanyPhone       string
	CompanyEmail       string
	ContactPersonName  string
	ContactPersonEmail string
	ContactPersonPhone string
	ContactPersonRole  string
	Note               string
	Status             enum.CustomerStatus
	Data               any
	CreatedBy          string
}
