package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"golang.org/x/tools/imports"
)

var (
	pkgName       string
	outDir        string
	connectionUri string
)

func init() {
	flag.StringVar(&pkgName, "pkg", "repository", "pkgName")
	flag.StringVar(&outDir, "o", "internal/aurora/repository/", "output dir")
	flag.StringVar(&connectionUri, "db", "postgres://root@localhost:26257/postgres?sslmode=disable", "db connection")
	flag.Parse()

	if pkgName == "" || outDir == "" || connectionUri == "" {
		log.Fatal(flag.ErrHelp)
	}
}

func main() {
	connConfig, err := pgx.ParseConnectionString(connectionUri)
	if err != nil {
		log.Panic("cannot parse PG_CONNECTION_URI", err)
	}

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: 32,
	}

	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Panic("cannot create new connection pool to postgres", err.Error())
	}

	rows, err := pool.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		log.Panic(err.Error())
	}

	for rows.Next() {
		var tableName pgtype.Text
		if err := rows.Scan(&tableName); err != nil {
			log.Panic(err)
		}

		if tableName.String == "schema_lock" || tableName.String == "schema_migrations" {
			continue
		}
		if tableName.String != "production_order_device_config" {
			continue
		}
		genRepository(pool, tableName)
	}
}

var nameValues = map[string]string{
	"_bool":        "[]bool",
	"_date":        "[]time.Time",
	"_float4":      "[]float32",
	"_float8":      "[]float64",
	"_int2":        "[]int16",
	"_int4":        "[]int32",
	"_int8":        "[]int",
	"_text":        "[]string",
	"_timestamp":   "[]time.Time",
	"_timestamptz": "[]time.Time",
	"_uuid":        "[]string",
	"_varchar":     "[]string",
	"bool":         "bool",
	"date":         "time.Time",
	"float4":       "float32",
	"float8":       "float64",
	"int2":         "int16",
	"int4":         "int",
	"int8":         "int64",
	"interval":     "time.Duration",
	"json":         "map[string]interface{}",
	"jsonb":        "map[string]interface{}",
	"text":         "string",
	"timestamp":    "time.Time",
	"timestamptz":  "time.Time",
}

var nullValues = map[string]string{
	"text":        "sql.NullString",
	"timestamp":   "sql.NullTime",
	"timestamptz": "sql.NullTime",
	//"float4":      "sql.NullFloat64",
	//"float8":      "sql.NullFloat64",
}

func getType(fieldType, isNUll string) string {
	value, ok := nameValues[fieldType]
	if !ok {
		if isNUll == "YES" {
			return "sql.NullString"
		}
		return "string"
	}

	if isNUll == "YES" {
		v, ok := nullValues[fieldType]
		if ok {
			return v
		}
	}

	return value
}

type repository struct {
	PkgName     string
	DbName      string
	Name        string
	DbFieldName []string
	ConstName   []string
	FieldName   []string
	FieldType   []string
}

func genRepository(db *pgx.ConnPool, tableName pgtype.Text) {
	fmt.Print(time.Now().Format("2006-01-02T15:04:05"), " generate repository ", tableName.String, "...")

	rows, err := db.Query("SELECT COLUMN_NAME, UDT_NAME, is_nullable FROM information_schema.COLUMNS WHERE TABLE_NAME = $1;", &tableName)
	if err != nil {
		log.Panic(err.Error())
	}

	var singularName = tableName.String
	if strings.HasSuffix(tableName.String, "ies") {
		singularName = tableName.String[0:len(tableName.String)-3] + "y"
	} else if strings.HasSuffix(tableName.String, "sses") {
		singularName = tableName.String[0 : len(tableName.String)-2]
	} else if strings.HasSuffix(tableName.String, "s") {
		singularName = tableName.String[0 : len(tableName.String)-1]
	}

	m := repository{
		PkgName: pkgName,
		DbName:  tableName.String,
		Name:    strcase.ToCamel(singularName),
	}
	for rows.Next() {
		var fieldName, fieldType pgtype.Text
		var isNUll pgtype.Text
		if err := rows.Scan(&fieldName, &fieldType, &isNUll); err != nil {
			log.Panic(err)
		}

		m.DbFieldName = append(m.DbFieldName, fieldName.String)
		f := strcase.ToCamel(fieldName.String)
		f = strings.ReplaceAll(f, "Id", "ID")
		m.FieldName = append(m.FieldName, f)
		m.FieldType = append(m.FieldType, getType(fieldType.String, isNUll.String))
		m.ConstName = append(m.ConstName, strcase.ToCamel((singularName)+"Field"+strings.Title(f)))
	}
	t := template.Must(template.New("tmpl").Parse(tmpl))

	var outFile bytes.Buffer
	err = t.Execute(&outFile, m)
	if err != nil {
		panic(err)
	}
	res, err := imports.Process("", []byte(outFile.String()), &imports.Options{Comments: true})
	if err != nil {

		panic(err)
	}

	fileName := outDir + singularName + ".go"
	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("file %s existed, skip\n", fileName)
		return
	}

	err = ioutil.WriteFile(fileName, res, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println(" done. ")
}

const tmpl = `package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

type {{.Name}}Repo interface {
	Insert(ctx context.Context, e *model.{{.Name}}) error
	Update(ctx context.Context, e *model.{{.Name}}) error
	SoftDelete(ctx context.Context, id string) error
	Search(ctx context.Context, s *Search{{.Name}}Opts) ([]*{{.Name}}Data, error)
	Count(ctx context.Context, s *Search{{.Name}}Opts) (*CountResult, error)
}

type s{{.Name}}Repo struct {
}

func New{{.Name}}Repo() {{.Name}}Repo {
	return &s{{.Name}}Repo{}
}

func (r *s{{.Name}}Repo) Insert(ctx context.Context, e *model.{{.Name}}) error {
	err := cockroach.Create(ctx, e)
	if err != nil {
		return fmt.Errorf("r.baseRepo.Create: %w", err)
	}

	return nil
}

func (r *s{{.Name}}Repo) Update(ctx context.Context, e *model.{{.Name}}) error {
	e.UpdatedAt = time.Now()
	return cockroach.Update(ctx, e)
}

func (r *s{{.Name}}Repo) SoftDelete(ctx context.Context, id string) error {
	sql := "UPDATE {{.DbName}} SET deleted_at = NOW() WHERE id = $1;"

	cmd, err := cockroach.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("{{.DbName}} cockroach.Exec: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("*s{{.Name}}Repo not found any records to delete")
	}

	return nil
}

// Search{{.Name}}Opts all params is options
type Search{{.Name}}Opts struct {
	IDs    []string
	// todo add more search options
	Limit  int64
	Offset int64
	Sort   *Sort
}

func (s *Search{{.Name}}Opts) buildQuery(isCount bool) (string, []interface{}) {
	var args []interface{}
	conds := ""
	joins := ""

	if len(s.IDs) > 0 {
		args = append(args, s.IDs)
		conds += fmt.Sprintf(" AND b.%s = ANY($1)", model.{{.Name}}FieldID)
	}
	// todo add more search options example:
	//if s.Name != "" {
	//	args = append(args, "%"+s.Name+"%")
	//	conds += fmt.Sprintf(" AND (b.%[2]s ILIKE $%[1]d OR b.%[3]s ILIKE $%[1]d)",
	//		len(args), model.{{.Name}}FieldName, model.{{.Name}}FieldCode)
	//}
	//if s.Code != "" {
	//	args = append(args, s.Code)
	//	conds += fmt.Sprintf(" AND b.%s ILIKE $%d", model.{{.Name}}FieldCode, len(args))
	//}

	b := &model.{{.Name}}{}
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

type {{.Name}}Data struct {
	*model.{{.Name}}
}

func (r *s{{.Name}}Repo) Search(ctx context.Context, s *Search{{.Name}}Opts) ([]*{{.Name}}Data, error) {
	{{.Name}} := make([]*{{.Name}}Data, 0)
	sql, args := s.buildQuery(false)
	err := cockroach.Select(ctx, sql, args...).ScanAll(&{{.Name}})
	if err != nil {
		return nil, fmt.Errorf("cockroach.Select: %w", err)
	}

	return {{.Name}}, nil
}

func (r *s{{.Name}}Repo) Count(ctx context.Context, s *Search{{.Name}}Opts) (*CountResult, error) {
	countResult := &CountResult{}
	sql, args := s.buildQuery(true)
	err := cockroach.Select(ctx, sql, args...).ScanOne(countResult)
	if err != nil {
		return nil, fmt.Errorf("s{{.Name}}Repo.Count: %w", err)
	}

	return countResult, nil
}
`
