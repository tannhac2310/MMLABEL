package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MasterDataSelectionRepo interface {
	Insert(ctx context.Context, e *model.MasterDataSelection) error
	Update(ctx context.Context, e *model.MasterDataSelection) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*MasterDataSelectionData, error)
	Search(ctx context.Context, s *SearchMasterDataSelectionOpts) ([]*MasterDataSelectionData, error)
	Count(ctx context.Context, s *SearchMasterDataSelectionOpts) (*CountResult, error)
}

type sMasterDataSelectionRepo struct {
}

func NewMasterDataSelectionRepo() MasterDataSelectionRepo {
	return &sMasterDataSelectionRepo{}
}

func (r *sMasterDataSelectionRepo) Insert(ctx context.Context, e *model.MasterDataSelection) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMasterDataSelectionRepo) FindByID(ctx context.Context, id string) (*MasterDataSelectionData, error) {
	e := &MasterDataSelectionData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMasterDataSelectionRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMasterDataSelectionRepo) Update(ctx context.Context, e *model.MasterDataSelection) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *sMasterDataSelectionRepo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE master_data_selection SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("master_data_selection cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*sMasterDataSelectionRepo not found any records to delete")
	}

	return nil
}

// SearchMasterDataSelectionOpts all params is options
type SearchMasterDataSelectionOpts struct {
	IDs []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *SearchMasterDataSelectionOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MasterDataSelectionFieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.MasterDataSelectionFieldName, model.MasterDataSelectionFieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.MasterDataSelectionFieldCode, len(args))
	//}

	b := &model.MasterDataSelection{}
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

type MasterDataSelectionData struct {
	*model.MasterDataSelection
}

func (r *sMasterDataSelectionRepo) Search(ctx context.Context, s *SearchMasterDataSelectionOpts) ([]*MasterDataSelectionData, error) {
	MasterDataSelection := make([]*MasterDataSelectionData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MasterDataSelection)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MasterDataSelection, nil
}

func (r *sMasterDataSelectionRepo) Count(ctx context.Context, s *SearchMasterDataSelectionOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMasterDataSelectionRepo.Count: %w", err)
	}

	return countResult, nil
}
