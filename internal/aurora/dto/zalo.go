package dto

import (
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
)

type ZaloAuthLinkGeneratorRequest struct {
	AppID     string `json:"appId" binding:"required"`
	SecretKey string `json:"secretKey" binding:"required"`
	Redirect  string `json:"redirect" binding:"required"`
}
type ZaloAuthLinkGeneratorResponse struct {
	Link  string `json:"link"`
	State string `json:"state"`
}

type ZaloAuthCallbackRequest struct {
	State         string `form:"state" binding:"required"`
	CodeChallenge string `form:"code_challenge" binding:"required"`
	OaID          string `form:"oa_id" binding:"required"`
	Code          string `form:"code" binding:"required"`
}
type ZaloAuthCallbackResponse struct {
}

type ZaloRenewAccessTokenRequest struct {
	AppID string `json:"appId" binding:"required"`
	OaID  string `json:"oaId" binding:"required"`
}
type ZaloRenewAccessTokenResponse struct{}

type DeleteZaloOARequest struct {
	ID string `json:"id"`
}
type DeleteZaloResponse struct {
}
type OaFilter struct {
	IDs    []string `json:"ids"`
	AppID  string   `json:"appId"`
	OaID   string   `json:"oaId"`
	Search string   `json:"search"`
}

type FindZaloOaRequest struct {
	Filter *OaFilter         `json:"filter"`
	Paging *commondto.Paging `json:"paging" binding:"required"`
}
type OaData struct {
	ID            string    `json:"id"`
	AppID         string    `json:"appId"`
	OaID          string    `json:"oaId"`
	OaName        string    `json:"oaName"`
	CreatedBy     string    `json:"createdBy"`
	UpdatedBy     string    `json:"updatedBy"`
	CreatedByName string    `json:"createdByName"`
	UpdatedByName string    `json:"updatedByName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type FindZaloOaResponse struct {
	Oas      []*OaData         `json:"oas"`
	NextPage *commondto.Paging `json:"nextPage"`
	Total    int64             `json:"total"`
}
