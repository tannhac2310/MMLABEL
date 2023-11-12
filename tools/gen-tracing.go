package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"golang.org/x/tools/imports"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/group"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/role"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/upload"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/user"
)

const tmpl = `//DO NOT EDIT: code generated from 'tools/gen-tracing.go'
package {{.packageName}}

import (
	"context"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/tracingutil"
	{{ range $pkg := .pkgs }}  "{{ $pkg }}"
	{{end}}
)

type {{ .structName  }}Trace struct {
	{{ .fromPkg }}.{{ .IName }}
}

func New{{upperCase .structName  }}Trace( {{ range $a := .Args }}  {{ $a.name }} {{ $a.type }},
{{end}}
) {{ .fromPkg }}.{{ .IName }} {
	return &{{ .structName  }}Trace{
		{{ .IName }}: {{ .fromPkg }}.New{{ .IName  }} (
		{{ range $a := .Args }}  {{ $a.name }},
		{{end}}
		),
	}
}

{{ range $m := .methods }}
func (s *{{ $.structName }}Trace) {{$m.name}}({{ range $in := $m.in }}{{ $in.name }} {{ $in.type }}, {{end}})({{ range $out := $m.out }}{{ $out.type }},{{end}}) {
	ctx, span := tracingutil.Start(ctx, "{{ $.structName }}.{{ $m.name }}")
	defer span.End()

	{{ $length := len $m.out }}
	{{ if ne $length 0 }} return 	{{ end }} s.{{ $.IName }}.{{ $m.name }}({{ range $in := $m.in }}{{ $in.name }}, {{end}})
}
{{end}}
`

func main() {
	genTracing(repository.NewGroupRepo, "repository", "./pkg/tracingutil/repository/group_trace.go", "groupRepo")
	genTracing(repository.NewRoleRepo, "repository", "./pkg/tracingutil/repository/role_trace.go", "roleRepo")
	genTracing(repository.NewUserRepo, "repository", "./pkg/tracingutil/repository/user_trace.go", "userRepo")
	genTracing(repository.NewUserFCMTokenRepo, "repository", "./pkg/tracingutil/repository/user_fcm_token_trace.go", "userFCMTokenRepo")
	genTracing(repository.NewUserFirebaseRepo, "repository", "./pkg/tracingutil/repository/user_firebase_trace.go", "userFirebaseRepo")
	genTracing(repository.NewUserGroupRepo, "repository", "./pkg/tracingutil/repository/user_group_trace.go", "userGroupRepo")
	genTracing(repository.NewUserNotificationRepo, "repository", "./pkg/tracingutil/repository/user_notification_trace.go", "userNotificationRepo")
	genTracing(repository.NewUserRoleRepo, "repository", "./pkg/tracingutil/repository/user_role.go", "userRoleRepo")
	genTracing(repository.NewUserNamePasswordRepo, "repository", "./pkg/tracingutil/repository/username_password_trace.go", "userNamePasswordRepo")

	genTracing(user.NewService, "service", "./pkg/tracingutil/service/user_trace.go", "userService")
	genTracing(group.NewService, "service", "./pkg/tracingutil/service/group_trace.go", "groupService")
	genTracing(role.NewService, "service", "./pkg/tracingutil/service/role_trace.go", "roleService")
	genTracing(upload.NewService, "service", "./pkg/tracingutil/service/upload_trace.go", "uploadService")
	genTracing(user.NewService, "service", "./pkg/tracingutil/service/user_trace.go", "userService")
}

func genTracing(
	s interface{},
	pkgName,
	outDir string,
	structName string,
) {
	t := reflect.TypeOf(s)

	args := []map[string]interface{}{
		{"name": "ctx", "type": "context.Context"},
	}

	methods := []map[string]interface{}{
		{"name": "ctx", "in": map[string]interface{}{}, "out": map[string]interface{}{}},
	}

	var (
		interfaceName = t.Out(0).String()
		fromPkg       = interfaceName[:strings.LastIndex(interfaceName, ".")]
	)
	interfaceName = interfaceName[strings.LastIndex(interfaceName, ".")+1:]

	mapData := map[string]interface{}{
		"packageName": pkgName,
		"structName":  structName,
		"IName":       interfaceName,
		"fromPkg":     fromPkg,
		"Args":        args,
		"methods":     methods,
	}

	pkgs := []string{}
	toArgs := func(t reflect.Type) []map[string]interface{} {
		args = []map[string]interface{}{}
		for i := 0; i < t.NumIn(); i++ {
			name := fmt.Sprintf("arg%d", i+1)
			tType := t.In(i).String()
			if tType == "context.Context" {
				name = "ctx"
			}

			args = append(args, map[string]interface{}{"name": name, "type": tType})

			func() {
				defer func() {
					if r := recover(); r != nil {
						//fmt.Println("Recovered in f", r)
					}
				}()

				pkgs = append(pkgs, t.In(i).Elem().PkgPath())
			}()
		}

		return args
	}

	toOuts := func(t reflect.Type) []map[string]interface{} {
		args = []map[string]interface{}{}
		for i := 0; i < t.NumOut(); i++ {
			args = append(args, map[string]interface{}{"type": t.Out(i).String()})
		}

		return args
	}

	mapData["Args"] = toArgs(t)
	mapData["pkgs"] = pkgs

	methods = make([]map[string]interface{}, 0)
	out := t.Out(0)
	for i := 0; i < out.NumMethod(); i++ {
		method := out.Method(i)

		methods = append(methods, map[string]interface{}{
			"name": method.Name,
			"in":   toArgs(method.Type),
			"out":  toOuts(method.Type),
		})
	}

	mapData["methods"] = methods

	te := template.Must(template.New("tmpl").Funcs(template.FuncMap{
		"upperCase": func(v string) string {
			return strings.Title(v)
		},
	}).Parse(tmpl))

	var outFile bytes.Buffer
	err := te.Execute(&outFile, mapData)
	if err != nil {
		panic(err)
	}

	res, err := imports.Process(outDir, outFile.Bytes(), &imports.Options{Comments: true})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(outDir, res, 0644)
	if err != nil {
		panic(err)
	}

	log.Println("generated succcess ", outDir)
}
