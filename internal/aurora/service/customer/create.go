package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *customerService) CreateCustomer(ctx context.Context, opt *CreateCustomerOpts) (string, error) {
	now := time.Now()
	b, _ := json.Marshal(opt)
	searchContent := fmt.Sprintf("%s %s %s", opt.Name, opt.CompanyEmail, string(b))
	customer := &model.Customer{
		ID:                 idutil.ULIDNow(),
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
		Note:               cockroach.String(opt.Note),
		SearchContent:      cockroach.String(searchContent),
		Data:               opt.Data,
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
