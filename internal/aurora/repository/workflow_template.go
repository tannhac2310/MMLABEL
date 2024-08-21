package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type WorkflowTemplateRepo interface {
	Insert(ctx context.Context, e *model.WorkflowTemplate) error
	Update(ctx context.Context, e *model.WorkflowTemplate) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.WorkflowTemplate, error)
	Search(ctx context.Context, s *SearchWorkflowTemplateOpts) ([]*WorkflowTemplateData, error)
	Count(ctx context.Context, s *SearchWorkflowTemplateOpts) (*CountResult, error)
}

type sWorkflowTemplateRepo struct {
}

func NewWorkflowTemplateRepo() WorkflowTemplateRepo {
	return &sWorkflowTemplateRepo{}
}

func (r *sWorkflowTemplateRepo) Insert(ctx context.Context, e *model.WorkflowTemplate) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sWorkflowTemplateRepo) FindByID(ctx context.Context, id string) (*model.WorkflowTemplate, error) {
	e := &model.WorkflowTemplate{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sWorkflowTemplateRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sWorkflowTemplateRepo) Update(ctx context.Context, e *model.WorkflowTemplate) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sWorkflowTemplateRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE workflow_templates SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("workflow_templates cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sWorkflowTemplateRepo not found any records to delete")
	}

	return nil
}

// SearchWorkflowTemplateOpts all params is options
type SearchWorkflowTemplateOpts struct {
	IDs    []string
	Name   string
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchWorkflowTemplateOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.WorkflowTemplateFieldID)
	}
	// todo add more search options example:
	if s.Name != "" {
		args = append(args, "%"+s.Name+"%")
		conds += fmt.Sprintf(" AND b.%[2]s ILIKE $%[1]d",
			len(args), model.WorkflowTemplateFieldName)

	}

	b := &model.WorkflowTemplate{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf("SELECT count(*) as cnt FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL", b.TableName(), joins, conds), args
	}

	order := " ORDER BY b.id DESC "
	if s.Sort != nil {
		order = fmt.Sprintf(" ORDER BY b.%s %s", s.Sort.By, s.Sort.Order)
	}
	return fmt.Sprintf("SELECT b.%s FROM %s AS b %s WHERE TRUE %s AND b.deleted_at IS NULL %s LIMIT %d OFFSET %d", strings.Join(fields, ", b."), b.TableName(), joins, conds, order, s.Limit, s.Offset), args
}

type WorkflowTemplateData struct {
	*model.WorkflowTemplate
}

func (r *sWorkflowTemplateRepo) Search(ctx context.Context, s *SearchWorkflowTemplateOpts) ([]*WorkflowTemplateData, error) {
	WorkflowTemplate := make([]*WorkflowTemplateData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&WorkflowTemplate)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return WorkflowTemplate, nil
}

func (r *sWorkflowTemplateRepo) Count(ctx context.Context, s *SearchWorkflowTemplateOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sWorkflowTemplateRepo.Count: %w", err)
	}

	return countResult, nil
}
