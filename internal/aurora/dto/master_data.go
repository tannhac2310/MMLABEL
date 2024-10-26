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

type MasterData struct {
	ID          string                 `json:"id"`
	Type        enum.MasterDataType    `json:"type"`
	Name        string                 `json:"name"`
	UserFields  []*MasterDataUserField `json:"userFields"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	CreatedBy   string                 `json:"createdBy"`
	UpdatedBy   string                 `json:"updatedBy"`
}

// CreateMasterDataRequest create
type CreateMasterDataRequest struct {
	Type        enum.MasterDataType    `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	UserFields  []*MasterDataUserField `json:"userFields"`
}

type CreateMasterDataResponse struct {
	ID string `json:"id"`
}

// UpdateMasterDataRequest update
type UpdateMasterDataRequest struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	UserFields  []*MasterDataUserField `json:"userFields"`
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
	ID   string              `json:"id"`
	Type enum.MasterDataType `json:"type"`
}

type SearchMasterDataRequest struct {
	Filter SearchMasterDataFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging      `json:"paging" binding:"required"`
}

type GetMasterDataResponse struct {
	MasterData []*MasterData `json:"masterData"`
	Total      int64         `json:"total"`
}
