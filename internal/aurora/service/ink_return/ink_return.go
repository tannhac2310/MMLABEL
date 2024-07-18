package ink_return

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/generic"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type EditInkReturnDetailOpts struct {
	InkID             string
	InkExportDetailID string
	Quantity          float64
	ColorDetail       map[string]interface{}
	Description       string
	Data              map[string]interface{}
}

type EditInkReturnOpts struct {
	ID              string
	UpdatedBy       string
	Description     string
	InkReturnDetail []*EditInkReturnDetailOpts
}

type CreateInkReturnDetailOpts struct {
	InkID             string
	InkExportDetailID string
	Quantity          float64
	ColorDetail       map[string]interface{}
	Description       string
	Data              map[string]interface{}
}
type CreateInkReturnOpts struct {
	Name            string
	Code            string
	InkExportID     string
	Description     string
	Data            map[string]interface{}
	InkReturnDetail []*CreateInkReturnDetailOpts
	CreatedBy       string
}
type FindInkReturnOpts struct {
	Name         string
	ID           string
	InkExportID  string
	InkExportIDs []string
	Status       enum.InventoryCommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditInkReturnOpts) error
	Create(ctx context.Context, opt *CreateInkReturnOpts) (string, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, opt *FindInkReturnOpts, sort *repository.Sort, limit, offset int64) ([]*InkReturnData, *repository.CountResult, error)
}

type inkReturnService struct {
	inkReturnRepo       repository.InkReturnRepo
	inkReturnDetailRepo repository.InkReturnDetailRepo
	inkExportDetailRepo repository.InkExportDetailRepo
	inkRepo             repository.InkRepo
}

func (p inkReturnService) Edit(ctx context.Context, opt *EditInkReturnOpts) error {
	now := time.Now()

	inkReturn, err := p.inkReturnRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return err
	}
	inkReturn.Description = cockroach.String(opt.Description)
	inkReturn.UpdatedBy = opt.UpdatedBy
	inkReturn.UpdatedAt = now

	inkReturnItems, err := p.inkReturnDetailRepo.Search(ctx, &repository.SearchInkReturnDetailOpts{
		InkReturnID: inkReturn.ID,
		Limit:       10000,
		Offset:      0,
	})
	if err != nil {
		return err
	}

	updatingItems := generic.ToMap(opt.InkReturnDetail, func(v *EditInkReturnDetailOpts) string {
		return v.InkID
	})

	execTx := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// modify changing ink return detail
		for i := range inkReturnItems {
			updatingItem, ok := updatingItems[inkReturnItems[i].InkID]
			if !ok {
				continue
			}

			// update returning ink
			existingQuality := inkReturnItems[i].Quantity
			inkReturnItems[i].Quantity = updatingItem.Quantity
			inkReturnItems[i].Description = cockroach.String(updatingItem.Description)
			inkReturnItems[i].UpdatedAt = now
			if err := p.inkReturnDetailRepo.Update(ctx, inkReturnItems[i].InkReturnDetail); err != nil {
				return fmt.Errorf("error when update return ink detail: %w", err)
			}

			// remove updated item out of map
			delete(updatingItems, inkReturnItems[i].InkID)

			// correct ink quality with updating quantity
			inkData, err := p.inkRepo.FindByID(c, inkReturnItems[i].InkID)
			if err != nil {
				return fmt.Errorf("màu mực không tồn tại. id: %w", err)
			}
			inkData.Quantity = inkData.Quantity - existingQuality + inkReturnItems[i].Quantity
			if inkData.Quantity < 0 {
				return fmt.Errorf("không đủ số lượng để xuất kho: trong kho còn %v", inkData.Quantity)
			}
			if err := p.inkRepo.Update(c, inkData.Ink); err != nil {
				return fmt.Errorf("error when insert ink: %w", err)
			}
		}

		// add new changing ink return detail
		for _, inkReturnDetail := range updatingItems {
			newInkReturnDetail := &model.InkReturnDetail{
				ID:                idutil.ULIDNow(),
				InkReturnID:       inkReturn.ID,
				InkExportID:       inkReturn.InkExportID, // todo remove this field
				InkID:             inkReturnDetail.InkID,
				InkExportDetailID: inkReturnDetail.InkExportDetailID,
				Quantity:          inkReturnDetail.Quantity,
				ColorDetail:       inkReturnDetail.ColorDetail,
				Description:       cockroach.String(inkReturnDetail.Description),
				Data:              inkReturnDetail.Data,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			if err := p.inkReturnDetailRepo.Insert(c, newInkReturnDetail); err != nil {
				return fmt.Errorf("error when insert ink import detail: %w", err)
			}

			// check exist inkID and inkExportID in ink_export_detail
			inkExportDetail, err := p.inkExportDetailRepo.Search(c, &repository.SearchInkExportDetailOpts{
				InkExportID: inkReturn.InkExportID,
				InkID:       inkReturnDetail.InkID,
				Limit:       1,
				Offset:      0,
			})
			if err != nil {
				return fmt.Errorf("phiếu xuất kho và màu mực không tồn tại: %w", err)
			}
			if len(inkExportDetail) == 0 {
				return fmt.Errorf("phiếu xuất kho và màu mực không tồn tại")
			}

			// check if inkExportDetail.InkID exist in ink.ID
			inkData, err := p.inkRepo.FindByID(c, inkReturnDetail.InkID)
			if err != nil {
				return fmt.Errorf("màu mực không tồn tại. id: %w", err)
			}

			// update ink quantity
			inkData.Quantity = inkData.Quantity + inkReturnDetail.Quantity
			if err := p.inkRepo.Update(c, inkData.Ink); err != nil {
				return fmt.Errorf("error when insert ink: %w", err)
			}
		}

		// update ink return
		if err := p.inkReturnRepo.Update(ctx, inkReturn); err != nil {
			return fmt.Errorf("error when update ink return: %w", err)
		}

		return nil
	})
	if execTx != nil {
		return execTx
	}

	return nil
}

func (p inkReturnService) Create(ctx context.Context, opt *CreateInkReturnOpts) (string, error) {
	// write code to insert to ink_import table and insert to ink_import_detail table in transaction
	returnId := idutil.ULIDNow()
	now := time.Now()
	startOfDay := now.Truncate(time.Hour * 24)
	endOfDay := startOfDay.Add(time.Hour*24 - time.Nanosecond)

	err := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// get count ink return in day
		count, err := p.inkReturnRepo.Count(c, &repository.SearchInkReturnOpts{
			DateFrom: startOfDay,
			DateTo:   endOfDay,
		})
		if err != nil {
			return fmt.Errorf("error when count ink return: %w", err)
		}
		// insert to ink_import
		err1 := p.inkReturnRepo.Insert(c, &model.InkReturn{
			ID:          returnId,
			Name:        opt.Name,
			Code:        opt.Code + fmt.Sprintf("%03d", count.Count+1),
			InkExportID: opt.InkExportID,
			ReturnDate:  cockroach.Time(now),
			Description: cockroach.String(opt.Description),
			Status:      enum.InventoryCommonStatusStatusCompleted, // default status is completed
			Data:        opt.Data,
			CreatedBy:   opt.CreatedBy,
			UpdatedBy:   opt.CreatedBy,
			CreatedAt:   now,
			UpdatedAt:   now,
		})

		if err1 != nil {
			return fmt.Errorf("error when insert ink import: %w", err1)
		}
		// write code to insert to ink_import_detail table
		for _, inkReturnDetail := range opt.InkReturnDetail {
			err2 := p.inkReturnDetailRepo.Insert(c, &model.InkReturnDetail{
				ID:                idutil.ULIDNow(),
				InkReturnID:       returnId,
				InkID:             inkReturnDetail.InkID,
				InkExportID:       opt.InkExportID, // todo remove this field
				InkExportDetailID: inkReturnDetail.InkExportDetailID,
				Quantity:          inkReturnDetail.Quantity,
				ColorDetail:       inkReturnDetail.ColorDetail,
				Description:       cockroach.String(inkReturnDetail.Description),
				Data:              inkReturnDetail.Data,
				CreatedAt:         now,
				UpdatedAt:         now,
			})
			if err2 != nil {
				return fmt.Errorf("error when insert ink import detail: %w", err2)
			}
			// todo check if import status is completed, insert to ink table
			// check exist inkID and inkExportID in ink_export_detail
			inkExportDetail, err := p.inkExportDetailRepo.Search(c, &repository.SearchInkExportDetailOpts{
				InkExportID: opt.InkExportID,
				InkID:       inkReturnDetail.InkID,
				Limit:       1,
				Offset:      0,
			})
			if err != nil {
				return fmt.Errorf("phiếu xuất kho và màu mực không tồn tại: %w", err)
			}
			if len(inkExportDetail) == 0 {
				return fmt.Errorf("phiếu xuất kho và màu mực không tồn tại")
			}
			// check if inkExportDetail.InkID exist in ink.ID
			inkData, err := p.inkRepo.FindByID(c, inkReturnDetail.InkID)
			if err != nil {
				return fmt.Errorf("màu mực không tồn tại. id: %w", err)
			}
			// update ink quantity
			inkData.Quantity = inkData.Quantity + inkReturnDetail.Quantity

			err3 := p.inkRepo.Update(c, inkData.Ink)
			if err3 != nil {
				return fmt.Errorf("error when insert ink: %w", err3)
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}
	return returnId, nil
}

func (p inkReturnService) Delete(ctx context.Context, id string) error {
	return p.inkReturnRepo.SoftDelete(ctx, id)
}

type InkReturnData struct {
	ID              string
	Name            string
	Code            string
	InkExportID     string
	ReturnDate      time.Time
	ReturnWarehouse string
	Description     string
	Status          enum.InventoryCommonStatus
	Data            map[string]interface{}
	CreatedBy       string
	UpdatedBy       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CreatedByName   string
	UpdatedByName   string
	InkReturnDetail []*InkReturnDetail
}

type InkReturnDetail struct {
	ID                string
	InkReturnID       string // FK to ink_return
	InkExportDetailID string // FK to ink_export_detail
	InkID             string
	InkData           *repository.InkData
	Quantity          float64
	ColorDetail       map[string]interface{}
	Description       string
	Data              map[string]interface{}
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p inkReturnService) Find(ctx context.Context, opt *FindInkReturnOpts, sort *repository.Sort, limit, offset int64) ([]*InkReturnData, *repository.CountResult, error) {
	filter := &repository.SearchInkReturnOpts{
		Name:         opt.Name,
		Status:       opt.Status,
		InkExportID:  opt.InkExportID,
		InkExportIDs: opt.InkExportIDs,
		ID:           opt.ID,
		Limit:        limit,
		Offset:       offset,
		Sort:         sort,
	}
	inkReturns, err := p.inkReturnRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*InkReturnData, 0)
	// write code to get ink_import_detail
	for _, inkReturn := range inkReturns {
		data := &InkReturnData{
			ID:              inkReturn.ID,
			Name:            inkReturn.Name,
			Code:            inkReturn.Code,
			Description:     inkReturn.Description.String,
			Status:          inkReturn.Status,
			Data:            inkReturn.Data,
			InkExportID:     inkReturn.InkExportID,
			CreatedBy:       inkReturn.CreatedBy,
			UpdatedBy:       inkReturn.UpdatedBy,
			CreatedByName:   inkReturn.CreatedByName,
			UpdatedByName:   inkReturn.UpdatedByName,
			CreatedAt:       inkReturn.CreatedAt,
			UpdatedAt:       inkReturn.UpdatedAt,
			InkReturnDetail: nil,
		}
		inkReturnDetails, err := p.inkReturnDetailRepo.Search(ctx, &repository.SearchInkReturnDetailOpts{
			InkReturnID: inkReturn.ID,
			Limit:       10000,
		})
		if err != nil {
			return nil, nil, err
		}
		inkReturnDetailResults := make([]*InkReturnDetail, 0)
		for _, inkReturnDetail := range inkReturnDetails {
			// find inkdata
			inkData, _ := p.inkRepo.FindByID(ctx, inkReturnDetail.InkID)
			dataDetail := &InkReturnDetail{
				ID:                inkReturnDetail.ID,
				InkReturnID:       inkReturnDetail.InkReturnID,
				InkData:           inkData,
				InkExportDetailID: inkReturnDetail.InkExportDetailID,
				InkID:             inkReturnDetail.InkID,
				Quantity:          inkReturnDetail.Quantity,
				ColorDetail:       inkReturnDetail.ColorDetail,
				Description:       inkReturnDetail.Description.String,
				Data:              inkReturnDetail.Data,
				CreatedAt:         inkReturnDetail.CreatedAt,
				UpdatedAt:         inkReturnDetail.UpdatedAt,
			}
			inkReturnDetailResults = append(inkReturnDetailResults, dataDetail)
		}
		data.InkReturnDetail = inkReturnDetailResults
		results = append(results, data)
	}
	total, err := p.inkReturnRepo.Count(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	return results, total, nil
}

func NewService(
	inkReturnRepo repository.InkReturnRepo,
	inkReturnDetailRepo repository.InkReturnDetailRepo,
	inkRepo repository.InkRepo,
	inkExportDetailRepo repository.InkExportDetailRepo,
) Service {
	return &inkReturnService{
		inkExportDetailRepo: inkExportDetailRepo,
		inkReturnRepo:       inkReturnRepo,
		inkReturnDetailRepo: inkReturnDetailRepo,
		inkRepo:             inkRepo,
	}
}
