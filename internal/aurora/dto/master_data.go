package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type MasterDataUserField struct {
	ID           string `json:"id"`
	MasterDataID string `json:"masterDataID"`
	FieldName    string `json:"fieldName"`
	FieldValue   string `json:"fieldValue"`
}
type ShortProductionPlan struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	CustomData map[string]string `json:"customData"`
}
type MasterData struct {
	ID              string                 `json:"id"`
	Type            enum.MasterDataType    `json:"type"`
	Name            string                 `json:"name"`
	Code            string                 `json:"code"`
	UserFields      []*MasterDataUserField `json:"userFields"`
	Description     string                 `json:"description"`
	Status          enum.MasterDataStatus  `json:"status"`
	ProductionPlans []*ShortProductionPlan `json:"productionPlans"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	CreatedBy       string                 `json:"createdBy"`
	UpdatedBy       string                 `json:"updatedBy"`
}

type CreateMasterDataUserField struct {
	ID           string `json:"id"`
	MasterDataID string `json:"masterDataID"`
	FieldName    string `json:"fieldName"`
	FieldValue   string `json:"fieldValue"`
}

// CreateMasterDataRequest create
type CreateMasterDataRequest struct {
	Type        enum.MasterDataType          `json:"type" binding:"required"`
	Name        string                       `json:"name" binding:"required"`
	Description string                       `json:"description"`
	UserFields  []*CreateMasterDataUserField `json:"userFields" binding:"required"`
	Status      enum.MasterDataStatus        `json:"status" binding:"required"`
	Code        string                       `json:"code" binding:"required"`
}

type CreateMasterDataResponse struct {
	ID string `json:"id"`
}

// UpdateMasterDataRequest update
type UpdateMasterDataRequest struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Status      enum.MasterDataStatus        `json:"status"`
	UserFields  []*CreateMasterDataUserField `json:"userFields"`
	Code        string                       `json:"code"`
}

type UpdateMasterDataResponse struct {
}

// DeleteMasterDataRequest delete
type DeleteMasterDataRequest struct {
	ID string `json:"id"`
}

type DeleteMasterDataResponse struct {
}

// SearchMasterDataFilter get
type SearchMasterDataFilter struct {
	IDs    []string            `json:"ids"`
	Type   enum.MasterDataType `json:"type"`
	Search string              `json:"search"`
}

type SearchMasterDataRequest struct {
	Filter SearchMasterDataFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging      `json:"paging" binding:"required"`
}

type GetMasterDataResponse struct {
	MasterData []*MasterData `json:"masterData"`
	Total      int64         `json:"total"`
}
