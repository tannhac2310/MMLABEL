package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	ProductionOrderDeviceConfigFieldID                     = "id"
	ProductionOrderDeviceConfigFieldProductionOrderID      = "production_order_id"
	ProductionOrderDeviceConfigFieldDeviceID               = "device_id"
	ProductionOrderDeviceConfigFieldColor                  = "color"
	ProductionOrderDeviceConfigFieldDescription            = "description"
	ProductionOrderDeviceConfigFieldSearch                 = "search"
	ProductionOrderDeviceConfigFieldDeviceConfig           = "device_config"
	ProductionOrderDeviceConfigFieldCreatedBy              = "created_by"
	ProductionOrderDeviceConfigFieldCreatedAt              = "created_at"
	ProductionOrderDeviceConfigFieldUpdatedBy              = "updated_by"
	ProductionOrderDeviceConfigFieldUpdatedAt              = "updated_at"
	ProductionOrderDeviceConfigFieldDeletedAt              = "deleted_at"
	ProductionOrderDeviceConfigFieldProductionPlanID       = "production_plan_id"
	ProductionOrderDeviceConfigFieldDeviceType             = "device_type"
	ProductionOrderDeviceConfigFieldMaThongSoMay           = "ma_thong_so_may"
	ProductionOrderDeviceConfigFieldMaTaiLieuHuongDan      = "ma_tai_lieu_huong_dan"
	ProductionOrderDeviceConfigFieldNgayHieuLuc            = "ngay_hieu_luc"
	ProductionOrderDeviceConfigFieldStageID                = "stage_id"
	ProductionOrderDeviceConfigFieldMaSanPham              = "ma_san_pham"
	ProductionOrderDeviceConfigFieldMaSanPhamNoiBo         = "ma_san_pham_noi_bo"
	ProductionOrderDeviceConfigFieldChuKyIn                = "chu_ky_in"
	ProductionOrderDeviceConfigFieldThoiGianChuanBi        = "thoi_gian_chuan_bi"
	ProductionOrderDeviceConfigFieldTenMauMuc              = "ten_mau_muc"
	ProductionOrderDeviceConfigFieldTenLoaiMuc             = "ten_loai_muc"
	ProductionOrderDeviceConfigFieldSoThuTuIn              = "so_thu_tu_in"
	ProductionOrderDeviceConfigFieldMaPhim                 = "ma_phim"
	ProductionOrderDeviceConfigFieldMaMauMuc               = "ma_mau_muc"
	ProductionOrderDeviceConfigFieldTinhTrangMuc           = "tinh_trang_muc"
	ProductionOrderDeviceConfigFieldDienTichPhuMuc         = "dien_tich_phu_muc"
	ProductionOrderDeviceConfigFieldDungMoi                = "dung_moi"
	ProductionOrderDeviceConfigFieldNhietDoSay             = "nhiet_do_say"
	ProductionOrderDeviceConfigFieldThoiGianSay            = "thoi_gian_say"
	ProductionOrderDeviceConfigFieldGhiChu                 = "ghi_chu"
	ProductionOrderDeviceConfigFieldMaKhung                = "ma_khung"
	ProductionOrderDeviceConfigFieldThongSoLua             = "thong_so_lua"
	ProductionOrderDeviceConfigFieldKhoangCachKhungInBanIn = "khoang_cach_khung_in_ban_in"
	ProductionOrderDeviceConfigFieldCachIn                 = "cach_in"
	ProductionOrderDeviceConfigFieldCungDao                = "cung_dao"
	ProductionOrderDeviceConfigFieldDoBenDao               = "do_ben_dao"
	ProductionOrderDeviceConfigFieldDoNghiengDao           = "do_nghieng_dao"
	ProductionOrderDeviceConfigFieldTocDoDao               = "toc_do_dao"
	ProductionOrderDeviceConfigFieldTocDo                  = "toc_do"
	ProductionOrderDeviceConfigFieldSanPhamNguon           = "san_pham_nguon"
	ProductionOrderDeviceConfigFieldBanInNguon             = "ban_in_nguon"
	ProductionOrderDeviceConfigFieldInkID                  = "ink_id"
)

type ProductionOrderDeviceConfig struct {
	ID                     string                 `db:"id"`
	ProductionOrderID      sql.NullString         `db:"production_order_id"`
	DeviceID               sql.NullString         `db:"device_id"`
	Color                  sql.NullString         `db:"color"`
	Description            sql.NullString         `db:"description"`
	Search                 string                 `db:"search"`
	DeviceConfig           map[string]interface{} `db:"device_config"`
	CreatedBy              string                 `db:"created_by"`
	CreatedAt              time.Time              `db:"created_at"`
	UpdatedBy              string                 `db:"updated_by"`
	UpdatedAt              time.Time              `db:"updated_at"`
	DeletedAt              sql.NullTime           `db:"deleted_at"`
	ProductionPlanID       sql.NullString         `db:"production_plan_id"`
	DeviceType             enum.DeviceConfigType  `db:"device_type"`
	MaThongSoMay           string                 `db:"ma_thong_so_may"`
	MaTaiLieuHuongDan      sql.NullString         `db:"ma_tai_lieu_huong_dan"`
	NgayHieuLuc            time.Time              `db:"ngay_hieu_luc"`
	StageID                sql.NullString         `db:"stage_id"`
	MaSanPham              sql.NullString         `db:"ma_san_pham"`
	MaSanPhamNoiBo         sql.NullString         `db:"ma_san_pham_noi_bo"`
	ChuKyIn                sql.NullString         `db:"chu_ky_in"`
	ThoiGianChuanBi        sql.NullString         `db:"thoi_gian_chuan_bi"`
	TenMauMuc              sql.NullString         `db:"ten_mau_muc"`
	TenLoaiMuc             sql.NullString         `db:"ten_loai_muc"`
	SoThuTuIn              sql.NullString         `db:"so_thu_tu_in"`
	MaPhim                 sql.NullString         `db:"ma_phim"`
	MaMauMuc               sql.NullString         `db:"ma_mau_muc"`
	TinhTrangMuc           sql.NullString         `db:"tinh_trang_muc"`
	DienTichPhuMuc         sql.NullString         `db:"dien_tich_phu_muc"`
	DungMoi                sql.NullString         `db:"dung_moi"`
	NhietDoSay             sql.NullString         `db:"nhiet_do_say"`
	ThoiGianSay            sql.NullString         `db:"thoi_gian_say"`
	GhiChu                 sql.NullString         `db:"ghi_chu"`
	MaKhung                sql.NullString         `db:"ma_khung"`
	ThongSoLua             sql.NullString         `db:"thong_so_lua"`
	KhoangCachKhungInBanIn sql.NullString         `db:"khoang_cach_khung_in_ban_in"`
	CachIn                 sql.NullString         `db:"cach_in"`
	CungDao                sql.NullString         `db:"cung_dao"`
	DoBenDao               sql.NullString         `db:"do_ben_dao"`
	DoNghiengDao           sql.NullString         `db:"do_nghieng_dao"`
	TocDoDao               sql.NullString         `db:"toc_do_dao"`
	TocDo                  sql.NullString         `db:"toc_do"`
	SanPhamNguon           sql.NullString         `db:"san_pham_nguon"`
	BanInNguon             sql.NullString         `db:"ban_in_nguon"`
	InkID                  sql.NullString         `db:"ink_id"`
}

func (rcv *ProductionOrderDeviceConfig) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		ProductionOrderDeviceConfigFieldID,
		ProductionOrderDeviceConfigFieldProductionOrderID,
		ProductionOrderDeviceConfigFieldDeviceID,
		ProductionOrderDeviceConfigFieldColor,
		ProductionOrderDeviceConfigFieldDescription,
		ProductionOrderDeviceConfigFieldSearch,
		ProductionOrderDeviceConfigFieldDeviceConfig,
		ProductionOrderDeviceConfigFieldCreatedBy,
		ProductionOrderDeviceConfigFieldCreatedAt,
		ProductionOrderDeviceConfigFieldUpdatedBy,
		ProductionOrderDeviceConfigFieldUpdatedAt,
		ProductionOrderDeviceConfigFieldDeletedAt,
		ProductionOrderDeviceConfigFieldProductionPlanID,
		ProductionOrderDeviceConfigFieldDeviceType,
		ProductionOrderDeviceConfigFieldMaThongSoMay,
		ProductionOrderDeviceConfigFieldMaTaiLieuHuongDan,
		ProductionOrderDeviceConfigFieldNgayHieuLuc,
		ProductionOrderDeviceConfigFieldStageID,
		ProductionOrderDeviceConfigFieldMaSanPham,
		ProductionOrderDeviceConfigFieldMaSanPhamNoiBo,
		ProductionOrderDeviceConfigFieldChuKyIn,
		ProductionOrderDeviceConfigFieldThoiGianChuanBi,
		ProductionOrderDeviceConfigFieldTenMauMuc,
		ProductionOrderDeviceConfigFieldTenLoaiMuc,
		ProductionOrderDeviceConfigFieldSoThuTuIn,
		ProductionOrderDeviceConfigFieldMaPhim,
		ProductionOrderDeviceConfigFieldMaMauMuc,
		ProductionOrderDeviceConfigFieldTinhTrangMuc,
		ProductionOrderDeviceConfigFieldDienTichPhuMuc,
		ProductionOrderDeviceConfigFieldDungMoi,
		ProductionOrderDeviceConfigFieldNhietDoSay,
		ProductionOrderDeviceConfigFieldThoiGianSay,
		ProductionOrderDeviceConfigFieldGhiChu,
		ProductionOrderDeviceConfigFieldMaKhung,
		ProductionOrderDeviceConfigFieldThongSoLua,
		ProductionOrderDeviceConfigFieldKhoangCachKhungInBanIn,
		ProductionOrderDeviceConfigFieldCachIn,
		ProductionOrderDeviceConfigFieldCungDao,
		ProductionOrderDeviceConfigFieldDoBenDao,
		ProductionOrderDeviceConfigFieldDoNghiengDao,
		ProductionOrderDeviceConfigFieldTocDoDao,
		ProductionOrderDeviceConfigFieldTocDo,
		ProductionOrderDeviceConfigFieldSanPhamNguon,
		ProductionOrderDeviceConfigFieldBanInNguon,
		ProductionOrderDeviceConfigFieldInkID,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.ProductionOrderID,
		&rcv.DeviceID,
		&rcv.Color,
		&rcv.Description,
		&rcv.Search,
		&rcv.DeviceConfig,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedBy,
		&rcv.UpdatedAt,
		&rcv.DeletedAt,
		&rcv.ProductionPlanID,
		&rcv.DeviceType,
		&rcv.MaThongSoMay,
		&rcv.MaTaiLieuHuongDan,
		&rcv.NgayHieuLuc,
		&rcv.StageID,
		&rcv.MaSanPham,
		&rcv.MaSanPhamNoiBo,
		&rcv.ChuKyIn,
		&rcv.ThoiGianChuanBi,
		&rcv.TenMauMuc,
		&rcv.TenLoaiMuc,
		&rcv.SoThuTuIn,
		&rcv.MaPhim,
		&rcv.MaMauMuc,
		&rcv.TinhTrangMuc,
		&rcv.DienTichPhuMuc,
		&rcv.DungMoi,
		&rcv.NhietDoSay,
		&rcv.ThoiGianSay,
		&rcv.GhiChu,
		&rcv.MaKhung,
		&rcv.ThongSoLua,
		&rcv.KhoangCachKhungInBanIn,
		&rcv.CachIn,
		&rcv.CungDao,
		&rcv.DoBenDao,
		&rcv.DoNghiengDao,
		&rcv.TocDoDao,
		&rcv.TocDo,
		&rcv.SanPhamNguon,
		&rcv.BanInNguon,
		&rcv.InkID,
	}

	return
}

func (*ProductionOrderDeviceConfig) TableName() string {
	return "production_order_device_config"
}
