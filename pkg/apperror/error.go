package apperror

import (
	"bytes"
	"fmt"
	"html/template"
)

type AppError struct {
	Code         ErrCode
	debugMessage string
	args         []interface{}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Code: " + e.Code.String() + ", Message: " + e.debugMessage)
}

func (e *AppError) Translate(lang string) string {
	return fmt.Sprintf(e.Code.GetMessage(lang), e.args...)
}

func (e *AppError) WithDebugMessage(msg string) *AppError {
	e.debugMessage = msg
	return e
}

func (e *AppError) WithArgs(args ...interface{}) *AppError {
	e.args = args
	return e
}

const tmpl = `
# Error code
| Code      | Name      |  Message |
| :---        | :---        | :---  |
{{ range $code, $msg := . }}| {{ printf "%d" $code }} | {{ $msg }} | - en_US: {{ toMessage $code "en_US" }} </br> - vi_VN: {{ toMessage $code "vi_VN" }} |
{{ end }}
`

func ExposeDocs() string {
	t := template.Must(template.New("tmpl").Funcs(template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"toMessage": func(e ErrCode, lang string) string {
			return e.GetMessage(lang)
		},
	}).Parse(tmpl))

	var outFile bytes.Buffer
	_ = t.Execute(&outFile, ErrCode_name)

	return outFile.String()
}
