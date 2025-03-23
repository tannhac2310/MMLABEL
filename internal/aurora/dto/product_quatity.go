package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
)

type ProductQualityFilter struct {
	IDs               []string `json:"ids"`
	ProductionOrderID string   `json:"productionOrderID"`
	DeviceIDs         []string `json:"deviceIDs"`
	DefectTypes       []string `json:"defectType"`
	//ProductSearch     string    `json:"productSearch"`
	CreatedAtFrom time.Time `json:"createdAtFrom"`
	CreatedAtTo   time.Time `json:"createdAtTo"`
}

type FindProductQualityRequest struct {
	Filter *ProductQualityFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging     `json:"paging" binding:"required"`
}

type FindProductQualityResponse struct {
	ProductQuality []*ProductQuality         `json:"productQuality"`
	Total          int64                     `json:"total"`
	Analysis       []*ProductQualityAnalysis `json:"analysis"`
}
type ProductQualityAnalysis struct {
	DefectType string `json:"defectType"`
	Count      int64  `json:"count"`
}
type DeviceData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductQuality struct {
	ID                       string    `json:"id"`
	ProductionOrderID        string    `json:"productionOrderID"`
	ProductionOrderCode      string    `json:"productionOrderCode"`
	ProductionOrderName      string    `json:"productionOrderName"`
	ProductionOrderStartDate time.Time `json:"productionOrderStartDate"`
	InspectionDate           time.Time `json:"inspectionDate"`
	InspectorName            string    `json:"inspectorName"`
	Quantity                 int64     `json:"quantity"`
	ProductID                string    `json:"productID"`
	ProductName              string    `json:"productName"`
	ProductCode              string    `json:"productCode"`
	CustomerID               string    `json:"customerID"`
	CustomerName             string    `json:"customerName"`
	CustomerCode             string    `json:"customerCode"`
	SoLuongHopDong           int64     `json:"soLuongHopDong"`
	SoLuongIn                int64     `json:"soLuongIn"`
	MaDonDatHang             string    `json:"maDonDatHang"`
	NguoiKiemTra             string    `json:"nguoiKiemTra"`
	NguoiPheDuyet            string    `json:"nguoiPheDuyet"`
	SoLuongThanhPhamDat      int64     `json:"soLuongThanhPhamDat"`
	OrderData                struct {
		ID          string `json:"id"`
		MaDatHangMm string `json:"maDatHangMm"`
		Status      string `json:"status"`
	} `json:"orderData"`
	Note             string            `json:"note"`
	InspectionErrors []InspectionError `json:"inspectionErrors"`
	CreatedBy        string            `json:"createdBy"`
	UpdatedBy        string            `json:"updatedBy"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
}

type InspectionError struct {
	ID               string `json:"id"`
	DeviceID         string `json:"deviceID" binding:"required"`
	DeviceName       string `json:"deviceName" binding:"required"`
	InspectionFormID string `json:"inspectionFormID"`
	ErrorType        string `json:"errorType" binding:"required"`
	Quantity         int64  `json:"quantity" binding:"required"`
	Note             string `json:"note"`
	NhanVienThucHien string `json:"nhanVienThucHien"`
}

type CreateProductQualityRequest struct {
	ProductionOrderID string    `json:"productionOrderID" binding:"required"`
	InspectionDate    time.Time `json:"inspectionDate"  binding:"required"`
	InspectorName     string    `json:"inspectorName"`
	Quantity          int64     `json:"quantity" binding:"required"`
	ProductID         string    `json:"productID" binding:"required"`
	SoLuongHopDong    int64     `json:"soLuongHopDong"`
	SoLuongIn         int64     `json:"soLuongIn"`
	//MaDonDatHang        string            `json:"maDonDatHang"`
	NguoiKiemTra        string            `json:"nguoiKiemTra"`
	NguoiPheDuyet       string            `json:"nguoiPheDuyet"`
	SoLuongThanhPhamDat int64             `json:"soLuongThanhPhamDat" binding:"required"`
	Note                string            `json:"note"`
	InspectionErrors    []InspectionError `json:"inspectionErrors" binding:"required"`
}

type CreateProductQualityResponse struct {
	ID string `json:"id"`
}

type EditProductQualityRequest struct {
	ID                  string            `json:"id" binding:"required"`
	ProductionOrderID   string            `json:"productionOrderID" binding:"required"`
	InspectionDate      time.Time         `json:"inspectionDate"`
	InspectorName       string            `json:"inspectorName"`
	Quantity            int64             `json:"quantity" binding:"required"`
	ProductID           string            `json:"productID" binding:"required"`
	SoLuongHopDong      int64             `json:"soLuongHopDong" binding:"required"`
	SoLuongIn           int64             `json:"soLuongIn" binding:"required"`
	NguoiKiemTra        string            `json:"nguoiKiemTra"`
	NguoiPheDuyet       string            `json:"nguoiPheDuyet"`
	SoLuongThanhPhamDat int64             `json:"soLuongThanhPhamDat" binding:"required"`
	Note                string            `json:"note"`
	InspectionErrors    []InspectionError `json:"inspectionErrors" binding:"required"`
}

type EditProductQualityResponse struct {
}

type DeleteProductQualityRequest struct {
	ID string `json:"id"`
}

type DeleteProductQualityResponse struct {
}
