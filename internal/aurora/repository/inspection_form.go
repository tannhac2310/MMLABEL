package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type InspectionFormRepo interface {
	Insert(ctx context.Context, e *model.InspectionForm) error
	Update(ctx context.Context, e *model.InspectionForm) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*InspectionFormData, error)
	Search(ctx context.Context, s *SearchInspectionFormOpts) ([]*InspectionFormData, error)
	Count(ctx context.Context, s *SearchInspectionFormOpts) (*CountResult, error)
}

type sInspectionFormRepo struct {
}

func NewInspectionFormRepo() InspectionFormRepo {
	return &sInspectionFormRepo{}
}

func (r *sInspectionFormRepo) Insert(ctx context.Context, e *model.InspectionForm) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sInspectionFormRepo) FindByID(ctx context.Context, id string) (*InspectionFormData, error) {
	e := &InspectionFormData{
		InspectionForm: &model.InspectionForm{},
	}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sInspectionFormRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sInspectionFormRepo) Update(ctx context.Context, e *model.InspectionForm) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sInspectionFormRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE inspection_forms SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("inspection_forms cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sInspectionFormRepo not found any records to delete")
	}

	return nil
}

// SearchInspectionFormOpts all params is options
type SearchInspectionFormOpts struct {
	IDs               []string
	ProductionOrderID string
	DefectType        []string
	CreatedAtFrom     time.Time
	CreatedAtTo       time.Time
	Limit             int64
	Offset            int64
	Sort              *Sort
}

func (s *SearchInspectionFormOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := " JOIN production_orders AS po ON b.production_order_id = po.id "
	//joins += " JOIN products AS p ON po.product_id = p.id "
	//if s.ProductSearch != "" {
	//	args = append(args, "%"+s.ProductSearch+"%")
	//	conds += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.code ILIKE $%d)", len(args), len(args))
	//}

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.InspectionFormFieldID)
	}
	if s.ProductionOrderID != "" {
		args = append(args, s.ProductionOrderID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.InspectionFormFieldProductionOrderID, len(args))
	}
	if len(s.DefectType) > 0 {
		args = append(args, s.DefectType)
		joins += " JOIN inspection_errors AS ie ON b.id = ie.inspection_form_id "
		conds += fmt.Sprintf(" AND ie.%s = ANY($%d)", model.InspectionErrorFieldErrorType, len(args))
	}

	b := &model.InspectionForm{}
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

type InspectionFormData struct {
	*model.InspectionForm
}

func (r *sInspectionFormRepo) Search(ctx context.Context, s *SearchInspectionFormOpts) ([]*InspectionFormData, error) {
	InspectionForm := make([]*InspectionFormData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&InspectionForm)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return InspectionForm, nil
}

func (r *sInspectionFormRepo) Count(ctx context.Context, s *SearchInspectionFormOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sInspectionFormRepo.Count: %w", err)
	}

	return countResult, nil
}
