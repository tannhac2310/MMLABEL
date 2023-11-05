package customer

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *customerService) CreateCustomer(ctx context.Context, opt *CreateCustomerOpts) (string, error) {
	now := time.Now()

	customer := &model.Customer{
		ID:          idutil.ULIDNow(),
		Name:        opt.Name,
		Avatar:      cockroach.String(opt.Avatar),
		PhoneNumber: cockroach.String(opt.PhoneNumber),
		Email:       cockroach.String(opt.Email),
		Status:      opt.Status,
		Type:        opt.Type,
		Address:     cockroach.String(opt.Address),
		CreatedBy:   opt.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
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
	Name        string
	Avatar      string
	PhoneNumber string
	Email       string
	Status      int16
	Type        int16
	Address     string
	CreatedBy   string
}
