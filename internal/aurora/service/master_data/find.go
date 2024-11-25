package master_data

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

func (s *masterDataService) FindMasterData(ctx context.Context, opt *FindMasterDataOpts) ([]*MasterData, int64, error) {
	filter := repository.SearchMasterDataOpts{
		Type:   opt.Type,
		IDs:    opt.IDs,
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
	// extra
	ufFieldValues, err := s.customFieldRepo.Search(ctx, &repository.SearchCustomFieldsOpts{
		EntityType: enum.CustomFieldTypeProductionPlan,
		Values:     masterDataIDs,
		Limit:      10000, // getall
		Offset:     0,
		Sort:       nil,
	})

	if err != nil {
		return nil, 0, fmt.Errorf("s.customFieldRepo.Search: %w", err)
	}
	// ufFieldValuesMap
	baninnguonIDMapper := make(map[string][]string)
	for _, uf := range ufFieldValues {
		if _, ok := baninnguonIDMapper[uf.Value]; !ok {
			baninnguonIDMapper[uf.Value] = make([]string, 0)
		}
		baninnguonIDMapper[uf.Value] = append(baninnguonIDMapper[uf.Value], uf.EntityID)
	}
	// find baninnguonIDs from production_order_device_config.go
	deviceConfigItems, err := s.productionOrderDeviceConfigRepo.Search(ctx, &repository.SearchProductionOrderDeviceConfigOpts{
		MasterDataIDS: masterDataIDs,
		Limit:         10000,
		Offset:        0,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("tìm bản in nguồn từ cấu hình thiết bị %w", err)
	}

	for _, item := range deviceConfigItems {
		if !item.ProductionPlanID.Valid { // TODO add && production order id
			continue
		}
		if item.MaKhung.Valid {
			if _, ok := baninnguonIDMapper[item.MaKhung.String]; !ok {
				baninnguonIDMapper[item.MaKhung.String] = make([]string, 0)
			}
			baninnguonIDMapper[item.MaKhung.String] = append(baninnguonIDMapper[item.MaKhung.String], item.ProductionPlanID.String)
		}

		if item.MaPhim.Valid {
			if _, ok := baninnguonIDMapper[item.MaPhim.String]; !ok {
				baninnguonIDMapper[item.MaPhim.String] = make([]string, 0)
			}
			baninnguonIDMapper[item.MaPhim.String] = append(baninnguonIDMapper[item.MaPhim.String], item.ProductionPlanID.String)
		}
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

		// extra
		binnguonIDs, ok := baninnguonIDMapper[md.ID]
		if !ok {
			//return nil, 0, fmt.Errorf("user fields not found for master data id: %s", md.ID)
		}

		result = append(result, &MasterData{
			ID:                md.ID,
			Type:              md.Type,
			Name:              md.Name,
			Code:              md.Code,
			Status:            md.Status,
			Description:       md.Description,
			CreatedAt:         md.CreatedAt,
			CreatedBy:         md.CreatedBy,
			UpdatedAt:         md.UpdatedAt,
			UpdatedBy:         md.UpdatedBy,
			ProductionPlanIDs: binnguonIDs,
			UserFields:        ufData,
		})
	}
	return result, count.Count, nil
}
