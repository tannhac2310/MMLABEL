package ink

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type InkFormulation struct {
	ID          string
	InkID       string
	Quantity    float64
	Description string
	InkCode     string
	InkName     string
}
type CreateFormulation struct {
	InkID       string
	Quantity    float64
	Description string
}
type CreateInkMixingOpts struct {
	Name           string
	Code           string
	ProductCodes   []string
	Manufacturer   string
	Quantity       float64
	ExpirationDate string
	Position       string
	Location       string
	Description    string
	InkFormula     []CreateFormulation
	Status         enum.CommonStatus
	CreatedBy      string
}

type InkMixingData struct {
	ID             string
	Name           string
	Code           string
	InkID          string
	ProductCodes   []string
	Quantity       float64
	ExpirationDate string
	Position       string
	Location       string
	Description    string
	InkFormula     []InkFormulation
	Status         enum.CommonStatus
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedByName  string
	UpdatedByName  string
}

func (p inkService) MixInk(ctx context.Context, opt *CreateInkMixingOpts) (string, error) {
	// Create ink mixing
	now := time.Now()
	//nowDate := now.Format("2006-01-02")
	newInkID := idutil.ULIDNow()
	inkMixingID := idutil.ULIDNow()
	errTx := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// create new ink
		err := p.inkRepo.Insert(c, &model.Ink{
			ID:             newInkID,
			Name:           opt.Name,
			Manufacturer:   opt.Manufacturer,
			Code:           opt.Code,
			MixingID:       cockroach.String(inkMixingID),
			ProductCodes:   opt.ProductCodes,
			Position:       opt.Position,
			Location:       opt.Location,
			Quantity:       opt.Quantity,
			ExpirationDate: opt.ExpirationDate,
			Description:    cockroach.String(opt.Description),
			Status:         opt.Status,
			CreatedBy:      opt.CreatedBy,
			UpdatedBy:      opt.CreatedBy,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		if err != nil {
			return fmt.Errorf("error creating ink: %w", err)
		}

		err = p.inkMixingRepo.Insert(c, &model.InkMixing{
			ID:    inkMixingID,
			Name:  opt.Name,
			Code:  opt.Code,
			InkID: newInkID,
			//MixingDate:  nowDate,
			Description: opt.Description,
			Status:      enum.CommonStatusActive,
			CreatedBy:   opt.CreatedBy,
			UpdatedBy:   opt.CreatedBy,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		if err != nil {
			return fmt.Errorf("error creating ink mixing: %w", err)
		}

		// create ink mixing detail
		for _, ink := range opt.InkFormula {
			// minus ink quantity
			inkData, err := p.inkRepo.FindByID(c, ink.InkID)
			if err != nil {
				return fmt.Errorf("Không tìm thấy mực: %w", err)
			}
			fmt.Println("inkData: ", inkData.Ink.Name, inkData.Ink.Quantity, ink.Quantity)
			newValue := inkData.Ink.Quantity - ink.Quantity
			if newValue < 0 {
				return fmt.Errorf("số lượng mực không đủ, inkName: %s, Tồn kho: %v, Cần dùng: %v", inkData.Ink.Name, inkData.Ink.Quantity, ink.Quantity)
			}
			inkData.Ink.Quantity = newValue
			err = p.inkRepo.Update(c, inkData.Ink)
			if err != nil {
				return fmt.Errorf("error updating ink: %w", err)
			}

			// create ink mixing detail
			err = p.inkMixingDetailRepo.Insert(c, &model.InkMixingDetail{
				ID:          idutil.ULIDNow(),
				InkMixingID: inkMixingID,
				InkID:       ink.InkID,
				Quantity:    ink.Quantity,
				Description: ink.Description,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err != nil {
				return fmt.Errorf("error creating ink mixing detail: %w", err)
			}
		}

		return nil
	})
	if errTx != nil {
		return "", fmt.Errorf("Pha mực lỗi, %w", errTx)
	}
	return inkMixingID, nil
}

// find ink mixing
type FindInkMixingOpts struct {
	IDs []string
	//Search string
	//InkID  string
	//InkIDs []string
	Limit  int64
	Offset int64
}

func (p inkService) FindInkMixing(ctx context.Context, opt *FindInkMixingOpts) ([]*InkMixingData, *repository.CountResult, error) {
	filter := &repository.SearchInkMixingOpts{
		IDs: opt.IDs,
		//Search: opt.Search,
		//InkID:  opt.InkID,
		//InkIDs: opt.InkIDs,
		Limit:  opt.Limit,
		Offset: opt.Offset,
	}
	inkMixing, err := p.inkMixingRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching ink mixing: %w", err)
	}
	count, err := p.inkMixingRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("error counting ink mixing: %w", err)
	}
	inkIDs := make([]string, 0)
	for _, ink := range inkMixing {
		inkIDs = append(inkIDs, ink.InkID)
	}
	inkData, err := p.inkRepo.Search(ctx, &repository.SearchInkOpts{
		IDs:    inkIDs,
		Limit:  int64(len(inkIDs)),
		Offset: 0,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error searching ink: %w", err)
	}
	inkDataMap := make(map[string]*repository.InkData)
	for _, ink := range inkData {
		inkDataMap[ink.ID] = ink
	}
	// ink mixing detail
	inkMixingDetail, err := p.inkMixingDetailRepo.Search(ctx, &repository.SearchInkMixingDetailOpts{
		InkMixingIDs: opt.IDs,
		Limit:        int64(len(inkMixing) * 100),
		Offset:       0,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error searching ink mixing detail: %w", err)
	}
	inkMixingDetailMap := make(map[string][]*repository.InkMixingDetailData)
	for _, detail := range inkMixingDetail {
		if _, ok := inkMixingDetailMap[detail.InkMixingID]; !ok {
			inkMixingDetailMap[detail.InkMixingID] = make([]*repository.InkMixingDetailData, 0)
		}
		inkMixingDetailMap[detail.InkMixingID] = append(inkMixingDetailMap[detail.InkMixingID], detail)
	}
	results := make([]*InkMixingData, 0)
	for _, im := range inkMixing {
		_inkDetail, ok := inkDataMap[im.InkID]
		if !ok {
			continue
		}
		inkFormula := make([]InkFormulation, 0)
		if details, ok := inkMixingDetailMap[im.ID]; ok {
			for _, detail := range details {
				inkFormula = append(inkFormula, InkFormulation{
					ID:          detail.ID,
					InkID:       detail.InkID,
					Quantity:    detail.Quantity,
					Description: detail.Description,
					InkName:     detail.InkName,
					InkCode:     detail.InkCode,
				})
			}
		}

		results = append(results, &InkMixingData{
			ID:             im.ID,
			Name:           im.Name,
			Code:           im.Code,
			InkID:          im.InkID,
			ProductCodes:   _inkDetail.ProductCodes,
			Quantity:       _inkDetail.Quantity,
			ExpirationDate: _inkDetail.ExpirationDate,
			Position:       _inkDetail.Position,
			Location:       _inkDetail.Location,
			Description:    im.Description,
			InkFormula:     inkFormula,
			Status:         im.Status,
			CreatedBy:      im.CreatedBy,
			CreatedAt:      im.CreatedAt,
			UpdatedAt:      im.UpdatedAt,
			CreatedByName:  im.CreatedByName,
			UpdatedByName:  im.UpdatedByName,
		})
	}
	return results, count, nil
}
