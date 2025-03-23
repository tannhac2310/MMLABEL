package product_quality

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *productQualityService) CreateProductQuality(ctx context.Context, opt *CreateProductQualityOpts) (string, error) {
	formID := ""
	countAll, err := c.inspectionFormRepo.CountAll(ctx)
	if err != nil {
		return "", fmt.Errorf("c.inspectionFormRepo.CountAll: %w", err)
	}
	formID = fmt.Sprintf("OQC-%d", *countAll+1)
	now := time.Now()
	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.inspectionFormRepo.Insert(ctx2, &model.InspectionForm{
			ID:                  formID,
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
			CreatedBy:           opt.CreatedBy,
			UpdatedBy:           opt.CreatedBy,
			CreatedAt:           now,
			UpdatedAt:           now,
		})

		if err != nil {
			return fmt.Errorf("c.inspectionFormRepo.Insert with id %s: %w", formID, err)
		}

		for _, e := range opt.InspectionErrors {
			err2 := c.inspectionErrorRepo.Insert(ctx2, &model.InspectionError{
				ID:               idutil.ULIDNow(),
				DeviceID:         e.DeviceID,
				DeviceName:       e.DeviceName,
				InspectionFormID: formID,
				ErrorType:        e.ErrorType,
				Quantity:         e.Quantity,
				Note:             e.Note,
				NhanVienThucHien: e.NhanVienThucHien,
				CreatedBy:        opt.CreatedBy,
				UpdatedBy:        opt.CreatedBy,
				CreatedAt:        now,
				UpdatedAt:        now,
			})
			if err2 != nil {
				return fmt.Errorf("c.inspectionErrorRepo.Insert: %w", err2)
			}
		}

		return nil
	})

	if errTx != nil {
		return "", fmt.Errorf("cockroach.ExecInTx: %w", errTx)
	}
	return formID, nil
}

type CreateProductQualityOpts struct {
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
	CreatedBy           string
	InspectionErrors    []*InspectionError
}

type InspectionError struct {
	DeviceID         string
	DeviceName       string
	InspectionFormID string
	ErrorType        string
	Quantity         int64
	Note             string
	NhanVienThucHien string
}
