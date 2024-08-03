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

type EditInkFormulation struct {
	ID          string
	InkID       string
	Quantity    float64
	Description string
}
type EditInkMixingOpts struct {
	ID             string
	Name           string
	Code           string
	ProductCodes   []string
	Quantity       float64
	Manufacturer   string
	ExpirationDate string
	ColorDetail    map[string]interface{}
	Position       string
	Location       string
	Description    string
	Data           map[string]interface{}
	InkFormula     []EditInkFormulation
	Status         enum.CommonStatus
	UpdatedBy      string
}

func (p inkService) EditInkMixing(ctx context.Context, opt *EditInkMixingOpts) error {
	// Create ink mixing
	now := time.Now()
	errTx := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// find ink mixing
		inkMixing, err := p.inkMixingRepo.FindByID(ctx, opt.ID)
		if err != nil {
			return fmt.Errorf("find ink mixing, %w", err)
		}
		inkMixing.Code = opt.Code
		inkMixing.Name = opt.Name
		inkMixing.Description = opt.Description
		inkMixing.UpdatedBy = opt.UpdatedBy
		inkMixing.UpdatedAt = now
		// update ink mixing
		err = p.inkMixingRepo.Update(ctx, inkMixing)
		if err != nil {
			return fmt.Errorf("update ink mixing, %w", err)
		}

		// update ink
		err = p.Edit(ctx, &EditInkOpts{
			ID:             inkMixing.InkID,
			Name:           opt.Name,
			Code:           opt.Code,
			ProductCodes:   opt.ProductCodes,
			Position:       opt.Position,
			Location:       opt.Location,
			Quantity:       opt.Quantity,
			Manufacturer:   opt.Manufacturer,
			ColorDetail:    nil,
			ExpirationDate: opt.ExpirationDate,
			Description:    opt.Description,
			Data:           nil,
			Status:         opt.Status,
			UpdatedBy:      opt.UpdatedBy,
		})
		if err != nil {
			return fmt.Errorf("update ink, %w", err)
		}

		oldInkMixingDetail, err := p.inkMixingDetailRepo.Search(ctx, &repository.SearchInkMixingDetailOpts{
			InkMixingID: opt.ID,
			Offset:      0,
			Limit:       1000,
		})
		if err != nil {
			return fmt.Errorf("search ink mixing detail, %w", err)
		}
		oldInkMixingDetailMap := map[string]repository.InkMixingDetailData{}
		// delete ink mixing detail
		for _, ink := range oldInkMixingDetail {
			oldInkMixingDetailMap[ink.InkMixingDetail.InkID] = *ink // save old ink mixing detail
			err = p.inkMixingDetailRepo.SoftDelete(ctx, ink.InkMixingDetail.ID)
			if err != nil {
				return fmt.Errorf("delete ink mixing detail, %w", err)
			}
		}

		for _, ink := range opt.InkFormula {
			// create ink mixing detail because old ink mixing detail has been deleted
			err = p.inkMixingDetailRepo.Insert(ctx, &model.InkMixingDetail{
				ID:          idutil.ULIDNow(),
				InkMixingID: opt.ID,
				InkID:       ink.InkID,
				Quantity:    ink.Quantity,
				Description: ink.Description,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err != nil {
				return fmt.Errorf("error creating ink mixing detail: %w", err)
			}

			// update ink quantity
			inkData, err := p.inkRepo.FindByID(ctx, ink.InkID)
			if err != nil {
				return fmt.Errorf("find ink, %w", err)
			}
			inkData.Ink.UpdatedAt = now
			inkData.Ink.Quantity -= ink.Quantity
			oldValue, ok := oldInkMixingDetailMap[ink.InkID]
			if ok {
				inkData.Ink.Quantity += oldValue.InkMixingDetail.Quantity
				delete(oldInkMixingDetailMap, ink.InkID)
			}
			err = p.inkRepo.Update(ctx, inkData.Ink)
			if err != nil {
				return fmt.Errorf("update ink quanlity, %w", err)
			}
		}

		// update ink quantity
		for _, ink := range oldInkMixingDetailMap {
			inkData, err := p.inkRepo.FindByID(ctx, ink.InkMixingDetail.InkID)
			if err != nil {
				return fmt.Errorf("find ink, %w", err)
			}
			inkData.Ink.UpdatedAt = now
			inkData.Ink.Quantity += ink.InkMixingDetail.Quantity
			err = p.inkRepo.Update(ctx, inkData.Ink)
			if err != nil {
				return fmt.Errorf("update ink quanlity, %w", err)
			}
		}

		return nil
	})
	if errTx != nil {
		return fmt.Errorf("edit ink mixing, %w", errTx)
	}
	return nil
}
