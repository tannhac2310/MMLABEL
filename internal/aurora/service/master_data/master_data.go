package masterdata

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
)

type Service interface {
	FindMasterDatas(ctx context.Context, opt *FindMasterDatasOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error)
	CreateMasterData(ctx context.Context, opts *CreateProductionOrderOpts) (string, error)
}

var _ Service = (*masterDataService)(nil)

type masterDataService struct {
	mKhachHangRepo      repository.MKhachHangRepo
	mKhungInRepo        repository.MKhungInRepo
	mKhuongBeRepo       repository.MKhuonBeRepo
	mKhuonDapRepo       repository.MKhuonDapRepo
	mNguyenVatLieuRepo  repository.MNguyenVatLieuRepo
	mPhimRepo           repository.MPhimRepo
	mThongSoMayInRepo   repository.MThongSoMayInRepo
	mThongSoMayKhacRepo repository.MThongSoMayKhacRepo
}

func NewMasterDataService(
	mKhachHangRepo repository.MKhachHangRepo,
	mKhungInRepo repository.MKhungInRepo,
	mKhuongBeRepo repository.MKhuonBeRepo,
	mKhuonDapRepo repository.MKhuonDapRepo,
	mNguyenVatLieuRepo repository.MNguyenVatLieuRepo,
	mPhimRepo repository.MPhimRepo,
	mThongSoMayInRepo repository.MThongSoMayInRepo,
	mThongSoMayKhacRepo repository.MThongSoMayKhacRepo,
) *masterDataService {
	return &masterDataService{
		mKhachHangRepo:      mKhachHangRepo,
		mKhungInRepo:        mKhungInRepo,
		mKhuongBeRepo:       mKhuongBeRepo,
		mKhuonDapRepo:       mKhuonDapRepo,
		mNguyenVatLieuRepo:  mNguyenVatLieuRepo,
		mPhimRepo:           mPhimRepo,
		mThongSoMayInRepo:   mThongSoMayInRepo,
		mThongSoMayKhacRepo: mThongSoMayKhacRepo,
	}
}

type Data struct {
	any
}
