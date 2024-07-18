package ink_export

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/ink_return"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/generic"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type EditInkExportDetailOpts struct {
	InkID       string
	Quantity    float64
	Description string
	Data        map[string]interface{}
}

type EditInkExportOpts struct {
	ID              string
	Description     string
	UpdatedBy       string
	InkExportDetail []*EditInkExportDetailOpts
}

type CreateInkExportDetailOpts struct {
	InkID       string
	Quantity    float64
	Description string
	Data        map[string]interface{}
}

type CreateInkExportOpts struct {
	Name              string
	Code              string
	ProductionOrderID string
	ExportDate        string
	Description       string
	Data              map[string]interface{}
	CreatedBy         string
	InkExportDetail   []*CreateInkExportDetailOpts
}

type FindInkExportOpts struct {
	Search      string
	InkCode     string
	ID          string
	ProductName string
	Status      enum.InventoryCommonStatus
}

type Service interface {
	Edit(ctx context.Context, opt *EditInkExportOpts) error
	Create(ctx context.Context, opt *CreateInkExportOpts) (string, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, opt *FindInkExportOpts, sort *repository.Sort, limit, offset int64) ([]*InkExportData, *repository.CountResult, error)
	FindImportDetailByPOID(ctx context.Context, poID string) ([]*InkExportDetail, error)
}

type inkExportService struct {
	inkExportRepo       repository.InkExportRepo
	inkExportDetailRepo repository.InkExportDetailRepo
	inkRepo             repository.InkRepo
	productionOrderRepo repository.ProductionOrderRepo
	inkReturnSvc        ink_return.Service
}

func (p inkExportService) FindImportDetailByPOID(ctx context.Context, poID string) ([]*InkExportDetail, error) {
	inkExportDetails, err := p.inkExportDetailRepo.Search(ctx, &repository.SearchInkExportDetailOpts{
		ProductionOrderID: poID,
		Limit:             10000,
	})
	if err != nil {
		return nil, err
	}
	results := make([]*InkExportDetail, 0)
	for _, inkExportDetail := range inkExportDetails {
		inkData, _ := p.inkRepo.FindByID(ctx, inkExportDetail.InkID)
		dataDetail := &InkExportDetail{
			ID:          inkExportDetail.ID,
			InkExportID: inkExportDetail.InkExportID,
			InkID:       inkExportDetail.InkID,
			InkData:     inkData,
			Quantity:    inkData.Quantity,
			ColorDetail: inkData.ColorDetail,
			Description: inkExportDetail.Description.String,
			Data:        inkExportDetail.Data,
			CreatedAt:   inkExportDetail.CreatedAt,
			UpdatedAt:   inkExportDetail.UpdatedAt,
		}
		results = append(results, dataDetail)
	}
	return results, nil
}

func (p inkExportService) Edit(ctx context.Context, opt *EditInkExportOpts) error {
	now := time.Now()

	inkExport, err := p.inkExportRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return err
	}
	inkExport.Description = cockroach.String(opt.Description)
	inkExport.UpdatedBy = opt.UpdatedBy
	inkExport.UpdatedAt = now

	inkExportItems, err := p.inkExportDetailRepo.Search(ctx, &repository.SearchInkExportDetailOpts{
		InkExportID: inkExport.ID,
		Limit:       10000,
		Offset:      0,
	})
	if err != nil {
		return err
	}

	inkReturns, _, err := p.inkReturnSvc.Find(ctx, &ink_return.FindInkReturnOpts{
		InkExportID: inkExport.ID,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, 1000, 0)
	if err != nil {
		return fmt.Errorf("error when find ink return data: %w", err)
	}

	// inkReturnQuantity is a map of ink_id and quality
	inkReturnQuantity := make(map[string]float64, 0)
	for _, inkReturn := range inkReturns {
		for _, inkReturnDetail := range inkReturn.InkReturnDetail {
			inkReturnItem, ok := inkReturnQuantity[inkReturnDetail.InkID]
			if !ok {
				inkReturnItem = 0
			}
			inkReturnItem += inkReturnDetail.Quantity
		}
	}

	updatingItems := generic.ToMap(opt.InkExportDetail, func(v *EditInkExportDetailOpts) string {
		return v.InkID
	})

	execTx := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// modify changing ink export detail
		for i := range inkExportItems {
			updatingItem, ok := updatingItems[inkExportItems[i].InkID]
			if !ok {
				continue
			}

			// update export ink
			returnQuantity := inkReturnQuantity[inkExportItems[i].InkID]
			if returnQuantity > updatingItem.Quantity {
				return fmt.Errorf("số lượng màu mực đang chỉnh sửa phải lớn hơn hoặc bằng số lượng đã trả")
			}

			inkData, err := p.inkRepo.FindByID(c, inkExportItems[i].InkID)
			if err != nil {
				return fmt.Errorf("màu mực %s không tồn tại: %w", inkExportItems[i].InkID, err)
			}
			inkData.Quantity = inkData.Quantity + inkExportItems[i].Quantity - updatingItem.Quantity
			if inkData.Quantity < 0 {
				return fmt.Errorf("không đủ số lượng để xuất kho: trong kho còn %v", inkData.Quantity)
			}

			inkExportItems[i].Quantity = updatingItem.Quantity
			inkExportItems[i].Description = cockroach.String(updatingItem.Description)
			inkExportItems[i].UpdatedAt = now
			if err := p.inkExportDetailRepo.Update(ctx, inkExportItems[i].InkExportDetail); err != nil {
				return fmt.Errorf("error when update return ink detail: %w", err)
			}

			// remove updated item out of map
			delete(updatingItems, inkExportItems[i].InkID)

			// correct ink quality with updating quantity
			if err := p.inkRepo.Update(c, inkData.Ink); err != nil {
				return fmt.Errorf("error when insert ink: %w", err)
			}
		}

		// add new changing ink export detail
		for _, inkExportDetail := range updatingItems {
			newInkExportDetail := &model.InkExportDetail{
				ID:          idutil.ULIDNow(),
				InkExportID: inkExport.ID, // todo remove this field
				InkID:       inkExportDetail.InkID,
				Quantity:    inkExportDetail.Quantity,
				Description: cockroach.String(inkExportDetail.Description),
				Data:        inkExportDetail.Data,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			if err := p.inkExportDetailRepo.Insert(c, newInkExportDetail); err != nil {
				return fmt.Errorf("error when insert ink export detail: %w", err)
			}

			// correct ink quality with updating quantity
			inkData, err := p.inkRepo.FindByID(c, inkExportDetail.InkID)
			if err != nil {
				return fmt.Errorf("error when find ink by id: %w", err)
			}
			inkData.Quantity = inkData.Quantity - inkExportDetail.Quantity
			if inkData.Quantity < 0 {
				return fmt.Errorf("không đủ số lượng để xuất kho: trong kho còn %v", inkData.Quantity)
			}
			if err := p.inkRepo.Update(c, inkData.Ink); err != nil {
				return fmt.Errorf("error when insert ink: %w", err)
			}
		}

		// update ink return
		if err := p.inkExportRepo.Update(ctx, inkExport); err != nil {
			return fmt.Errorf("error when update ink export: %w", err)
		}

		return nil
	})
	if execTx != nil {
		return execTx
	}

	return nil
}

func (p inkExportService) Create(ctx context.Context, opt *CreateInkExportOpts) (string, error) {
	// write code to insert to ink_export table and insert to ink_export_detail table in transaction
	exportId := idutil.ULIDNow()
	now := time.Now()
	startOfDay := now.Truncate(time.Hour * 24)
	endOfDay := startOfDay.Add(time.Hour*24 - time.Nanosecond)

	err := cockroach.ExecInTx(ctx, func(c context.Context) error {
		// get count ink_export by date
		count, err := p.inkExportRepo.Count(c, &repository.SearchInkExportOpts{
			ExportDateFrom: startOfDay,
			ExportDateTo:   endOfDay,
		})
		if err != nil {
			return fmt.Errorf("error when count ink export by date: %w", err)
		}
		// insert to ink_export
		err1 := p.inkExportRepo.Insert(c, &model.InkExport{
			ID:                exportId,
			Name:              opt.Name,
			Code:              opt.Code + fmt.Sprintf("%03d", count.Count+1), // hopefully we don't face race condition
			ProductionOrderID: opt.ProductionOrderID,
			ExportDate:        cockroach.Time(now),
			Description:       cockroach.String(opt.Description),
			Status:            enum.InventoryCommonStatusStatusCompleted, // default status is completed
			Data:              opt.Data,
			CreatedBy:         opt.CreatedBy,
			UpdatedBy:         opt.CreatedBy,
			CreatedAt:         now,
			UpdatedAt:         now,
		})

		if err1 != nil {
			return fmt.Errorf("error when insert ink import: %w", err1)
		}
		// write code to insert to ink_export_detail table
		for _, inkExportDetail := range opt.InkExportDetail {
			err2 := p.inkExportDetailRepo.Insert(c, &model.InkExportDetail{
				ID:          idutil.ULIDNow(),
				InkExportID: exportId,
				InkID:       inkExportDetail.InkID,
				Quantity:    inkExportDetail.Quantity,
				Description: cockroach.String(inkExportDetail.Description),
				Data:        inkExportDetail.Data,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err2 != nil {
				return fmt.Errorf("error when insert ink import detail: %w", err2)
			}
			// todo check if import status is completed, insert to ink table
			// check if inkExportDetail.InkID exist in ink.ID
			inkData, err := p.inkRepo.FindByID(c, inkExportDetail.InkID)
			if err != nil {
				return fmt.Errorf("error when find ink import detail by id: %w", err)
			}
			// update ink quantity
			inkData.Quantity = inkData.Quantity - inkExportDetail.Quantity
			if inkData.Quantity < 0 {
				return fmt.Errorf("không đủ số lượng để xuất kho: trong kho còn %v", inkData.Quantity)
			}
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
	return exportId, nil
}

func (p inkExportService) Delete(ctx context.Context, id string) error {
	return p.inkExportRepo.SoftDelete(ctx, id)
}

type InkExportData struct {
	ID                  string
	ProductionOrderID   string
	Code                string
	Name                string
	Status              enum.InventoryCommonStatus
	Description         string
	Data                map[string]interface{}
	CreatedAt           time.Time
	CreatedBy           string
	UpdatedAt           time.Time
	UpdatedBy           string
	CreatedByName       string
	UpdatedByName       string
	InkExportDetail     []*InkExportDetail
	ProductionOrderData *repository.ProductionOrderData
	InkReturnData       []*ink_return.InkReturnData
}
type InkExportDetail struct {
	ID          string
	InkExportID string
	InkID       string
	InkData     *repository.InkData
	Quantity    float64
	ColorDetail map[string]interface{}
	Description string
	Data        map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p inkExportService) Find(ctx context.Context, opt *FindInkExportOpts, sort *repository.Sort, limit, offset int64) ([]*InkExportData, *repository.CountResult, error) {
	filter := &repository.SearchInkExportOpts{
		Search:              opt.Search,
		InkCode:             opt.InkCode,
		ProductionOrderName: opt.ProductName,
		Status:              opt.Status,
		ID:                  opt.ID,
		Limit:               limit,
		Offset:              offset,
		Sort:                sort,
	}
	inkExports, err := p.inkExportRepo.Search(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	results := make([]*InkExportData, 0)

	//InkExportIDs
	inkExportIDs := make([]string, 0)
	for _, inkExport := range inkExports {
		inkExportIDs = append(inkExportIDs, inkExport.ID)
	}
	inkReturnData, _, err := p.inkReturnSvc.Find(ctx, &ink_return.FindInkReturnOpts{
		InkExportIDs: inkExportIDs,
	}, &repository.Sort{
		Order: repository.SortOrderDESC,
		By:    "ID",
	}, 1000, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("error when find ink return data: %w", err)
	}
	inkReturnDataMap := make(map[string][]*ink_return.InkReturnData)
	for _, inkReturn := range inkReturnData {
		if _, ok := inkReturnDataMap[inkReturn.InkExportID]; !ok {
			inkReturnDataMap[inkReturn.InkExportID] = make([]*ink_return.InkReturnData, 0)
		}
		inkReturnDataMap[inkReturn.InkExportID] = append(inkReturnDataMap[inkReturn.InkExportID], inkReturn)
	}

	// write code to get ink_export_detail
	for _, inkExport := range inkExports {
		data := &InkExportData{
			ID:                  inkExport.ID,
			ProductionOrderID:   inkExport.ProductionOrderID,
			Code:                inkExport.Code,
			Name:                inkExport.Name,
			Status:              inkExport.Status,
			Description:         inkExport.Description.String,
			Data:                inkExport.Data,
			CreatedAt:           inkExport.CreatedAt,
			CreatedBy:           inkExport.CreatedBy,
			CreatedByName:       inkExport.CreatedByName,
			UpdatedAt:           inkExport.UpdatedAt,
			UpdatedBy:           inkExport.UpdatedBy,
			UpdatedByName:       inkExport.UpdatedByName,
			InkExportDetail:     nil,
			ProductionOrderData: nil,
		}
		inkExportDetails, err := p.inkExportDetailRepo.Search(ctx, &repository.SearchInkExportDetailOpts{
			InkExportID: inkExport.ID,
			Limit:       10000,
		})
		if err != nil {
			return nil, nil, err
		}
		inkExportDetailResults := make([]*InkExportDetail, 0)
		for _, inkExportDetail := range inkExportDetails {
			inkData, _ := p.inkRepo.FindByID(ctx, inkExportDetail.InkID)
			dataDetail := &InkExportDetail{
				ID:          inkExportDetail.ID,
				InkExportID: inkExportDetail.InkExportID,
				InkID:       inkExportDetail.InkID,
				InkData:     inkData,
				Quantity:    inkExportDetail.Quantity,
				//ColorDetail: inkExportDetail.ColorDetail,
				Description: inkExportDetail.Description.String,
				Data:        inkExportDetail.Data,
				CreatedAt:   inkExportDetail.CreatedAt,
				UpdatedAt:   inkExportDetail.UpdatedAt,
			}
			inkExportDetailResults = append(inkExportDetailResults, dataDetail)
		}
		data.InkExportDetail = inkExportDetailResults

		// find production order data
		productionOrderData, err := p.productionOrderRepo.Search(ctx, &repository.SearchProductionOrdersOpts{
			IDs:    []string{inkExport.ProductionOrderID},
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			return nil, nil, err
		}
		if len(productionOrderData) == 1 {
			data.ProductionOrderData = productionOrderData[0]
		}
		// find ink return data
		data.InkReturnData = inkReturnDataMap[inkExport.ID]

		results = append(results, data)
	}
	total, err := p.inkExportRepo.Count(ctx, filter)

	return results, total, nil
}

func NewService(
	inkExportRepo repository.InkExportRepo,
	inkExportDetailRepo repository.InkExportDetailRepo,
	inkRepo repository.InkRepo,
	productionOrderRepo repository.ProductionOrderRepo,
	inkReturnSvc ink_return.Service,
) Service {
	return &inkExportService{
		inkExportRepo:       inkExportRepo,
		inkExportDetailRepo: inkExportDetailRepo,
		inkRepo:             inkRepo,
		productionOrderRepo: productionOrderRepo,
		inkReturnSvc:        inkReturnSvc,
	}

}
