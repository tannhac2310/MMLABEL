package cockroach

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

func Select(ctx context.Context, sql string, args ...interface{}) *Rows {
	rows, err := Query(ctx, sql, args...)
	if err != nil {
		err = fmt.Errorf("err db.Query: %w", err)
	}

	return &Rows{
		err:     err,
		pgxRows: rows,
	}
}

func Create(ctx context.Context, e Entity) error {
	fields, values := e.FieldMap()
	fieldNames := strings.Join(fields, ",")
	placeHolders := generatePlaceholders(len(fields))

	stmt := "INSERT INTO " + e.TableName() + " (" + fieldNames + ") VALUES (" + placeHolders + ");"
	cmdTag, err := Exec(ctx, stmt, values...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return fmt.Errorf("cannot insert new " + e.TableName())
	}

	return nil
}

func Update(ctx context.Context, e Entity) error {
	fields, values := e.FieldMap()

	// get primary field at first
	primaryField := fields[0]
	// remove primary key from fields
	fields = fields[1:]

	placeHolders := generateUpdatePlaceholders(fields)

	stmt := fmt.Sprintf("UPDATE %s SET %s WHERE "+primaryField+" = $%d;", e.TableName(), placeHolders, len(fields)+1)

	// move primary field to last, matching where id = $
	args := append(values[1:], values[0])

	cmdTag, err := Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return fmt.Errorf("cannot update " + e.TableName())
	}

	return nil
}

func Upsert(ctx context.Context, e Entity, pkField string, upsertFields ...string) error {
	fields, _ := e.FieldMap()
	fieldNames := strings.Join(fields, ",")
	placeHolders := generatePlaceholders(len(fields))

	var upsertStatement strings.Builder
	space := ", "

	if len(upsertFields) > 0 {
		for i, field := range fields {
			for _, upsertField := range upsertFields {
				if upsertField == field {
					upsertStatement.WriteString(upsertField + " = $" + strconv.Itoa(i+1) + space)
					break
				}
			}
		}
	} else {
		for i, field := range fields {
			if !(field == "id" || field == "created_at") {
				upsertStatement.WriteString(field + " = $" + strconv.Itoa(i+1) + space)
			}
		}
	}

	stmt := "INSERT INTO " + e.TableName() + " (" + fieldNames + ") VALUES (" + placeHolders + ") ON CONFLICT (" + pkField + ") DO UPDATE SET " + upsertStatement.String()[:upsertStatement.Len()-2]
	args := getScanFields(e, fields)
	cmdTag, err := Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return fmt.Errorf("cannot insert new " + e.TableName())
	}

	return nil
}

type Updater interface {
	Set(fieldName string, value interface{})
}

type update struct {
	tableName  string
	primaryKey string
	primaryVal interface{}
	data       map[string]interface{}
}

func NewUpdater(tableName string, primaryKey string, primaryVal interface{}) Updater {
	return &update{tableName: tableName, primaryKey: primaryKey, primaryVal: primaryVal, data: make(map[string]interface{})}
}

func (rcv *update) Set(fieldName string, value interface{}) {
	rcv.data[fieldName] = value
}

func UpdateFields(ctx context.Context, updater Updater) error {
	// get primary field at first
	var (
		fields []string
		args   []interface{}
		u      = updater.(*update)
	)
	for k, v := range u.data {
		fields = append(fields, k)
		args = append(args, v)
	}

	placeHolders := generateUpdatePlaceholders(fields)

	stmt := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d;", u.tableName, placeHolders, u.primaryKey, len(fields)+1)

	// move primary field to last, matching where id = $
	args = append(args, u.primaryVal)
	cmdTag, err := Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func FindOne(ctx context.Context, e Entity, conds string, args ...interface{}) error {
	fields, values := e.FieldMap()
	sql := fmt.Sprintf(`SELECT %s
		FROM %s
		WHERE %s AND deleted_at IS NULL
	`, strings.Join(fields, ","), e.TableName(), conds)

	return QueryRow(ctx, sql, args...).Scan(values...)
}

func FindMany(ctx context.Context, e Entity, result interface{}, conds string, args ...interface{}) error {
	fields, _ := e.FieldMap()
	sql := fmt.Sprintf(`SELECT %s
		FROM %s
		WHERE %s AND deleted_at IS NULL
	`, strings.Join(fields, ","), e.TableName(), conds)

	return Select(ctx, sql, args...).ScanAll(result)
}

func generatePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}

	var builder strings.Builder
	sep := ", "
	for i := 1; i <= n; i++ {
		if i == n {
			sep = ""
		}
		builder.WriteString("$" + strconv.Itoa(i) + sep)
	}

	return builder.String()
}

// getScanFields return scan fields from model.Entity
func getScanFields(e Entity, names []string) []interface{} {
	fields, values := e.FieldMap()

	n := len(values)
	if len(names) < n {
		n = len(names)
	}

	result := make([]interface{}, 0, n)
	for _, name := range names {
		for i, fieldName := range fields {
			if name == fieldName {
				result = append(result, values[i])
				break
			}
		}
	}

	return result
}

func generateUpdatePlaceholders(fields []string) string {
	var builder strings.Builder
	sep := ", "

	totalField := len(fields)
	for i, field := range fields {
		if i == totalField-1 {
			sep = ""
		}

		builder.WriteString(field + " = $" + strconv.Itoa(i+1) + sep)
	}

	return builder.String()
}

type QueryOption func(query *string)

func WithShareLock() QueryOption {
	return func(query *string) {
		*query += " FOR SHARE"
	}
}

func WithUpdateLock() QueryOption {
	return func(query *string) {
		*query += " FOR UPDATE"
	}
}
