package generator

import (
	"fmt"
	"go/ast"
	"reflect"
)

type parameterType int

const (
	INT parameterType = iota
	STRING
	BOOL
	STRUCT_POINTER
	HTTP_WRITER
	UUID
	FLOAT64
	CONTEXT
)

const ControlEndpointName = "isUp"

var parameterTypeToString = []string{
	"int",
	"string",
	"bool",
	"*",
	"",
	"",
	"float64",
	"",
}

type parameter struct {
	Name      string
	TypeName  string
	ParamType parameterType
}

func (p *parameter) String() string {
	return p.Name + " " + parameterTypeToString[p.ParamType] + p.TypeName
}

type resultType int

const (
	RESPONSE_CODE resultType = iota
	ERROR
	BODY_JSON
	TRANSACTION
	COMMON_TRANSACTION
	AIRLINE_CONFIG
)

const (
	CustomTypeTransaction = "Tx"
	CustomTypeConfig      = "Config"
)

type endpoint struct {
	Name               string
	Parameters         []*parameter
	Results            []resultType
	RestResults        []resultType
	ResultJsonTypeName string
	IsPost             bool
	HasTransaction     bool
	HasContext         bool
	// Control endpoint is used while initializing rest client to check if rest server is up
	IsControl bool
}

func (e *endpoint) String() string {
	r := e.Name + "("
	for i, param := range e.Parameters {
		r += param.String()
		if i+1 < len(e.Parameters) {
			r += ","
		}
	}
	r += ")"
	if len(e.Results) > 1 {
		r += " ("
	}
	for i, result := range e.Results {
		switch result {
		case RESPONSE_CODE:
			r += "response_code"
		case ERROR:
			r += "error"
		case BODY_JSON:
			r += "*" + e.ResultJsonTypeName
		}
		if i+1 < len(e.Results) {
			r += ", "
		}
	}
	if len(e.Results) > 1 {
		r += ")"
	}
	return r
}

func interfaceTypeToEndpoints(filename string, serviceName string, interfaceType *ast.InterfaceType) ([]*endpoint, error) {
	endpointResult := []*endpoint{}
	if interfaceType != nil {
		for _, method := range interfaceType.Methods.List {
			switch method.Type.(type) {
			case *ast.FuncType:
				methodName := getNameFromIdent(method.Names)
				endp := &endpoint{Name: methodName, HasTransaction: false, HasContext: false}
				endpointResult = append(endpointResult, endp)

				funcType := method.Type.(*ast.FuncType)
				if funcType.Params != nil {
					hasStruct := false
					for _, param := range funcType.Params.List {
						parameter, err := paramFieldToParameter(param)
						if err != nil {
							return nil, fmt.Errorf("could not parse parameter of method %s because %w", methodName, err)
						}
						if parameter.ParamType == STRUCT_POINTER {
							if hasStruct {
								return nil, fmt.Errorf("function %s has defined more than one structure arguemnt", methodName)
							}
							hasStruct = true
						} else if parameter.ParamType == CONTEXT {
							endp.HasContext = true
						}
						endp.Parameters = append(endp.Parameters, parameter)
					}
				}
				if funcType.Results != nil {
					hasResponseCode := false
					hasError := false
					hasBodyJson := false
					for _, result := range funcType.Results.List {
						t, name, err := resultToType(result)
						if err != nil {
							return nil, fmt.Errorf("could not parse Results for function %s because: %w", methodName, err)
						}
						switch t {
						case RESPONSE_CODE:
							if hasResponseCode {
								return nil, fmt.Errorf("function %s has defined more than one int response", methodName)
							}
							hasResponseCode = true
						case ERROR:
							if hasError {
								return nil, fmt.Errorf("function %s has defined more than one error response", methodName)
							}
							hasError = true
						case BODY_JSON:
							if hasBodyJson {
								return nil, fmt.Errorf("function %s has defined more than one hasBodyJson response", methodName)
							}
							hasBodyJson = true
							endp.ResultJsonTypeName = name
						case TRANSACTION:
							endp.HasTransaction = true
						case COMMON_TRANSACTION:
							endp.HasTransaction = true
						case AIRLINE_CONFIG:

						default:
							return nil, fmt.Errorf("unknown reponse type for function %s", methodName)
						}
						endp.Results = append(endp.Results, t)
					}
					if !hasError {
						fmt.Printf("WARNING: interface %s method %s from file %s does not return error.\n", serviceName, methodName, filename)
					}
				}
			default:
				return nil, fmt.Errorf("unexpected method type: %s", reflect.TypeOf(method.Type).Name())
			}
		}
	}
	for _, endpoint := range endpointResult {
		endpoint.IsPost = false
		for _, param := range endpoint.Parameters {
			if param.ParamType == 3 {
				endpoint.IsPost = true
				break
			}
		}
	}
	return endpointResult, nil
}

func resultToType(field *ast.Field) (resultType, string, error) {
	switch field.Type.(type) {
	case *ast.Ident:
		ident := field.Type.(*ast.Ident)
		switch ident.Name {
		case "int":
			return RESPONSE_CODE, "", nil
		case "error":
			return ERROR, "", nil
		default:
			return ERROR, "", fmt.Errorf("invalid argument type")
		}
	case *ast.SelectorExpr:
		selector := field.Type.(*ast.SelectorExpr)
		switch selector.X.(type) {
		case *ast.Ident:
			ident2 := selector.X.(*ast.Ident)
			if selector.Sel == nil {
				return ERROR, "", fmt.Errorf("selector has no Sel")
			}
			switch selName := selector.Sel.Name; selName {
			case CustomTypeTransaction:
				return COMMON_TRANSACTION, ident2.Name + "." + selName, nil
			case CustomTypeConfig:
				return AIRLINE_CONFIG, ident2.Name + "." + selName, nil
			default:
				return ERROR, "", fmt.Errorf("unexpected Sel Name %s", selName)
			}
		default:
			return ERROR, "", fmt.Errorf("unexpected selector.X type %s", reflect.TypeOf(selector.X))
		}
	case *ast.StarExpr:
		starExpr := field.Type.(*ast.StarExpr)
		switch starExpr.X.(type) {
		case *ast.Ident:
			ident := starExpr.X.(*ast.Ident)
			return BODY_JSON, ident.Name, nil
		case *ast.SelectorExpr:
			selector := starExpr.X.(*ast.SelectorExpr)
			switch selector.X.(type) {
			case *ast.Ident:
				ident2 := selector.X.(*ast.Ident)
				if selector.Sel == nil {
					return ERROR, "", fmt.Errorf("selector has no Sel")
				}
				switch selName := selector.Sel.Name; selName {
				case CustomTypeTransaction:
					return TRANSACTION, ident2.Name + "." + selector.Sel.Name, nil
				default:
					return BODY_JSON, ident2.Name + "." + selector.Sel.Name, nil
				}
			default:
				return ERROR, "", fmt.Errorf("unexpected selector.X type %s", reflect.TypeOf(selector.X))
			}
		default:
			return ERROR, "", fmt.Errorf("unexpected startExpr.X type %s", reflect.TypeOf(starExpr.X))
		}
	default:
		return ERROR, "", fmt.Errorf("unexpected field.Type type %s", reflect.TypeOf(field.Type))
	}
}

func paramFieldToParameter(field *ast.Field) (*parameter, error) {
	param := &parameter{
		Name: getNameFromIdent(field.Names),
	}

	switch field.Type.(type) {
	case *ast.Ident:
		ident := field.Type.(*ast.Ident)
		switch ident.Name {
		case "int":
			param.ParamType = INT
		case "string":
			param.ParamType = STRING
		case "bool":
			param.ParamType = BOOL
		case "float64":
			param.ParamType = FLOAT64
		default:
			return nil, fmt.Errorf("invalid argument type")
		}
	case *ast.StarExpr:
		starExpr := field.Type.(*ast.StarExpr)
		switch starExpr.X.(type) {
		case *ast.Ident:
			ident := starExpr.X.(*ast.Ident)
			param.ParamType = STRUCT_POINTER
			param.TypeName = ident.Name
		case *ast.SelectorExpr:
			selector := starExpr.X.(*ast.SelectorExpr)
			switch selector.X.(type) {
			case *ast.Ident:
				ident2 := selector.X.(*ast.Ident)
				if selector.Sel == nil {
					return nil, fmt.Errorf("selector has no Sel")
				}
				param.ParamType = STRUCT_POINTER
				param.TypeName = ident2.Name + "." + selector.Sel.Name
			default:
				return nil, fmt.Errorf("unexpected selector.X type %s", reflect.TypeOf(selector.X))
			}
		default:
			return nil, fmt.Errorf("unexpected startExpr.X type %s", reflect.TypeOf(starExpr.X))
		}
	case *ast.SelectorExpr:
		selName := field.Type.(*ast.SelectorExpr).Sel.Name
		switch selName {
		case "Writer":
			param.Name = "writer"
			param.TypeName = "io.Writer"
			param.ParamType = HTTP_WRITER
		case "UUID":
			param.TypeName = "uuid.UUID"
			param.ParamType = UUID
		case "Context":
			param.Name = "ctx"
			param.TypeName = "context.Context"
			param.ParamType = CONTEXT

		}

	default:
		return nil, fmt.Errorf("unexpected field.Type type %s", reflect.TypeOf(field.Type))
	}
	return param, nil
}
