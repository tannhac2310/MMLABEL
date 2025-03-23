package product_quality

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
)

func (c *productQualityService) EditProductQuality(ctx context.Context, opt *EditProductQualityOpts) error {
	inspectData, err := c.inspectionFormRepo.FindByID(ctx, opt.ID)
	if err != nil {
		return fmt.Errorf("c.inspectionFormRepo.FindByID: %w", err)
	}
	return cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {

		err := c.inspectionFormRepo.Update(ctx2, &model.InspectionForm{
			ID:                  opt.ID,
			ProductionOrderID:   opt.ProductionOrderID,
			InspectionDate:      opt.InspectionDate,
			InspectorName:       opt.InspectorName,
			Quantity:            opt.Quantity,
			ProductID:           opt.ProductID,
			SoLuongHopDong:      opt.SoLuongHopDong,
			SoLuongIn:           opt.SoLuongIn,
			NguoiKiemTra:        opt.NguoiKiemTra,
			NguoiPheDuyet:       opt.NguoiPheDuyet,
			SoLuongThanhPhamDat: opt.SoLuongThanhPhamDat,
			Note:                opt.Note,
			CreatedBy:           inspectData.CreatedBy,
			UpdatedBy:           opt.CreatedBy,
			CreatedAt:           inspectData.CreatedAt,
			UpdatedAt:           time.Now(),
		})
		if err != nil {
			return fmt.Errorf("c.inspectionFormRepo.Update: %w", err)
		}

		err = c.inspectionErrorRepo.SoftDeleteByFormID(ctx2, opt.ID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			//return fmt.Errorf("c.inspectionErrorRepo.DeleteByFormID: %w", err)
		}

		for _, e := range opt.InspectionErrors {
			err = c.inspectionErrorRepo.Insert(ctx2, &model.InspectionError{
				ID:               idutil.ULIDNow(),
				DeviceID:         e.DeviceID,
				DeviceName:       e.DeviceName,
				InspectionFormID: opt.ID,
				ErrorType:        e.ErrorType,
				Quantity:         e.Quantity,
				NhanVienThucHien: e.NhanVienThucHien,
				Note:             e.Note,
				CreatedBy:        opt.CreatedBy,
				UpdatedBy:        opt.CreatedBy,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			})
		}
		return nil
	})
}

type EditProductQualityOpts struct {
	ID                  string
	ProductionOrderID   string
	InspectionDate      time.Time
	InspectorName       string
	Quantity            int64
	Note                string
	ProductID           string
	SoLuongHopDong      int64
	SoLuongIn           int64
	NguoiKiemTra        string
	NguoiPheDuyet       string
	SoLuongThanhPhamDat int64
	InspectionErrors    []*EditInspectionError
	CreatedBy           string
}

type EditInspectionError struct {
	DeviceID         string
	DeviceName       string
	ErrorType        string
	Quantity         int64
	Note             string
	NhanVienThucHien string
}
