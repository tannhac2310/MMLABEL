package master_data

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (s *masterDataService) UpdateMasterData(ctx context.Context, opt *UpdateMasterDataOpts) error {
	if opt == nil || opt.ID == "" {
		return fmt.Errorf("opt is required")
	}

	errTx := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		// 1. Update master data
		table := model.MasterData{}
		updater := cockroach.NewUpdater(table.TableName(), model.MasterDataFieldID, opt.ID)
		updater.Set(model.MasterDataFieldName, opt.Name)
		updater.Set(model.MasterDataFieldDescription, opt.Description)
		updater.Set(model.MasterDataFieldUpdatedBy, opt.UpdateBy)
		updater.Set(model.MasterDataFieldUpdatedAt, time.Now())

		if err := cockroach.UpdateFields(ctx, updater); err != nil {
			return fmt.Errorf("s.masterDataRep.Update: %w", err)
		}

		// 2. Update user fields
		if err := s.masterDataUserField.DeleteByMasterDataIDs(ctx, []string{opt.ID}); err != nil {
			return fmt.Errorf("s.masterDataUserField.DeleteByMasterDataIDs: %w", err)
		}

		for _, userField := range opt.UserFields {
			modelData := &model.MasterDataUserField{
				ID:           idutil.ULIDNow(),
				MasterDataID: opt.ID,
				FieldName:    userField.FieldName,
				FieldValue:   userField.FieldValue,
				CreatedBy:    opt.UpdateBy,
				CreatedAt:    time.Now(),
			}
			if err := s.masterDataUserField.Insert(ctx, modelData); err != nil {
				return fmt.Errorf("s.masterDataUserField.Insert: %w", err)
			}
		}

		return nil
	})

	if errTx != nil {
		return fmt.Errorf("update master data: %w", errTx)
	}

	return nil
}
