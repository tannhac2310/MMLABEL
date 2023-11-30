package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type OptionFilter struct {
	Name string `json:"name"`
	Entity string `json:"entity"`
}

type FindOptionRequest struct {
	Filter *OptionFilter     `json:"filter" binding:"required"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}

type FindOptionResponse struct {
	Options []*Option `json:"options"`
	Total   int64     `json:"total"`
}
type Option struct {
	ID        string                 `json:"id"`
	Entity    string                 `json:"entity"`
	Name      string                 `json:"name"`
	Code      string                 `json:"code"`
	OptionID  string                 `json:"optionID"`
	Data      map[string]interface{} `json:"data"`
	Status    enum.CommonStatus      `json:"status"`
	CreatedBy string                 `json:"createdBy"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type CreateOptionRequest struct {
	Name     string                 `json:"name"`
	Entity    string                `json:"entity"`
	Code     string                 `json:"code"`
	OptionID string                 `json:"optionID"`
	Data     map[string]interface{} `json:"data"`
	Status   enum.CommonStatus      `json:"status"`
}

type CreateOptionResponse struct {
	ID string `json:"id"`
}

type EditOptionRequest struct {
	ID       string                 `json:"id" binding:"required"`
	Name     string                 `json:"name"`
	Code     string                 `json:"code"`
	OptionID string                 `json:"optionID"`
	Data     map[string]interface{} `json:"data"`
	Status   enum.CommonStatus      `json:"status"`
}

type EditOptionResponse struct {
}

type DeleteOptionRequest struct {
	ID string `json:"id"`
}

type DeleteOptionResponse struct {
}
