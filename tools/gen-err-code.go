package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"html/template"
	"io/ioutil"
	"log"

	"golang.org/x/tools/imports"
)

var input, output string

func init() {
	flag.StringVar(&input, "i", "./deployments/err_code.csv", "input csv file dir")
	flag.StringVar(&output, "o", "./pkg/apperror/err_code.go", "output config yaml")
	flag.Parse()

	if input == "" || output == "" {
		log.Fatal(flag.ErrHelp)
	}

}

const tmpl = `//DO NOT EDIT: code generated from 'tools/gen-err-code.go'
package apperror

import "strconv"

type ErrCode int32

const (
{{ range $i, $lang := .languages }} Lang{{ $lang }} = "{{ $lang }}"
{{end}}
)


var (
	{{ range $i, $line := .data }}
	{{index $line 0}} = &AppError{Code: {{index $line 1}}}{{end}}
)

var ErrCode_name = map[ErrCode]string{
	{{ range $i, $line := .data }} {{index $line 1}}: "{{index $line 0}}",
{{ end }}
}

var ErrCode_value = map[string]ErrCode{
{{ range $i, $line := .data }}"{{index $line 0}}": {{index $line 1}},
{{ end }}
}

func (x ErrCode) String() string {
	s, ok := ErrCode_name[x]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

{{ $indexLang := 1}}
{{ range $i, $lang := .languages }}
{{ $indexLang = inc $indexLang }}
var ErrCode_{{$lang}} = map[ErrCode]string{
	{{ range $line := $.data }}{{index $line 1}}: "{{index $line $indexLang}}",
	{{ end }}
}
{{end}}


func (x ErrCode) GetMessage(lang string) string {
	switch lang { {{ range $lang := .languages }}
	case Lang{{$lang}}:
		s, ok := ErrCode_{{$lang}}[x]
		if ok {
			return s
		}{{end}}
	default:
		s, ok := ErrCode_{{index .languages 0}}[x]
		if ok {
			return s
		}
	}

	return strconv.Itoa(int(x))
}
`

func main() {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}

	lines, err := csv.NewReader(bytes.NewBuffer(data)).ReadAll()
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("tmpl").Funcs(template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}).Parse(tmpl))

	mapData := map[string]interface{}{
		"title":     lines[0:1],
		"languages": lines[0][2:],
		"data":      lines[1:],
	}

	var outFile bytes.Buffer
	err = t.Execute(&outFile, mapData)
	if err != nil {
		panic(err)
	}

	res, err := imports.Process("", outFile.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(output, res, 0644)
	if err != nil {
		panic(err)
	}

	log.Print("done => ", output)
}
