package customer

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *customerService) EditCustomer(ctx context.Context, opt *EditCustomerOpts) error {
	var err error
	table := model.Customer{}
	updater := cockroach.NewUpdater(table.TableName(), model.CustomerFieldID, opt.ID)

	updater.Set(model.CustomerFieldName, opt.Name)
	updater.Set(model.CustomerFieldAvatar, opt.Avatar)
	updater.Set(model.CustomerFieldPhoneNumber, opt.PhoneNumber)
	updater.Set(model.CustomerFieldEmail, opt.Email)
	updater.Set(model.CustomerFieldStatus, opt.Status)
	updater.Set(model.CustomerFieldType, opt.Type)
	updater.Set(model.CustomerFieldAddress, opt.Address)

	updater.Set(model.CustomerFieldUpdatedAt, time.Now())

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update customer failed %w", err)
	}
	return nil
}

type EditCustomerOpts struct {
	ID          string
	Name        string
	Avatar      string
	PhoneNumber string
	Email       string
	Status      int16
	Type        int16
	Address     string
}
