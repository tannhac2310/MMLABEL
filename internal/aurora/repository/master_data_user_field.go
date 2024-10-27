package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type MasterDataUserFieldRepo interface {
	Insert(ctx context.Context, e *model.MasterDataUserField) error
	Update(ctx context.Context, e *model.MasterDataUserField) error
	DeleteByMasterDataIDs(ctx context.Context, masterDataIds []string) error
	FindByID(ctx context.Context, id string) (*MasterDataUserFieldData, error)
	Search(ctx context.Context, s *SearchMasterDataUserFieldOpts) ([]*MasterDataUserFieldData, error)
	Count(ctx context.Context, s *SearchMasterDataUserFieldOpts) (*CountResult, error)
}

type sMasterDataUserFieldRepo struct {
}

func NewMasterDataUserFieldRepo() MasterDataUserFieldRepo {
	return &sMasterDataUserFieldRepo{}
}

func (r *sMasterDataUserFieldRepo) Insert(ctx context.Context, e *model.MasterDataUserField) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *sMasterDataUserFieldRepo) DeleteByMasterDataIDs(ctx context.Context, masterDataIds []string) error {
	if len(masterDataIds) == 0 {
		return fmt.Errorf("ids is empty")
	}

	sql := fmt.Sprintf("DELETE FROM master_data_user_field WHERE master_data_id = ANY($1)")
	cmd, err := cockroach.Exec(ctx, sql, masterDataIds)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

func (r *sMasterDataUserFieldRepo) FindByID(ctx context.Context, id string) (*MasterDataUserFieldData, error) {
	e := &MasterDataUserFieldData{}
	err := cockroach.FindOne(ctx, e, "id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("sMasterDataUserFieldRepo.cockroach.FindOne: %w", err)
	}

	return e, nil
}
func (r *sMasterDataUserFieldRepo) Update(ctx context.Context, e *model.MasterDataUserField) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

// SearchMasterDataUserFieldOpts all params is options
type SearchMasterDataUserFieldOpts struct {
	IDs           []string
	MasterDataID  string
	MasterDataIDs []string
	Limit         int64
	Offset        int64
	Sort          *Sort
}

func (s *SearchMasterDataUserFieldOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.MasterDataUserFieldFieldID)
	}

	if s.MasterDataID != "" {
		args = append(args, s.MasterDataID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.MasterDataUserFieldFieldMasterDataID, len(args))
	}

	if len(s.MasterDataIDs) > 0 {
		args = append(args, s.MasterDataIDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($%d)", model.MasterDataUserFieldFieldMasterDataID, len(args))
	}
	b := &model.MasterDataUserField{}
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

type MasterDataUserFieldData struct {
	*model.MasterDataUserField
}

func (r *sMasterDataUserFieldRepo) Search(ctx context.Context, s *SearchMasterDataUserFieldOpts) ([]*MasterDataUserFieldData, error) {
	MasterDataUserField := make([]*MasterDataUserFieldData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&MasterDataUserField)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return MasterDataUserField, nil
}

func (r *sMasterDataUserFieldRepo) Count(ctx context.Context, s *SearchMasterDataUserFieldOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("sMasterDataUserFieldRepo.Count: %w", err)
	}

	return countResult, nil
}
