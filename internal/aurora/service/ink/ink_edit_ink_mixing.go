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
		inkMixing.InkMixing.Description = opt.Description
		inkMixing.InkMixing.UpdatedBy = opt.UpdatedBy
		inkMixing.InkMixing.UpdatedAt = now
		// update ink mixing
		err = p.inkMixingRepo.Update(ctx, inkMixing.InkMixing)
		if err != nil {
			return fmt.Errorf("update ink mixing, %w", err)
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
			oldInkMixingDetailMap[ink.InkMixingDetail.ID] = *ink
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
			}
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
