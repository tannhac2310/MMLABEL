package device_config

import (
	"context"
	"fmt"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (c *deviceConfigService) EditDeviceConfig(ctx context.Context, opt *EditDeviceConfigOpts) error {
	var err error
	table := model.ProductionOrderDeviceConfig{}
	updater := cockroach.NewUpdater(table.TableName(), model.ProductionOrderDeviceConfigFieldID, opt.ID)

	updater.Set(model.ProductionOrderDeviceConfigFieldDeviceID, opt.DeviceID)
	updater.Set(model.ProductionOrderDeviceConfigFieldDeviceConfig, opt.DeviceConfig)
	updater.Set(model.ProductionOrderDeviceConfigFieldColor, opt.Color)
	updater.Set(model.ProductionOrderDeviceConfigFieldDescription, opt.Description)
	updater.Set(model.ProductionOrderDeviceConfigFieldSearch, opt.Search)
	updater.Set(model.ProductionOrderDeviceConfigFieldMaThongSoMay, opt.MaThongSoMay)
	updater.Set(model.ProductionOrderDeviceConfigFieldMaTaiLieuHuongDan, opt.MaTaiLieuHuongDan)
	updater.Set(model.ProductionOrderDeviceConfigFieldNgayHieuLuc, opt.NgayHieuLuc)
	updater.Set(model.ProductionOrderDeviceConfigFieldStageID, opt.StageID)
	updater.Set(model.ProductionOrderDeviceConfigFieldMaSanPham, opt.MaSanPham)
	updater.Set(model.ProductionOrderDeviceConfigFieldMaSanPhamNoiBo, opt.MaSanPhamNoiBo)
	updater.Set(model.ProductionOrderDeviceConfigFieldChuKyIn, opt.ChuKyIn)
	updater.Set(model.ProductionOrderDeviceConfigFieldThoiGianChuanBi, opt.ThoiGianChuanBi)
	updater.Set(model.ProductionOrderDeviceConfigFieldTenMauMuc, opt.TenMauMuc)
	updater.Set(model.ProductionOrderDeviceConfigFieldTenLoaiMuc, opt.TenLoaiMuc)
	updater.Set(model.ProductionOrderDeviceConfigFieldSoThuTuIn, opt.SoThuTuIn)
	updater.Set(model.ProductionOrderDeviceConfigFieldMaPhim, opt.MaPhim)
	updater.Set(model.ProductionOrderDeviceConfigFieldMaMauMuc, opt.MaMauMuc)
	updater.Set(model.ProductionOrderDeviceConfigFieldTinhTrangMuc, opt.TinhTrangMuc)
	updater.Set(model.ProductionOrderDeviceConfigFieldDienTichPhuMuc, opt.DienTichPhuMuc)
	updater.Set(model.ProductionOrderDeviceConfigFieldDungMoi, opt.DungMoi)
	updater.Set(model.ProductionOrderDeviceConfigFieldNhietDoSay, opt.NhietDoSay)
	updater.Set(model.ProductionOrderDeviceConfigFieldThoiGianSay, opt.ThoiGianSay)
	updater.Set(model.ProductionOrderDeviceConfigFieldGhiChu, opt.GhiChu)
	updater.Set(model.ProductionOrderDeviceConfigFieldMaKhung, opt.MaKhung)
	updater.Set(model.ProductionOrderDeviceConfigFieldThongSoLua, opt.ThongSoLua)
	updater.Set(model.ProductionOrderDeviceConfigFieldKhoangCachKhungInBanIn, opt.KhoangCachKhungInBanIn)
	updater.Set(model.ProductionOrderDeviceConfigFieldCachIn, opt.CachIn)
	updater.Set(model.ProductionOrderDeviceConfigFieldCungDao, opt.CungDao)
	updater.Set(model.ProductionOrderDeviceConfigFieldDoBenDao, opt.DoBenDao)
	updater.Set(model.ProductionOrderDeviceConfigFieldDoNghiengDao, opt.DoNghiengDao)
	updater.Set(model.ProductionOrderDeviceConfigFieldTocDoDao, opt.TocDoDao)
	updater.Set(model.ProductionOrderDeviceConfigFieldTocDo, opt.TocDo)
	updater.Set(model.ProductionOrderDeviceConfigFieldInkID, opt.InkID)

	updater.Set(model.ProductionOrderDeviceConfigFieldUpdatedAt, time.Now())
	updater.Set(model.ProductionOrderDeviceConfigFieldUpdatedBy, opt.UpdatedBy)

	err = cockroach.UpdateFields(ctx, updater)
	if err != nil {
		return fmt.Errorf("update deviceConfig failed %w", err)
	}
	return nil
}

type EditDeviceConfigOpts struct {
	ID                     string
	DeviceConfig           map[string]interface{}
	Color                  string
	Description            string
	Search                 string
	DeviceID               string
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
	InkID                  string
	UpdatedBy              string
}
