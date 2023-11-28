package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"time"
)

type StageFilter struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type FindStagesRequest struct {
	Filter *StageFilter      `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindStagesResponse struct {
	Stages []*Stage `json:"stages"`
	Total  int64    `json:"total"`
}
type Stage struct {
	ID             string                 `json:"id"`
	ParentID       string                 `json:"parentID"`
	DepartmentCode string                 `json:"departmentCode"`
	Name           string                 `json:"name"`
	ShortName      string                 `json:"shortName"`
	Code           string                 `json:"code"`
	Sorting        int16                  `json:"sorting"`
	ErrorCode      string                 `json:"errorCode"`
	Data           map[string]interface{} `json:"data"`
	Note           string                 `json:"note"`
	Status         enum.StageStatus       `json:"status"`
	CreatedBy      string                 `json:"createdBy"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}

type CreateStageRequest struct {
	ParentID       string                 `json:"parentID"`
	DepartmentCode string                 `json:"departmentCode"`
	Name           string                 `json:"name" binding:"required"`
	ShortName      string                 `json:"shortName" binding:"required"`
	Code           string                 `json:"code"`
	Sorting        int16                  `json:"sorting"`
	ErrorCode      string                 `json:"errorCode"`
	Data           map[string]interface{} `json:"data"`
	Note           string                 `json:"note"`
	Status         enum.StageStatus       `json:"status" binding:"required"`
}

type CreateStageResponse struct {
	ID string `json:"id"`
}

type EditStageRequest struct {
	ID             string                 `json:"id" binding:"required"`
	ParentID       string                 `json:"parentID"`
	DepartmentCode string                 `json:"departmentCode"`
	Name           string                 `json:"name" binding:"required"`
	ShortName      string                 `json:"shortName" binding:"required"`
	Code           string                 `json:"code"`
	Sorting        int16                  `json:"sorting"`
	ErrorCode      string                 `json:"errorCode"`
	Data           map[string]interface{} `json:"data"`
	Note           string                 `json:"note"`
	Status         enum.StageStatus       `json:"status" binding:"required"`
}

type EditStageResponse struct {
}

type DeleteStageRequest struct {
	ID string `json:"id"`
}

type DeleteStageResponse struct {
}
