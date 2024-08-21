package model

import (
	"database/sql"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

const (
	WorkflowTemplateFieldID        = "id"
	WorkflowTemplateFieldName      = "name"
	WorkflowTemplateFieldConfig    = "config"
	WorkflowTemplateFieldStatus    = "status"
	WorkflowTemplateFieldCreatedBy = "created_by"
	WorkflowTemplateFieldCreatedAt = "created_at"
	WorkflowTemplateFieldUpdatedAt = "updated_at"
	WorkflowTemplateUpdatedBy      = "updated_by"
	WorkflowTemplateFieldDeletedAt = "deleted_at"
)

type WorkflowTemplate struct {
	ID        string            `db:"id"`
	Name      string            `db:"name"`
	Config    any               `db:"config"`
	Status    enum.CommonStatus `db:"status"`
	CreatedBy string            `db:"created_by"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedBy string            `db:"updated_by"`
	UpdatedAt time.Time         `db:"updated_at"`
	DeletedAt sql.NullTime      `db:"deleted_at"`
}

func (rcv *WorkflowTemplate) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		WorkflowTemplateFieldID,
		WorkflowTemplateFieldName,
		WorkflowTemplateFieldConfig,
		WorkflowTemplateFieldStatus,
		WorkflowTemplateFieldCreatedBy,
		WorkflowTemplateFieldCreatedAt,
		WorkflowTemplateFieldUpdatedAt,
		WorkflowTemplateUpdatedBy,
		WorkflowTemplateFieldDeletedAt,
	}

	values = []interface{}{
		&rcv.ID,
		&rcv.Name,
		&rcv.Config,
		&rcv.Status,
		&rcv.CreatedBy,
		&rcv.CreatedAt,
		&rcv.UpdatedAt,
		&rcv.UpdatedBy,
		&rcv.DeletedAt,
	}

	return
}

func (*WorkflowTemplate) TableName() string {
	return "workflow_templates"
}
