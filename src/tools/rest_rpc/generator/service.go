package generator

import (
	"go/ast"
)

type service struct {
	Endpoints       []*endpoint
	ControlEndpoint *endpoint
	Prefix          string
	Name            string
}

func (s *service) String() string {
	r := "Service " + s.Name + "\n"
	for _, e := range s.Endpoints {
		r += " " + e.String() + "\n"
	}
	return r
}

func parseServices(filename string, decls []ast.Decl, prefix string, interfaces ...string) ([]*service, error) {
	result := []*service{}
	for _, decl := range decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl != nil {
				for _, spec := range genDecl.Specs {
					switch spec.(type) {
					case *ast.TypeSpec:
						typeSpec := spec.(*ast.TypeSpec)
						switch typeSpec.Type.(type) {
						case *ast.InterfaceType:
							if len(interfaces) > 0 {
								skip := true
								for _, elem := range interfaces {
									if elem == typeSpec.Name.Name {
										skip = false
									}
								}
								if skip {
									break
								}
							}
							endpoints, err := interfaceTypeToEndpoints(filename, typeSpec.Name.Name, typeSpec.Type.(*ast.InterfaceType))
							if err != nil {
								return nil, err
							}
							endpoints = append(endpoints, &endpoint{
								Name:        ControlEndpointName,
								Results:     []resultType{RESPONSE_CODE},
								RestResults: []resultType{RESPONSE_CODE},
								IsControl:   true,
							})
							result = append(result, &service{
								Endpoints: endpoints,
								Name:      typeSpec.Name.Name,
								Prefix:    prefix,
							})
						}
					}
				}
			}
		}
	}
	return result, nil
}
