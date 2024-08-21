package dto

import (
	"mmlabel.gitlab.com/mm-printing-backend/pkg/commondto"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type WorkflowTemplate struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Config    any               `json:"config"`
	Status    enum.CommonStatus `json:"status"`
	CreatedBy string            `json:"created_by"`
	CreatedAt string            `json:"created_at"`
	UpdatedBy string            `json:"updated_by"`
	UpdatedAt string            `json:"updated_at"`
}

type CreateWorkflowTemplateRequest struct {
	Name       string            `json:"name"`
	ConfigData any               `json:"configData"`
	Status     enum.CommonStatus `json:"status"`
}

type CreateWorkflowTemplateResponse struct {
	ID string `json:"id"`
}
type EditWorkflowTemplateRequest struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	ConfigData any               `json:"configData"`
	Status     enum.CommonStatus `json:"status"`
}
type EditWorkflowTemplateResponse struct{}

type WorkflowTemplatesFilter struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type FindWorkflowTemplatesRequest struct {
	Filter *WorkflowTemplatesFilter `json:"filter" binding:"required"`
	Paging *commondto.Paging        `json:"paging" binding:"required"`
}

type FindWorkflowTemplatesResponse struct {
	WorkflowTemplates []*WorkflowTemplate `json:"workflowTemplates"`
	Total             int64               `json:"total"`
}

type DeleteWorkflowTemplateRequest struct {
	ID string `json:"id"`
}
type DeleteWorkflowTemplateResponse struct{}
