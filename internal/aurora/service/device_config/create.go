package device_config

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"
)

func (c *deviceConfigService) CreateDeviceConfig(ctx context.Context, opt *CreateDeviceConfigOpts) (string, error) {
	now := time.Now()

	deviceConfig := &model.ProductionOrderDeviceConfig{
		ID:                     idutil.ULIDNow(),
		ProductionOrderID:      cockroach.String(opt.ProductionOrderID),
		DeviceID:               cockroach.String(opt.DeviceID),
		Color:                  cockroach.String(opt.Color),
		Description:            cockroach.String(opt.Description),
		Search:                 fmt.Sprintf("%s ", opt.Color),
		DeviceConfig:           opt.DeviceConfig,
		CreatedBy:              opt.CreatedBy,
		CreatedAt:              now,
		UpdatedBy:              opt.CreatedBy,
		UpdatedAt:              now,
		ProductionPlanID:       cockroach.String(opt.ProductionPlanID),
		DeviceType:             opt.DeviceType,
		MaThongSoMay:           cockroach.String(opt.MaThongSoMay),
		MaTaiLieuHuongDan:      cockroach.String(opt.MaTaiLieuHuongDan),
		NgayHieuLuc:            cockroach.Time(opt.NgayHieuLuc),
		StageID:                cockroach.String(opt.StageID),
		MaSanPham:              cockroach.String(opt.MaSanPham),
		MaSanPhamNoiBo:         cockroach.String(opt.MaSanPhamNoiBo),
		ChuKyIn:                cockroach.String(opt.ChuKyIn),
		ThoiGianChuanBi:        cockroach.String(opt.ThoiGianChuanBi),
		TenMauMuc:              cockroach.String(opt.TenMauMuc),
		TenLoaiMuc:             cockroach.String(opt.TenLoaiMuc),
		SoThuTuIn:              cockroach.String(opt.SoThuTuIn),
		MaPhim:                 cockroach.String(opt.MaPhim),
		MaMauMuc:               cockroach.String(opt.MaMauMuc),
		TinhTrangMuc:           cockroach.String(opt.TinhTrangMuc),
		DienTichPhuMuc:         cockroach.String(opt.DienTichPhuMuc),
		DungMoi:                cockroach.String(opt.DungMoi),
		NhietDoSay:             cockroach.String(opt.NhietDoSay),
		ThoiGianSay:            cockroach.String(opt.ThoiGianSay),
		GhiChu:                 cockroach.String(opt.GhiChu),
		MaKhung:                cockroach.String(opt.MaKhung),
		ThongSoLua:             cockroach.String(opt.ThongSoLua),
		KhoangCachKhungInBanIn: cockroach.String(opt.KhoangCachKhungInBanIn),
		CachIn:                 cockroach.String(opt.CachIn),
		CungDao:                cockroach.String(opt.CungDao),
		DoBenDao:               cockroach.String(opt.DoBenDao),
		DoNghiengDao:           cockroach.String(opt.DoNghiengDao),
		TocDoDao:               cockroach.String(opt.TocDoDao),
		TocDo:                  cockroach.String(opt.TocDo),
		InkID:                  cockroach.String(opt.InkID),
		//SanPhamNguon:           sql.NullString{},
		//BanInNguon:             sql.NullString{},
	}
	fmt.Printf("deviceConfig: %+v\n", deviceConfig)

	errTx := cockroach.ExecInTx(ctx, func(ctx2 context.Context) error {
		err := c.deviceConfigRepo.Insert(ctx2, deviceConfig)
		if err != nil {
			return fmt.Errorf("c.deviceConfigRepo.Insert: %w", err)
		}

		return nil
	})
	if errTx != nil {
		return "", errTx
	}
	return deviceConfig.ID, nil
}

type CreateDeviceConfigOpts struct {
	ProductionOrderID      string
	ProductionPlanID       string
	DeviceID               string
	Color                  string
	MaThongSoMay           string
	MaTaiLieuHuongDan      string
	NgayHieuLuc            time.Time
	StageID                string
	MaSanPham              string
	MaSanPhamNoiBo         string
	ChuKyIn                string
	ThoiGianChuanBi        string
	TenMauMuc              string
	TenLoaiMuc             string
	SoThuTuIn              string
	MaPhim                 string
	MaMauMuc               string
	TinhTrangMuc           string
	DienTichPhuMuc         string
	DungMoi                string
	NhietDoSay             string
	ThoiGianSay            string
	GhiChu                 string
	MaKhung                string
	ThongSoLua             string
	KhoangCachKhungInBanIn string
	CachIn                 string
	CungDao                string
	DoBenDao               string
	DoNghiengDao           string
	TocDoDao               string
	TocDo                  string
	Description            string
	DeviceType             enum.DeviceConfigType
	DeviceConfig           map[string]interface{}
	InkID                  string
	CreatedBy              string
}
