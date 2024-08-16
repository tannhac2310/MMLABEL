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
	flag.StringVar(&pkgName, "pkg", "model", "pkgName")
	flag.StringVar(&outDir, "o", "internal/aurora/model/", "output dir")
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
		if tableName.String != "production_plans" {
			continue
		}
		genModel(pool, tableName)
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

type model struct {
	PkgName     string
	DbName      string
	Name        string
	DbFieldName []string
	ConstName   []string
	FieldName   []string
	FieldType   []string
}

func genModel(db *pgx.ConnPool, tableName pgtype.Text) {
	fmt.Print(time.Now().Format("2006-01-02T15:04:05"), " generate model ", tableName.String, "...")

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

	m := model{
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

const tmpl = `package {{.PkgName}}

const ({{ range $i, $name := .ConstName }}
	{{$name}}          = "{{index $.DbFieldName $i}}"{{ end }}
)

type {{.Name}} struct {
	{{ range $i, $name := .FieldName }}
	{{$name}}          {{index $.FieldType $i}} ` + "`db:\"{{index $.DbFieldName $i}}\"`" + ` {{ end }}
}

func (rcv *{{.Name}}) FieldMap() (fields []string, values []interface{}) {
	fields = []string{ {{ range $name := .ConstName }}
		{{$name}}, {{end}}
	}

	values = []interface{}{ {{ range $name := .FieldName }}
		&rcv.{{$name}}, {{end}}
	}

	return
}

func (*{{.Name}}) TableName() string {
	return "{{.DbName}}"
}`
