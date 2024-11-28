package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type DeviceConfigFilter struct {
	IDs               []string `json:"IDs"`
	Search            string   `json:"search"`
	ProductionOrderID string   `json:"productionOrderID"`
	ProductionPlanID  string   `json:"productionPlanID"`
	StageID           string   `json:"stageID"`
	DeviceType        string   `json:"deviceType"`
}

type FindDeviceConfigsRequest struct {
	Filter *DeviceConfigFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging   `json:"paging" binding:"required"`
}

type FindDeviceConfigsResponse struct {
	DeviceConfigs []*DeviceConfig `json:"deviceConfigs"`
	Total         int64           `json:"total"`
}
type DeviceConfig struct {
	ID                     string                `json:"id"`
	ProductionOrderID      string                `json:"productionOrderID" binding:"required"`
	ProductionOrderName    string                `json:"productionOrderName"`
	ProductionPlanID       string                `json:"productionPlanID"`
	DeviceID               string                `json:"deviceID" binding:"required"`
	DeviceName             string                `json:"deviceName"`
	DeviceType             enum.DeviceConfigType `json:"deviceType"`
	DeviceCode             string                `json:"deviceCode"`
	Color                  string                `json:"color"`
	MaThongSoMay           string                `json:"maThongSoMay"`
	MaTaiLieuHuongDan      string                `json:"maTaiLieuHuongDan"`
	NgayHieuLuc            time.Time             `json:"ngayHieuLuc"`
	StageID                string                `json:"stageID"`
	MaSanPham              string                `json:"maSanPham"`
	MaSanPhamNoiBo         string                `json:"maSanPhamNoiBo"`
	ChuKyIn                string                `json:"chuKyIn"`
	ThoiGianChuanBi        string                `json:"thoiGianChuanBi"`
	TenMauMuc              string                `json:"tenMauMuc"`
	TenLoaiMuc             string                `json:"tenLoaiMuc"`
	SoThuTuIn              string                `json:"soThuTuIn"`
	MaPhim                 string                `json:"maPhim"`
	MaMauMuc               string                `json:"maMauMuc"`
	TinhTrangMuc           string                `json:"tinhTrangMuc"`
	DienTichPhuMuc         string                `json:"dienTichPhuMuc"`
	DungMoi                string                `json:"dungMoi"`
	NhietDoSay             string                `json:"nhietDoSay"`
	ThoiGianSay            string                `json:"thoiGianSay"`
	GhiChu                 string                `json:"ghiChu"`
	MaKhung                string                `json:"maKhung"`
	ThongSoLua             string                `json:"thongSoLua"`
	KhoangCachKhungInBanIn string                `json:"khoangCachKhungInBanIn"`
	CachIn                 string                `json:"cachIn"`
	CungDao                string                `json:"cungDao"`
	DoBenDao               string                `json:"doBenDao"`
	DoNghiengDao           string                `json:"doNghiengDao"`
	TocDoDao               string                `json:"tocDoDao"`
	TocDo                  string                `json:"tocDo"`
	Description            string                `json:"description"`
	CreatedBy              string                `json:"createdBy"`
	CreatedAt              time.Time             `json:"createdAt"`
	UpdatedAt              time.Time             `json:"updatedAt"`
}

// ProductionOrderID      string
//
//	ProductionPlanID       string
//	DeviceID               string
//	Color                  string
//	MaThongSoMay           string
//	MaTaiLieuHuongDan      string
//	NgayHieuLuc            string
//	StageID                string
//	MaSanPham              string
//	MaSanPhamNoiBo         string
//	ChuKyIn                string
//	ThoiGianChuanBi        string
//	TenMauMuc              string
//	TenLoaiMuc             string
//	SoThuTuIn              string
//	MaPhim                 string
//	MaMauMuc               string
//	TinhTrangMuc           string
//	DienTichPhuMuc         string
//	DungMoi                string
//	NhietDoSay             string
//	ThoiGianSay            string
//	GhiChu                 string
//	MaKhung                string
//	ThongSoLua             string
//	KhoangCachKhungInBanIn string
//	CachIn                 string
//	CungDao                string
//	DoBenDao               string
//	DoNghiengDao           string
//	TocDoDao               string
//	TocDo                  string
//	Description            string
type CreateDeviceConfigRequest struct {
	ProductionOrderID      string                 `json:"productionOrderID"`
	ProductionPlanID       string                 `json:"productionPlanID"`
	DeviceID               string                 `json:"deviceID"`
	Color                  string                 `json:"color"`
	MaThongSoMay           string                 `json:"maThongSoMay"`
	MaTaiLieuHuongDan      string                 `json:"maTaiLieuHuongDan"`
	NgayHieuLuc            time.Time              `json:"ngayHieuLuc"`
	StageID                string                 `json:"stageID"`
	MaSanPham              string                 `json:"maSanPham"`
	MaSanPhamNoiBo         string                 `json:"maSanPhamNoiBo"`
	ChuKyIn                string                 `json:"chuKyIn"`
	ThoiGianChuanBi        string                 `json:"thoiGianChuanBi"`
	TenMauMuc              string                 `json:"tenMauMuc"`
	TenLoaiMuc             string                 `json:"tenLoaiMuc"`
	SoThuTuIn              string                 `json:"soThuTuIn"`
	MaPhim                 string                 `json:"maPhim"`
	MaMauMuc               string                 `json:"maMauMuc"`
	TinhTrangMuc           string                 `json:"tinhTrangMuc"`
	DienTichPhuMuc         string                 `json:"dienTichPhuMuc"`
	DungMoi                string                 `json:"dungMoi"`
	NhietDoSay             string                 `json:"nhietDoSay"`
	ThoiGianSay            string                 `json:"thoiGianSay"`
	GhiChu                 string                 `json:"ghiChu"`
	MaKhung                string                 `json:"maKhung"`
	ThongSoLua             string                 `json:"thongSoLua"`
	KhoangCachKhungInBanIn string                 `json:"khoangCachKhungInBanIn"`
	CachIn                 string                 `json:"cachIn"`
	CungDao                string                 `json:"cungDao"`
	DoBenDao               string                 `json:"doBenDao"`
	DoNghiengDao           string                 `json:"doNghiengDao"`
	TocDoDao               string                 `json:"tocDoDao"`
	TocDo                  string                 `json:"tocDo"`
	Description            string                 `json:"description"`
	DeviceType             enum.DeviceConfigType  `json:"deviceType" binding:"required"`
	DeviceConfig           map[string]interface{} `json:"deviceConfig"`
	InkID                  string                 `json:"inkID"`
}

type CreateDeviceConfigResponse struct {
	ID string `json:"id"`
}

type EditDeviceConfigRequest struct {
	ID                     string                 `json:"id" binding:"required"`
	ProductionOrderID      string                 `json:"productionOrderID" binding:"required"`
	ProductionPlanID       string                 `json:"productionPlanID"`
	DeviceID               string                 `json:"deviceID" binding:"required"`
	Color                  string                 `json:"color"`
	MaThongSoMay           string                 `json:"maThongSoMay"`
	MaTaiLieuHuongDan      string                 `json:"maTaiLieuHuongDan"`
	NgayHieuLuc            string                 `json:"ngayHieuLuc"`
	StageID                string                 `json:"stageID"`
	MaSanPham              string                 `json:"maSanPham"`
	MaSanPhamNoiBo         string                 `json:"maSanPhamNoiBo"`
	ChuKyIn                string                 `json:"chuKyIn"`
	ThoiGianChuanBi        string                 `json:"thoiGianChuanBi"`
	TenMauMuc              string                 `json:"tenMauMuc"`
	TenLoaiMuc             string                 `json:"tenLoaiMuc"`
	SoThuTuIn              string                 `json:"soThuTuIn"`
	MaPhim                 string                 `json:"maPhim"`
	MaMauMuc               string                 `json:"maMauMuc"`
	TinhTrangMuc           string                 `json:"tinhTrangMuc"`
	DienTichPhuMuc         string                 `json:"dienTichPhuMuc"`
	DungMoi                string                 `json:"dungMoi"`
	NhietDoSay             string                 `json:"nhietDoSay"`
	ThoiGianSay            string                 `json:"thoiGianSay"`
	GhiChu                 string                 `json:"ghiChu"`
	MaKhung                string                 `json:"maKhung"`
	ThongSoLua             string                 `json:"thongSoLua"`
	KhoangCachKhungInBanIn string                 `json:"khoangCachKhungInBanIn"`
	CachIn                 string                 `json:"cachIn"`
	CungDao                string                 `json:"cungDao"`
	DoBenDao               string                 `json:"doBenDao"`
	DoNghiengDao           string                 `json:"doNghiengDao"`
	TocDoDao               string                 `json:"tocDoDao"`
	TocDo                  string                 `json:"tocDo"`
	Description            string                 `json:"description"`
	DeviceConfig           map[string]interface{} `json:"deviceConfig"` // todo do we need this?
	InkID                  string                 `json:"inkID"`
}

type EditDeviceConfigResponse struct {
}

type DeleteDeviceConfigRequest struct {
	ID string `json:"id"`
}

type DeleteDeviceConfigResponse struct {
}
