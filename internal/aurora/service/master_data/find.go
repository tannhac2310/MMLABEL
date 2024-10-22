package masterdata

import (
	"context"
	"fmt"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type FindMasterDatasOpts struct {
	IDs  []string
	Type enum.MasterDataType
}

func (s *masterDataService) FindMasterDatas(ctx context.Context, opts *FindMasterDatasOpts, sort *repository.Sort, limit, offset int64) ([]*Data, *repository.CountResult, error) {
	switch opts.Type {
	case enum.MasterDataType_KhachHang:
		filter := &repository.SearchMKhachHangOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mKhachHangRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mKhachHangRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	case enum.MasterDataType_KhungIn:
		filter := &repository.SearchMKhungInOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mKhungInRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mKhungInRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	case enum.MasterDataType_KhuonBe:
		filter := &repository.SearchMKhuonBeOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mKhuongBeRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mKhuongBeRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	case enum.MasterDataType_KhuonDap:
		filter := &repository.SearchMKhuonDapOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mKhuonDapRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mKhuonDapRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	case enum.MasterDataType_NguyenVatLieu:
		filter := &repository.SearchMNguyenVatLieuOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mNguyenVatLieuRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mNguyenVatLieuRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	case enum.MasterDataType_Phim:
		filter := &repository.SearchMPhimOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mPhimRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mPhimRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	case enum.MasterDataType_ThongSoMayIn:
		filter := &repository.SearchMThongSoMayInOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mThongSoMayInRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mThongSoMayInRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	case enum.MasterDataType_ThongSoMayKhac:
		filter := &repository.SearchMThongSoMayKhacOpts{
			IDs:    opts.IDs,
			Limit:  limit,
			Offset: offset,
			Sort:   sort,
		}
		total, err := s.mThongSoMayKhacRepo.Count(ctx, filter)
		if err != nil {
			return nil, nil, err
		} else if total.Count == 0 {
			return []*Data{}, total, nil
		}

		list, err := s.mThongSoMayKhacRepo.Search(ctx, filter)
		if err != nil {
			return nil, nil, err
		}

		res := make([]*Data, 0)
		for _, item := range list {
			res = append(res, &Data{item})
		}
		return res, total, nil
	default:
		return nil, nil, fmt.Errorf("invalid master data type %s", enum.MasterDataTypeName[opts.Type])
	}
}
