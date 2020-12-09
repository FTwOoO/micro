package rpc

import (
	"bytes"
	"fmt"
	"github.com/fatih/structtag"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
)

type MethodDef struct {
	HandlerName string
	MethodName  string
	ReqName     string
	RespName    string
}

func templateText(tpl string, tplArgv map[string]interface{}) string {
	tpl1, err := template.New("x").Parse(tpl)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	_ = tpl1.Execute(&b, tplArgv)
	return b.String()
}

func camelCase(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func Generate(pkg string, filename string, src []byte, srcOut io.Writer, serverW io.Writer, clientW io.Writer) (err error) {

	if src == nil {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}

	collectStructs := func(x ast.Node) bool {
		ts, ok := x.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range s.Fields.List {
			if field.Tag == nil {
				field.Tag = &ast.BasicLit{
					Kind:  token.STRING,
					Value: "``",
				}
			}
			tag := field.Tag.Value
			tag = strings.Trim(tag, "`")
			tags, err := structtag.Parse(string(tag))
			if err != nil {
				return true
			}

			tagName := camelCase(field.Names[0].Name)

			if tags == nil {
				tags = &structtag.Tags{}
			}

			err = tags.Set(&structtag.Tag{
				Key:  "bson",
				Name: tagName,
			})

			if err != nil {
				return true
			}

			err = tags.Set(&structtag.Tag{
				Key:  "json",
				Name: tagName,
			})

			if err != nil {
				return true
			}

			field.Tag.Value = "`" + tags.String() + "`"

		}
		return false
	}

	ast.Inspect(file, collectStructs)
	err = printer.Fprint(srcOut, token.NewFileSet(), file)
	if err != nil {
		return err
	}

	methos := map[string][]MethodDef{}
	collectInterface := func(x ast.Node) bool {
		ts, ok := x.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}

		s, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		for _, method := range s.Methods.List {
			params := method.Type.(*ast.FuncType).Params.List
			if len(params) != 2 {
				panic("invalid req params")
			}

			reqName := params[1].Type.(*ast.StarExpr).X.(*ast.Ident).Name
			results := method.Type.(*ast.FuncType).Results.List

			var respName = ""
			if len(results) == 1 {
				if v, ok := results[0].Type.(*ast.Ident); ok {
					if v.Name == "error" {
						respName = ""
					}
				}
			} else {
				respName = results[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
				if _, ok := methos[ts.Name.Name]; !ok {
					methos[ts.Name.Name] = []MethodDef{}
				}
			}

			methos[ts.Name.Name] = append(methos[ts.Name.Name], MethodDef{
				HandlerName: ts.Name.Name,
				MethodName:  method.Names[0].Name,
				ReqName:     reqName,
				RespName:    respName,
			})
		}
		return false
	}

	ast.Inspect(file, collectInterface)

	out := fmt.Sprintf(`
      package %s
      import (
		"context"
		"encoding/json"
		"io/ioutil"
		"net/http"
		"time"
	)`, pkg)

	_, err = serverW.Write([]byte(out + "\n"))
	if err != nil {
		return
	}

	for serviceName, methoddefs := range methos {

		out := templateText(rpcServerCommonTpl, map[string]interface{}{
			"serviceName": camelCase(serviceName),
		})
		serverW.Write([]byte(out + "\n"))

		out = fmt.Sprintf(`func Register%sForHTTP(service %s) {`,
			serviceName, serviceName)
		_, err = serverW.Write([]byte(out + "\n"))
		if err != nil {
			return
		}

		for _, method := range methoddefs {
			if method.RespName == "" {
				out = templateText(httpRouteHandlerForNoRespTpl, map[string]interface{}{
					"serviceName":          camelCase(serviceName),
					"httpMethodName":       fmt.Sprintf("/%s/%s", camelCase(serviceName), camelCase(method.MethodName)),
					"rpcRequestStructName": method.ReqName,
					"rpcMethodName":        method.MethodName,
				})
			} else {
				out = templateText(httpRouteHandlerTpl, map[string]interface{}{
					"serviceName":           camelCase(serviceName),
					"httpMethodName":        fmt.Sprintf("/%s/%s", camelCase(serviceName), camelCase(method.MethodName)),
					"rpcRequestStructName":  method.ReqName,
					"rpcResponseStructName": method.RespName,
					"rpcMethodName":         method.MethodName,
				})
			}

			serverW.Write([]byte(out + "\n"))
		}

		_, err = serverW.Write([]byte("}\n"))
		if err != nil {
			return
		}
	}

	out = fmt.Sprintf(`
      package %s
      import (
     "bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	)`, pkg)
	_, err = clientW.Write([]byte(out + "\n"))
	if err != nil {
		return
	}

	for serviceName, methoddefs := range methos {
		out := templateText(RcpClientCommonTpl, map[string]interface{}{
			"structName": serviceName + "Client",
		})
		clientW.Write([]byte(out + "\n"))
		for _, method := range methoddefs {
			if method.RespName == "" {
				out = templateText(RpcClientMethodForNoRespTpl, map[string]interface{}{
					"serviceName":          camelCase(serviceName),
					"structName":           serviceName + "Client",
					"httpMethodName":       fmt.Sprintf("/%s/%s", camelCase(serviceName), camelCase(method.MethodName)),
					"rpcRequestStructName": method.ReqName,
					"rpcMethodName":        method.MethodName,
				})
			} else {
				out = templateText(RpcClientMethodTpl, map[string]interface{}{
					"serviceName":           camelCase(serviceName),
					"structName":            serviceName + "Client",
					"httpMethodName":        fmt.Sprintf("/%s/%s", camelCase(serviceName), camelCase(method.MethodName)),
					"rpcRequestStructName":  method.ReqName,
					"rpcResponseStructName": method.RespName,
					"rpcMethodName":         method.MethodName,
				})
			}
			clientW.Write([]byte(out + "\n"))
		}
	}

	return err
}
