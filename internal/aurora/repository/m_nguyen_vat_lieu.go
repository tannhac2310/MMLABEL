package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MNguyenVatLieuRepo interface {
	Insert(ctx context.Context, e *model.MNguyenVatLieu) error
	Update(ctx context.Context, e *model.MNguyenVatLieu) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MNguyenVatLieuData, error)
	Search(ctx context.Context, s *SearchMNguyenVatLieuOpts) ([]*MNguyenVatLieuData, error)
	Count(ctx context.Context, s *SearchMNguyenVatLieuOpts) (*CountResult, error)
}

type sMNguyenVatLieuRepo struct {
}

func NewMNguyenVatLieuRepo() MNguyenVatLieuRepo {
	return &sMNguyenVatLieuRepo{}
}

func (r *sMNguyenVatLieuRepo) Insert(ctx context.Context, e *model.MNguyenVatLieu) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMNguyenVatLieuRepo) FindByID(ctx context.Context, id string) (*MNguyenVatLieuData, error) {
	e := &MNguyenVatLieuData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMNguyenVatLieuRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMNguyenVatLieuRepo) Update(ctx context.Context, e *model.MNguyenVatLieu) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMNguyenVatLieuRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE m_nguyen_vat_lieu SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("m_nguyen_vat_lieu cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMNguyenVatLieuRepo not found any records to delete")
	}

	return nil
}

// SearchMNguyenVatLieuOpts all params is options
type SearchMNguyenVatLieuOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMNguyenVatLieuOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MNguyenVatLieuFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.MNguyenVatLieuFieldName, model.MNguyenVatLieuFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.MNguyenVatLieuFieldCode, len(args))
	//}

	b := &model.MNguyenVatLieu{}
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

type MNguyenVatLieuData struct {
	*model.MNguyenVatLieu
}

func (r *sMNguyenVatLieuRepo) Search(ctx context.Context, s *SearchMNguyenVatLieuOpts) ([]*MNguyenVatLieuData, error) {
	MNguyenVatLieu := make([]*MNguyenVatLieuData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MNguyenVatLieu)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MNguyenVatLieu, nil
}

func (r *sMNguyenVatLieuRepo) Count(ctx context.Context, s *SearchMNguyenVatLieuOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMNguyenVatLieuRepo.Count: %w", err)
	}

	return countResult, nil
}
