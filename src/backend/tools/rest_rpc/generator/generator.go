package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

var handlerTemplate *template.Template

func Generate(filename string, f *ast.File, prefix string, interfaces ...string) string {
	serv, err := parseServices(filename, f.Decls, prefix, interfaces...)
	if err != nil {
		panic(err)
	}
	additionalImports := []string{
		"net/http",
		"net/url",
		"strings",
		"github.com/go-chi/chi",
		"github.com/sirupsen/logrus",
		"e-clinic/src/backend/tools/http/http_client",
		"e-clinic/src/backend/tools/http/http_server",
		"e-clinic/src/backend/tools/metrics",
	}
	hasStructParam := false
	hasJsonResult := false
	hasErrResult := false
	hasBoolOrIntParam := false
	for _, service := range serv {
		for _, endpoint := range service.Endpoints {
			for _, parameter := range endpoint.Parameters {
				if parameter.ParamType == STRUCT_POINTER {
					hasStructParam = true
				}
				if parameter.ParamType == INT || parameter.ParamType == BOOL || parameter.ParamType == FLOAT64 {
					hasBoolOrIntParam = true
				}
			}
			for _, result := range endpoint.Results {
				if result != TRANSACTION || result != COMMON_TRANSACTION {
					endpoint.RestResults = append(endpoint.RestResults, result)
				}
				if result == ERROR {
					hasErrResult = true
				}
				if result == BODY_JSON {
					hasJsonResult = true
				}
			}
		}
	}
	if hasErrResult || hasJsonResult {
		additionalImports = append(additionalImports, "fmt")
	}
	if hasStructParam || hasJsonResult {
		additionalImports = append(additionalImports, "encoding/json")
	}
	if hasStructParam {
		additionalImports = append(additionalImports, "io/ioutil")
	}
	if hasBoolOrIntParam {
		additionalImports = append(additionalImports, "strconv")
	}
	errs, err := extractErrors(filename)
	if err != nil {
		panic(err)
	}
	imports := extractImports(f.Imports, additionalImports)
	b := bytes.NewBufferString("")
	err = handlerTemplate.Execute(b, struct {
		Package  string
		Imports  []string
		Errors   []string
		Services []*service
	}{
		Package:  f.Name.Name,
		Imports:  imports,
		Errors:   errs,
		Services: serv,
	})
	if err != nil {
		panic(err)
	}
	return b.String()
}

func extractImports(imports []*ast.ImportSpec, additionalPaths []string) []string {
	paths := map[string]*ast.ImportSpec{}
	importStrings := []string{}
	for _, imp := range imports {
		if imp.Path == nil {
			panic(fmt.Errorf("import path at pos %d is empty", imp.Path.ValuePos))
		}
		if imp.Name != nil {
			importStrings = append(importStrings, imp.Name.Name+" "+imp.Path.Value)
		} else {
			importStrings = append(importStrings, imp.Path.Value)
		}
		paths[imp.Path.Value] = imp
	}
	for _, path := range additionalPaths {
		imp, found := paths[path]
		if found && imp.Name != nil {
			panic(fmt.Errorf("path %s was already included with alias %s. please remove alias", path, imp.Name.Name))
		}
		if !found {
			importStrings = append(importStrings, "\""+path+"\"")
		}
	}
	return importStrings
}

func errFromObjects(objs map[string]*ast.Object) []string {
	var errorNames []string
	// extract errors from parsed file
	for _, v := range objs {
		if v.Kind == ast.Var && strings.HasPrefix(v.Name, "Err") {
			errorNames = append(errorNames, v.Name)
		}
	}
	return errorNames
}

func extractErrors(filename string) ([]string, error) {
	var errorNames []string
	// extract errors from files in the same directory
	dirPath := filepath.Dir(filename)

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, dirPath+"/"+file.Name(), nil, 0)
		if err != nil {
			return nil, err
		}
		errorNames = append(errorNames, errFromObjects(f.Scope.Objects)...)
	}
	sort.Strings(errorNames)
	return errorNames, nil
}

func getNameFromIdent(ident []*ast.Ident) string {
	result := ""
	for i, name := range ident {
		result += name.Name
		if i+1 < len(ident) {
			result += "."
		}
	}
	switch result {
	case
		"request",
		"instance",
		"log",
		"r",
		ControlEndpointName:
		panic(fmt.Errorf("invalid parameter name: %s, name is used in generated endpoint", result))
	}
	return result
}

func init() {
	var err error
	handlerTemplate, err = template.New("handler").Parse(handlerTemplateString)
	if err != nil {
		panic(err)
	}
}
