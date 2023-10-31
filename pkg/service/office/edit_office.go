package office

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/model"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (b *officeService) EditOffice(ctx context.Context, opt *EditOfficeOpts) error {
	var err error
	table := model.Office{}
	updater := cockroach.NewUpdater(table.TableName(), model.OfficeFieldID, opt.ID)

	updater.Set(model.OfficeFieldName, opt.Name)
	updater.Set(model.OfficeFieldPhone, opt.Phone)
	updater.Set(model.OfficeFieldAddress, opt.Address)
	updater.Set(model.OfficeFieldProvinceID, opt.ProvinceID)
	updater.Set(model.OfficeFieldDistrictID, opt.DistrictID)
	updater.Set(model.OfficeFieldStatus, opt.Status)
	updater.Set(model.OfficeFieldPhotoURL, opt.PhotoURL)
	updater.Set(model.OfficeFieldDescription, opt.Description)

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("err update office: %w", err)
	}

	return nil
}

type EditOfficeOpts struct {
	ID          string
	Name        string
	Phone       string
	Address     string
	ProvinceID  int64
	DistrictID  int64
	Status      enum.CommonStatus
	PhotoURL    string
	Description string
}
