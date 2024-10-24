package masterdata

import (
	"context"
	"fmt"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

type CreateProductionOrderOpts struct {
	Type      enum.MasterDataType
	Data      any
	CreatedBy string
}

func (s *masterDataService) CreateMasterData(ctx context.Context, opts *CreateProductionOrderOpts) (string, error) {
	now := time.Now()
	id := idutil.ULIDNow()

	switch opts.Type {
	case enum.MasterDataType_KhachHang:
		var newItem model.MKhachHangData
		if err := mapstructure.Decode(opts.Data, &newItem); err != nil {
			return "", fmt.Errorf("parse payload failed: %w", err)
		}

		err := s.mKhachHangRepo.Insert(ctx, &model.MKhachHang{
			ID:        id,
			Data:      newItem,
			CreatedBy: opts.CreatedBy,
			CreatedAt: now,
			UpdatedBy: opts.CreatedBy,
			UpdatedAt: now,
		})
		if err != nil {
			return "", err
		}

		return id, nil
	case enum.MasterDataType_KhungIn:
		return "", fmt.Errorf("unsupported yet")
	case enum.MasterDataType_KhuonBe:
		return "", fmt.Errorf("unsupported yet")
	case enum.MasterDataType_KhuonDap:
		return "", fmt.Errorf("unsupported yet")
	case enum.MasterDataType_NguyenVatLieu:
		return "", fmt.Errorf("unsupported yet")
	case enum.MasterDataType_Phim:
		return "", fmt.Errorf("unsupported yet")
	case enum.MasterDataType_ThongSoMayIn:
		return "", fmt.Errorf("unsupported yet")
	case enum.MasterDataType_ThongSoMayKhac:
		return "", fmt.Errorf("unsupported yet")
	default:
		return "", fmt.Errorf("invalid master data type %s", enum.MasterDataTypeName[opts.Type])
	}
}
