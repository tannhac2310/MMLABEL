package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MThongSoMayInRepo interface {
	Insert(ctx context.Context, e *model.MThongSoMayIn) error
	Update(ctx context.Context, e *model.MThongSoMayIn) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MThongSoMayInData, error)
	Search(ctx context.Context, s *SearchMThongSoMayInOpts) ([]*MThongSoMayInData, error)
	Count(ctx context.Context, s *SearchMThongSoMayInOpts) (*CountResult, error)
}

type sMThongSoMayInRepo struct {
}

func NewMThongSoMayInRepo() MThongSoMayInRepo {
	return &sMThongSoMayInRepo{}
}

func (r *sMThongSoMayInRepo) Insert(ctx context.Context, e *model.MThongSoMayIn) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMThongSoMayInRepo) FindByID(ctx context.Context, id string) (*MThongSoMayInData, error) {
	e := &MThongSoMayInData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMThongSoMayInRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMThongSoMayInRepo) Update(ctx context.Context, e *model.MThongSoMayIn) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMThongSoMayInRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE m_thong_so_may_in SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("m_thong_so_may_in cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMThongSoMayInRepo not found any records to delete")
	}

	return nil
}

// SearchMThongSoMayInOpts all params is options
type SearchMThongSoMayInOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMThongSoMayInOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MThongSoMayInFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.MThongSoMayInFieldName, model.MThongSoMayInFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.MThongSoMayInFieldCode, len(args))
	//}

	b := &model.MThongSoMayIn{}
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

type MThongSoMayInData struct {
	*model.MThongSoMayIn
}

func (r *sMThongSoMayInRepo) Search(ctx context.Context, s *SearchMThongSoMayInOpts) ([]*MThongSoMayInData, error) {
	MThongSoMayIn := make([]*MThongSoMayInData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MThongSoMayIn)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MThongSoMayIn, nil
}

func (r *sMThongSoMayInRepo) Count(ctx context.Context, s *SearchMThongSoMayInOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMThongSoMayInRepo.Count: %w", err)
	}

	return countResult, nil
}
