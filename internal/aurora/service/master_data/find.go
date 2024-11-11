package master_data

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

func (s *masterDataService) FindMasterData(ctx context.Context, opt *FindMasterDataOpts) ([]*MasterData, int64, error) {
	filter := repository.SearchMasterDataOpts{
		Type:   opt.Type,
		Limit:  opt.Limit,
		Search: opt.Search,
		Offset: opt.Offset,
		Sort:   nil,
	}
	// 1. find master data
	masterData, err := s.masterDataRep.Search(ctx, &filter)
	if err != nil {
		return nil, 0, fmt.Errorf("s.masterDataRep.Search: %w", err)
	}
	// 2. count master data
	count, err := s.masterDataRep.Count(ctx, &filter)
	if err != nil {
		return nil, 0, fmt.Errorf("s.masterDataRep.Count: %w", err)
	}
	// 3. get user fields
	masterDataIDs := make([]string, 0)
	for _, md := range masterData {
		masterDataIDs = append(masterDataIDs, md.ID)
	}
	userFields, err := s.masterDataUserField.Search(ctx, &repository.SearchMasterDataUserFieldOpts{
		MasterDataIDs: masterDataIDs,
		Offset:        0,
		Limit:         10000, // getall
	})
	if err != nil {
		return nil, 0, fmt.Errorf("s.masterDataUserField.Search: %w", err)
	}
	// 4. map user fields
	userFieldsMap := make(map[string][]*repository.MasterDataUserFieldData)
	for _, uf := range userFields {
		if _, ok := userFieldsMap[uf.MasterDataID]; !ok {
			userFieldsMap[uf.MasterDataID] = make([]*repository.MasterDataUserFieldData, 0)
		}
		userFieldsMap[uf.MasterDataID] = append(userFieldsMap[uf.MasterDataID], uf)
	}
	result := make([]*MasterData, 0)
	// 5. map user fields to master data
	for _, md := range masterData {
		uf, ok := userFieldsMap[md.ID]
		if !ok {
			//return nil, 0, fmt.Errorf("user fields not found for master data id: %s", md.ID)
		}

		ufData := make([]*MasterDataUserField, 0)
		for _, f := range uf {
			ufData = append(ufData, &MasterDataUserField{
				ID:           f.ID,
				MasterDataID: f.MasterDataID,
				FieldName:    f.FieldName,
				FieldValue:   f.FieldValue,
			})
		}
		result = append(result, &MasterData{
			ID:          md.ID,
			Type:        md.Type,
			Name:        md.Name,
			Code:        md.Code,
			Status:      md.Status,
			Description: md.Description,
			CreatedAt:   md.CreatedAt,
			CreatedBy:   md.CreatedBy,
			UpdatedAt:   md.UpdatedAt,
			UpdatedBy:   md.UpdatedBy,
			UserFields:  ufData,
		})
	}
	return result, count.Count, nil
}
