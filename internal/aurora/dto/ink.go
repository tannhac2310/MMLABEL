package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type Ink struct {
	ID                 string                 `json:"id"`
	ImportID           string                 `json:"importID"`
	Name               string                 `json:"name"`
	Code               string                 `json:"code"`
	MixInk             *MixInk                `json:"mixInk"`
	ProductCodes       []string               `json:"productCodes"`
	Position           string                 `json:"position"`
	Location           string                 `json:"location"`
	Manufacturer       string                 `json:"manufacturer"`
	ColorDetail        map[string]interface{} `json:"colorDetail"`
	Quantity           float64                `json:"quantity"`
	ExpirationDate     string                 `json:"expirationDate"`
	Description        string                 `json:"description"`
	Data               map[string]interface{} `json:"data"`
	Status             enum.CommonStatus      `json:"status"`
	ProductionPlanIDs  []string               `json:"productionPlanIDs"`
	ProductionOrderIDs []string               `json:"productionOrderIDs"`
	CreatedBy          string                 `json:"createdBy"`
	UpdatedBy          string                 `json:"updatedBy"`
	CreatedAt          time.Time              `json:"createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt"`
}

type InkFilter struct {
	Name   string            `json:"name"`
	ID     string            `json:"id"`
	NotIDs []string          `json:"notIDs"`
	Code   string            `json:"code"`
	Status enum.CommonStatus `json:"status"`
}
type FindInkRequest struct {
	Filter *InkFilter        `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}
type FindInksResponse struct {
	Ink   []*Ink `json:"ink"`
	Total int64  `json:"total"`
}

// dto for ink.edit

type EditInkRequest struct {
	ID             string                 `json:"id" binding:"required"`
	Name           string                 `json:"name"`
	Code           string                 `json:"code"`
	ProductCodes   []string               `json:"productCodes"`
	Position       string                 `json:"position"`
	Location       string                 `json:"location"`
	Manufacturer   string                 `json:"manufacturer"`
	ColorDetail    map[string]interface{} `json:"colorDetail"`
	ExpirationDate string                 `json:"expirationDate"`
	Description    string                 `json:"description"`
	Data           map[string]interface{} `json:"data"`
	Status         enum.CommonStatus      `json:"status"`
	Quantity       float64                `json:"quantity"`
}

type EditInkResponse struct{}

type CreateInkImportRequest struct {
	Name            string                   `json:"name" binding:"required"`
	Code            string                   `json:"code" binding:"required"`
	Description     string                   `json:"description"`
	Data            map[string]interface{}   `json:"data"`
	InkImportDetail []*CreateInkImportDetail `json:"inkImportDetail" binding:"required"`
}

type CreateInkImportDetail struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name" binding:"required"`
	Code           string                 `json:"code" binding:"required"`
	ProductCodes   []string               `json:"productCodes"`
	Position       string                 `json:"position"`
	Location       string                 `json:"location"`
	Manufacturer   string                 `json:"manufacturer"`
	ColorDetail    map[string]interface{} `json:"colorDetail"`
	Quantity       float64                `json:"quantity"`
	ExpirationDate string                 `json:"expirationDate"` // DD-MM-YYYY
	Description    string                 `json:"description"`
	Data           map[string]interface{} `json:"data"`
}

type InkImportDetail struct {
	ID             string                 `json:"id"`
	InkID          string                 `json:"inkID"`
	Name           string                 `json:"name" binding:"required"`
	Code           string                 `json:"code" binding:"required"`
	ProductCodes   []string               `json:"productCodes"`
	Position       string                 `json:"position"`
	Location       string                 `json:"location"`
	Manufacturer   string                 `json:"manufacturer"`
	ColorDetail    map[string]interface{} `json:"colorDetail"`
	Quantity       float64                `json:"quantity"`
	ExpirationDate string                 `json:"expirationDate"` // DD-MM-YYYY
	Description    string                 `json:"description"`
	Data           map[string]interface{} `json:"data"`
}

type CreateInkImportResponse struct {
	ID string `json:"id"`
}

// dto for ink_import.find
type FindInkImportsRequest struct {
	Filter *InkImportFilter  `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type InkImportFilter struct {
	Name   string                     `json:"name"`
	Code   string                     `json:"code"`
	ID     string                     `json:"id"`
	Status enum.InventoryCommonStatus `json:"status"`
}

type InkImport struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	Code            string                     `json:"code"`
	Description     string                     `json:"description"`
	Data            map[string]interface{}     `json:"data"`
	Status          enum.InventoryCommonStatus `json:"status"`
	InkImportDetail []*InkImportDetail         `json:"inkImportDetail"`
	CreatedBy       string                     `json:"createdBy"`
	CreatedAt       time.Time                  `json:"createdAt"`
}

type FindInkImportsResponse struct {
	InkImport []*InkImport `json:"inkImport"`
	Total     int64        `json:"total"`
}

// dto for in_export.edit
type EditInkExportDetail struct {
	InkID       string                 `json:"inkID" binding:"required"`
	Quantity    float64                `json:"quantity" binding:"required"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

type EditInkExportRequest struct {
	ID              string                 `json:"id" binding:"required"`
	Description     string                 `json:"description"`
	InkExportDetail []*EditInkExportDetail `json:"inkExportDetail" binding:"required"`
}

type EditInkExportResponse struct{}

// dto for ink_export.create
type CreateInkExportRequest struct {
	Name              string                   `json:"name" binding:"required"`
	Code              string                   `json:"code" binding:"required"`
	ProductionOrderID string                   `json:"productionOrderID" binding:"required"`
	ExportDate        string                   `json:"exportDate"` // DD-MM-YYYY
	Description       string                   `json:"description"`
	Data              map[string]interface{}   `json:"data"`
	InkExportDetail   []*CreateInkExportDetail `json:"inkExportDetail" binding:"required"`
}

type CreateInkExportDetail struct {
	InkID       string                 `json:"inkID" binding:"required"`
	Quantity    float64                `json:"quantity" binding:"required"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

type InkDataExportDetail struct {
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	ProductCodes []string               `json:"productCodes"`
	Quantity     float64                `json:"quantity"`
	Position     string                 `json:"position"`
	Location     string                 `json:"location"`
	Manufacturer string                 `json:"manufacturer"`
	ColorDetail  map[string]interface{} `json:"colorDetail"`
}
type InkExportDetail struct {
	ID          string                 `json:"id"`
	InkID       string                 `json:"inkID"`
	InkExportID string                 `json:"inkExportID"`
	InkData     *InkDataExportDetail   `json:"inkData"`
	Quantity    float64                `json:"quantity"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

type CreateInkExportResponse struct {
	ID string `json:"id"`
}

// dto for ink_export.find
type FindInkExportsRequest struct {
	Filter *InkExportFilter  `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type InkExportFilter struct {
	ProductName string `json:"productName"`
	InkCode     string `json:"inkCode"`
	Search      string `json:"search"`
	ID          string `json:"id"`
}

type ProductionOrderData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProductCode string `json:"productCode"`
	ProductName string `json:"productName"`
}

type InkExport struct {
	ID                  string                     `json:"id"`
	Name                string                     `json:"name"`
	Code                string                     `json:"code"`
	ProductionOrderID   string                     `json:"productionOrderID"`
	ProductionOrderData *ProductionOrderData       `json:"productionOrderData"`
	Description         string                     `json:"description"`
	Data                map[string]interface{}     `json:"data"`
	Status              enum.InventoryCommonStatus `json:"status"`
	CreatedBy           string                     `json:"createdBy"`
	CreatedAt           time.Time                  `json:"createdAt"`
	UpdatedBy           string                     `json:"updatedBy"`
	UpdatedAt           time.Time                  `json:"updatedAt"`
	CreatedByName       string                     `json:"createdByName"`
	UpdatedByName       string                     `json:"updatedByName"`
	InkExportDetail     []*InkExportDetail         `json:"inkExportDetail"`
	InkReturnData       []*InkReturn               `json:"inkReturnData"`
}

type FindInkExportsResponse struct {
	InkExport []*InkExport `json:"inkExport"`
	Total     int64        `json:"total"`
}

type CreateInkReturnRequest struct {
	Name            string                       `json:"name"`
	Code            string                       `json:"code"`
	Description     string                       `json:"description"`
	InkExportID     string                       `json:"inkExportID"`
	Data            map[string]interface{}       `json:"data"`
	InkReturnDetail []*CreateInkReturnDetailOpts `json:"inkReturnDetail" binding:"required"`
}

// dto for ink_return.edit
type EditInkReturnDetailOpts struct {
	InkID             string                 `json:"inkID"`
	InkExportDetailID string                 `json:"inkExportDetailID"`
	Quantity          float64                `json:"quantity"`
	ColorDetail       map[string]interface{} `json:"colorDetail"`
	Description       string                 `json:"description"`
	Data              map[string]interface{} `json:"data"`
}

type EditInkReturnRequest struct {
	ID              string                     `json:"id" binding:"required"`
	Description     string                     `json:"description"`
	InkReturnDetail []*EditInkReturnDetailOpts `json:"inkReturnDetail" binding:"required"`
}

type EditInkReturnResponse struct{}

type CreateInkReturnDetailOpts struct {
	InkID             string                 `json:"inkID"`
	InkExportDetailID string                 `json:"inkExportDetailID"`
	Quantity          float64                `json:"quantity"`
	ColorDetail       map[string]interface{} `json:"colorDetail"`
	Description       string                 `json:"description"`
	Data              map[string]interface{} `json:"data"`
}

type CreateInkReturnResponse struct {
}

type InkReturnFilter struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type FindInkReturnsRequest struct {
	Filter *InkReturnFilter  `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindInkReturnsResponse struct {
	InkReturn []*InkReturn `json:"inkReturn"`
	Total     int64        `json:"total"`
}

type InkReturn struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	Code            string                     `json:"code"`
	InkExportID     string                     `json:"inkExportID"`
	Description     string                     `json:"description"`
	Data            map[string]interface{}     `json:"data"`
	Status          enum.InventoryCommonStatus `json:"status"`
	InkReturnDetail []*InkReturnDetail         `json:"inkReturnDetail" binding:"required"`
	CreatedBy       string                     `json:"createdBy"`
	CreatedAt       time.Time                  `json:"createdAt"`
	UpdatedBy       string                     `json:"updatedBy"`
	UpdatedAt       time.Time                  `json:"updatedAt"`
	CreatedByName   string                     `json:"createdByName"`
	UpdatedByName   string                     `json:"updatedByName"`
}

type InkReturnDetail struct {
	ID                string                 `json:"id"`
	InkReturnID       string                 `json:"inkReturnID"`
	InkID             string                 `json:"inkID"`
	InkData           *InkDataExportDetail   `json:"inkData"`
	InkExportDetailID string                 `json:"inkExportDetailID"`
	Quantity          float64                `json:"quantity"`
	ColorDetail       map[string]interface{} `json:"colorDetail"`
	Description       string                 `json:"description"`
	Data              map[string]interface{} `json:"data"`
}

// dto for ink_export.find_by_po
type FindInkExportByPORequest struct {
	ProductionOrderID string `json:"productionOrderID" binding:"required"`
}

type FindInkExportByPOResponse struct {
	InkExportDetail []*InkExportDetail `json:"inkExportDetail"`
}

type CreateInkMixingFormulation struct {
	InkID       string  `json:"inkID"`
	Quantity    float64 `json:"quantity"`
	Description string  `json:"description"`
}

// create mix-ink
type CreateInkMixingRequest struct {
	Name           string                       `json:"name" binding:"required"`
	Code           string                       `json:"inkCode" binding:"required"`
	ProductCodes   []string                     `json:"productCodes"`
	Quantity       float64                      `json:"quantity" binding:"required"`
	ExpirationDate string                       `json:"expirationDate"`
	Position       string                       `json:"position"`
	Location       string                       `json:"location"`
	Description    string                       `json:"description"`
	InkFormulation []CreateInkMixingFormulation `json:"inkFormulation" binding:"required"`
}

type CreateInkMixingResponse struct {
	ID string `json:"id"`
}

type EditMixInkResponse struct{}

// find mix-ink
type FindInkMixingRequest struct {
	Filter *MixInkFilter     `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type MixInkFilter struct {
	Search string   `json:"search"`
	IDs    []string `json:"ids"`
	InkID  string   `json:"inkID"`
}

type MixInk struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Code           string                 `json:"code"`
	InkID          string                 `json:"inkID"`
	Quantity       float64                `json:"quantity"`
	ExpirationDate string                 `json:"expirationDate"`
	ProductCodes   []string               `json:"productCodes"`
	Position       string                 `json:"position"`
	Location       string                 `json:"location"`
	Description    string                 `json:"description"`
	CreatedBy      string                 `json:"createdBy"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedBy      string                 `json:"updatedBy"`
	UpdatedAt      time.Time              `json:"updatedAt"`
	CreatedByName  string                 `json:"createdByName"`
	UpdatedByName  string                 `json:"updatedByName"`
	InkFormulation []InkMixingFormulation `json:"inkFormulation"`
	Status         enum.CommonStatus      `json:"status"`
}

type FindInkMixingResponse struct {
	MixInk []*MixInk `json:"mixInk"`
	Total  int64     `json:"total"`
}

type InkMixingFormulation struct {
	ID          string  `json:"id"`
	InkID       string  `json:"inkID"`
	Quantity    float64 `json:"quantity"`
	Description string  `json:"description"`
	InkName     string  `json:"inkName"`
	InkCode     string  `json:"inkCode"`
}

type EditInkMixingRequest struct {
	ID             string                 `json:"id" binding:"required"`
	Name           string                 `json:"name" binding:"required"`
	Code           string                 `json:"inkCode" binding:"required"`
	ProductCodes   []string               `json:"productCodes"`
	Quantity       float64                `json:"quantity" binding:"required"`
	ExpirationDate string                 `json:"expirationDate"`
	Position       string                 `json:"position"`
	Location       string                 `json:"location"`
	Description    string                 `json:"description"`
	Manufacturer   string                 `json:"manufacturer"`
	InkFormulation []InkMixingFormulation `json:"inkFormulation" binding:"required"`
}

type EditInkMixingResponse struct{}
