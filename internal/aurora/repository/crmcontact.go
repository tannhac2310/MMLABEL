package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
)

type CrmContactRepo interface {
	Insert(ctx context.Context, e *model.CrmContact) error
	Update(ctx context.Context, e *model.CrmContact) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *SearchCrmContactsOpts) ([]*CrmContactData, error)
	Count(ctx context.Context, s *SearchCrmContactsOpts) (*CountResult, error)
	SearchOne(ctx context.Context, s *SearchCrmContactsOpts) (*CrmContactData, error)
}

type crmContactsRepo struct {
}

func NewCrmContactRepo() CrmContactRepo {
	return &crmContactsRepo{}
}

func (r *crmContactsRepo) Insert(ctx context.Context, e *model.CrmContact) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *crmContactsRepo) Update(ctx context.Context, e *model.CrmContact) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *crmContactsRepo) SoftDelete(ctx context.Context, id string) error {
	sql := `UPDATE message
		SET deleted_at = NOW()
		WHERE id = $1`

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("not found any records to delete")
	}

	return nil
}

// SearchCrmContactsOpts all params is options
type SearchCrmContactsOpts struct {
	IDs         []string
	ID          string
	SourceID    string
	SourceType  enum.CrmContactSource
	ExternalID  string
	PhoneNumber string
	Email       string
	Search      string
	Limit       int64
	Offset      int64
	Sort        *Sort
}

func (s *SearchCrmContactsOpts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.CrmContactFieldID)
	}

	if s.ID != "" {
		args = append(args, s.ID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.CrmContactFieldID, len(args))
	}
	if s.SourceID != "" {
		args = append(args, s.SourceID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.CrmContactFieldSourceID, len(args))
	}
	if s.SourceType > 0 {
		args = append(args, s.SourceType)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.CrmContactFieldSourceType, len(args))
	}
	if s.ExternalID != "" {
		args = append(args, s.ExternalID)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.CrmContactFieldExternalID, len(args))
	}
	if s.Email != "" {
		args = append(args, s.Email)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.CrmContactFieldEmail, len(args))
	}
	if s.PhoneNumber != "" {
		args = append(args, s.PhoneNumber)
		conds += fmt.Sprintf(" AND b.%s = $%d", model.CrmContactFieldPhoneNumber, len(args))
	}
	if s.Search != "" {
		// Email
		args = append(args, "%"+s.Search+"%")
		conds += fmt.Sprintf(" AND (c.%s ILIKE $%d", model.CrmContactFieldEmail, len(args))
		// Phone
		conds += fmt.Sprintf(" OR c.%s ILIKE $%d", model.CrmContactFieldPhoneNumber, len(args))
		// Name
		conds += fmt.Sprintf(" OR c.%s ILIKE $%d )", model.CrmContactFieldName, len(args))
	}

	b := &model.CrmContact{}
	fields, _ := b.FieldMap()
	if isCount {
		return fmt.Sprintf(`SELECT count(*) as cnt
		FROM %s AS b %s
		WHERE TRUE %s AND b.deleted_at IS NULL`, b.TableName(), joins, conds), args
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

type CrmContactData struct {
	*model.CrmContact
}

func (r *crmContactsRepo) Search(ctx context.Context, s *SearchCrmContactsOpts) ([]*CrmContactData, error) {
	message := make([]*CrmContactData, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return message, nil
}
func (r *crmContactsRepo) SearchOne(ctx context.Context, s *SearchCrmContactsOpts) (*CrmContactData, error) {
	message := &CrmContactData{}
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanOne(message)
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select1: %w", err)
	}

	return message, nil
}

func (r *crmContactsRepo) Count(ctx context.Context, s *SearchCrmContactsOpts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("crmcontact.Count: %w", err)
	}

	return countResult, nil
}
