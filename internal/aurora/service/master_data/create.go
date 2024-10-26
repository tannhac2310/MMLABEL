package master_data

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (s *masterDataService) CreateMasterData(ctx context.Context, opt *CreateMasterDataOpts) (string, error) {
	if opt == nil {
		return "", fmt.Errorf("opt is required")
	}

	masterDataID := idutil.ULIDNow()
	errTx := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		// 1. Insert master data
		masterData := &model.MasterData{
			ID:          masterDataID,
			Type:        opt.Type,
			Name:        opt.Name,
			Description: opt.Description,
			CreatedBy:   opt.CreatedBy,
		}
		if err := s.masterDataRep.Insert(ctx, masterData); err != nil {
			return fmt.Errorf("s.masterDataRep.Insert: %w", err)
		}

		// 2. Insert user fields
		for _, uf := range opt.UserFields {
			masterDataUserField := &model.MasterDataUserField{
				ID:           idutil.ULIDNow(),
				MasterDataID: masterDataID,
				FieldName:    uf.FieldName,
				FieldValue:   uf.FieldValue,
			}
			if err := s.masterDataUserField.Insert(ctx, masterDataUserField); err != nil {
				return fmt.Errorf("s.masterDataUserField.Insert: %w", err)
			}
		}
		return nil
	})
	if errTx != nil {
		return "", fmt.Errorf("insert master data: %w", errTx)
	}
	return masterDataID, nil
}
