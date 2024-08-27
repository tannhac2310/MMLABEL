package repository

import (
	"context"
	"fmt"
	"strings"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type CustomFieldRepo interface {
	Insert(ctx context.Context, e *model.CustomField) error
	Update(ctx context.Context, e *model.CustomField) error
	Delete(ctx context.Context, id string) error
	DeleteByEntity(ctx context.Context, entityType enum.CustomFieldType, entityId string) error
	Search(ctx context.Context, s *SearchCustomFieldsOpts) ([]*CustomFieldData, error)
	Count(ctx context.Context, s *SearchCustomFieldsOpts) (*CountResult, error)
}

type customFieldsRepo struct {
}

func NewCustomFieldRepo() CustomFieldRepo {
	return &customFieldsRepo{}
}

func (r *customFieldsRepo) Insert(ctx context.Context, e *model.CustomField) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *customFieldsRepo) Update(ctx context.Context, e *model.CustomField) error {
	return cockroach.Update(ctx, e)
}

func (r *customFieldsRepo) Delete(ctx context.Context, id string) error {
	//sql := "UPDATE custom_fields SET deleted_at = NOW() WHERE id = $1;"
	sql := "DELETE FROM custom_fields WHERE id = $1;" // should be hard delete
	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (r *customFieldsRepo) DeleteByEntity(ctx context.Context, entityType enum.CustomFieldType, entityId string) error {
	sql := `UPDATE custom_fields
		SET deleted_at = NOW()
		WHERE entity_type = $1 AND entity_id = $2`

	cmd, err := cockroach.Exec(ctx, sql, entityType, entityId)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

// SearchCustomFieldsOpts all params is options
type SearchCustomFieldsOpts struct {
	EntityId   string
	EntityType enum.CustomFieldType
	Field      string
	Code       string
	Limit      int64
	Offset     int64
	Sort       *Sort
}

func (s *SearchCustomFieldsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if s.Field != "" {
		args = append(args, s.Field)
		conds += fmt.Sprintf(" AND b.%[2]s = $%[1]d",
			len(args), model.CustomFieldFieldField)
	}
	if s.EntityId != "" {
		args = append(args, s.EntityId)
		conds += fmt.Sprintf(" AND b.%[2]s = $%[1]d",
			len(args), model.CustomFieldFieldEntityID)
	}
	if s.EntityType > 0 {
		args = append(args, s.EntityType)
		conds += fmt.Sprintf(" AND b.%[2]s = $%[1]d",
			len(args), model.CustomFieldFieldEntityType)
	}

	b := &model.CustomField{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s`, b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf(`SELECT b.%s
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL
		%s
		LIMIT %d
		OFFSET %d`, strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type CustomFieldData struct {
	*model.CustomField
}

func (r *customFieldsRepo) Search(ctx context.Context, s *SearchCustomFieldsOpts) ([]*CustomFieldData, error) {
	message := make([]*CustomFieldData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}

func (r *customFieldsRepo) Count(ctx context.Context, s *SearchCustomFieldsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("chat.Count: %w", err)
	}

	return countResult, nil
}
