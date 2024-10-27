package master_data

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (s *masterDataService) CreateMasterData(ctx context.Context, opt *CreateMasterDataOpts) (string, error) {
	if opt == nil {
		return "", fmt.Errorf("opt is required")
	}

	masterDataID := idutil.ULIDNow()
	now := time.Now()
	errTx := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		// 0. count master data
		count, err := s.masterDataRep.Count(ctx, &repository.SearchMasterDataOpts{
			Type:         opt.Type,
			IsIncludeDel: true,
		})
		fmt.Println("count", count)
		if err != nil {
			return fmt.Errorf("s.masterDataRep.Count: %w", err)
		}
		masterDataID = fmt.Sprintf("%s-%d", opt.Type, count.Count+1)
		// 1. Insert master data
		masterData := &model.MasterData{
			ID:          masterDataID,
			Type:        opt.Type,
			Name:        opt.Name,
			Description: opt.Description,
			Status:      opt.Status,
			CreatedBy:   opt.CreatedBy,
			CreatedAt:   now,
			UpdatedAt:   now,
			UpdatedBy:   opt.CreatedBy,
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
