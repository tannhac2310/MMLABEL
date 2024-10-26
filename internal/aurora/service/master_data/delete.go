package master_data

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (s *masterDataService) DeleteMasterData(ctx context.Context, opt *DeleteMasterDataOpts) error {
	if opt == nil || opt.ID == "" {
		return fmt.Errorf("opt is required")
	}
	errTx := cockroach.ExecInTx(ctx, func(ctx context.Context) error {
		// 1. Delete master data
		if err := s.masterDataRep.SoftDelete(ctx, opt.ID); err != nil {
			return fmt.Errorf("s.masterDataRep.SoftDelete: %w", err)
		}

		// 2. Delete user fields
		if err := s.masterDataUserField.DeleteByMasterDataIDs(ctx, []string{opt.ID}); err != nil {
			return fmt.Errorf("s.masterDataUserField.DeleteByMasterDataIDs: %w", err)
		}
		return nil
	})
	if errTx != nil {
		return fmt.Errorf("delete master data: %w", errTx)
	}
	return nil
}
