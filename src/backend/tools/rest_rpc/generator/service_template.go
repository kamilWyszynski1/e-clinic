package generator

const handlerTemplateString = `package {{.Package}}
// Code generated by rest_rpc. DO NOT EDIT.

import (
{{range .Imports}}	{{.}}
{{end}}
)


{{ $errs := .Errors }}

{{range .Services}}{{$serviceName := .Name}}
const {{.Name}}RoutePattern = "{{.Prefix}}/{{.Name}}"


func Register{{$serviceName}}(instance {{$serviceName}}, r *chi.Mux, log logrus.FieldLogger) {
	r.Route("{{.Prefix}}/{{.Name}}", func(r chi.Router) {{"{"}}{{range .Endpoints}}
{{- if .IsPost }}
		r.Post("/{{.Name}}", func(writer http.ResponseWriter, request *http.Request) {
{{- else }}
		r.Get("/{{.Name}}", func(writer http.ResponseWriter, request *http.Request) {
{{- end }}
{{- if .IsControl }}
	writer.WriteHeader(http.StatusOK)
{{- else }}
{{range .Parameters}}			
{{- if ne .ParamType 4}}// Reading argument {{.Name}}
{{if and (ne .ParamType 3) (ne .ParamType 7)}}			{{.Name}}{{if ne .ParamType 1}}Str{{end}} := request.URL.Query().Get("{{.Name}}")
{{end}}
{{- end}}{{if eq .ParamType 0}}			{{.Name}}, err := strconv.Atoi({{.Name}}Str)
			if err != nil {
				log.WithError(err).Errorf("Could not parse {{.Name}} parameter value. Expected integer got %s", {{.Name}}Str)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

{{end}}{{if eq .ParamType 2}}			{{.Name}}, err := strconv.ParseBool({{.Name}}Str)
			if err != nil {
				log.WithError(err).Errorf("Could not parse {{.Name}} parameter value. Expected bool got %s", {{.Name}}Str)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}


{{end}}{{if eq .ParamType 6}}			{{.Name}}, err := strconv.ParseFloat({{.Name}}Str, 64)
			if err != nil {
				log.WithError(err).Errorf("Could not parse {{.Name}} parameter value. Expected float got %s", {{.Name}}Str)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

{{end}}{{if eq .ParamType 7}}			{{.Name}} := request.Context()
{{end}}{{if eq .ParamType 3}}			rawBody, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.WithError(err).Error("can't read body")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			var {{.Name}} *{{.TypeName}}
			if len(rawBody) != 0 {
				{{.Name}} = &{{.TypeName}}{}
				if err := json.Unmarshal(rawBody, {{.Name}}); err != nil {
					log.WithError(err).Error("can't unmarshal body")
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
			}

{{end}}{{if eq .ParamType 5}}			{{.Name}}, err := uuid.FromString({{.Name}}Str)
			if err != nil {
				log.WithError(err).Error("can't parse uuid")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
{{end}}{{end}}			// Calling function {{.Name}}
			{
				{{- $hasTransaction := .HasTransaction -}}
				{{- if .Results}}
				{{range $i, $r := .Results}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}responseCode{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}response{{end}}{{if eq $r 3}}transaction{{end}}{{if eq $r 4}}transaction{{end}}{{if eq $r 5}}cfg{{end}}{{end}} := instance.{{.Name}}({{range $i, $p := .Parameters}}{{if gt $i 0}}, {{end}}{{$p.Name}}{{end}}){{$hasResponseCode := false}}{{$results := .Results}}{{range .Results}}{{if eq . 0}}
				{{if not $hasTransaction}}writer.WriteHeader(responseCode){{$hasResponseCode = true}}{{end}}{{end}}{{end}}
				{{- if .HasTransaction }}
				responseCode, data := common.HandleResponse(response, log, transaction, err)
				writer.WriteHeader(responseCode)
				_, err = writer.Write(data)
				if err != nil {
					log.WithError(err).Error("Could not write error response")
				}
				{{- else }}
				{{- range .Results}}{{if eq . 1}}
				if err != nil {{"{"}}{{if not $hasResponseCode}}
					writer.WriteHeader(http.StatusInternalServerError){{end}}
					if err := json.NewEncoder(writer).Encode(&struct{Error string}{Error: err.Error()}); err != nil {
						fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
						log.WithError(err).Error("Could not marshal response error")
					}
					return
				}{{end}}{{end}}{{range .Results}}{{if or (eq . 2) (eq . 5)}}
				{{- $name := ""}}
				{{- if eq . 2 }} {{$name = "response"}} {{- end}}
				{{- if eq . 5 }} {{$name = "cfg"}} {{- end}}
				if {{$name}} != nil {
					if err := json.NewEncoder(writer).Encode({{$name}}); err != nil {
						log.WithError(err).Error("Could not write {{$name}}")
						writer.WriteHeader(http.StatusInternalServerError)
						if err := json.NewEncoder(writer).Encode(&struct{Error string}{Error: err.Error()}); err != nil {
							fmt.Fprintf(writer, "{\"Error\": \"Could not marshal response error\"}")
							log.WithError(err).Error("Could not marshal response error")
						}
						return
					}
				}{{end}}{{end}}{{end}}{{end}}
			}{{end}}
		}){{end}}
	})
}

type {{.Name}}RestClient struct {
	address      string
	client       http_client.DefaultHTTPClient
	errs         []error
}

func New{{.Name}}RestClient(address string, client http_client.DefaultHTTPClient) (*{{.Name}}RestClient, error) {
	packageErrors := []error{ {{range $i, $err := $errs}} {{if gt $i 0}}, {{end}}{{$err}}{{end}} }
	
	restClient := &{{.Name}}RestClient{
		address: address,
		client: client,
		errs: packageErrors,
	}
	
	// restClient is returned to avoid panics, some services may work without working rest client
	if code := restClient.isUp(); code != http.StatusOK {
		return restClient, http_server.ServerNotUp
	}
	return restClient, nil
}


func New{{.Name}}RestClientErrWrap(address string, client http_client.DefaultHTTPClient, log logrus.FieldLogger) *{{.Name}}RestClient {
	restClient, err := New{{.Name}}RestClient(address, client)
	if err != nil {
		log.WithError(err).Warn("failed to create {{.Name}}RestClient")
	}
	return restClient
}

func (serviceInstance *{{$serviceName}}RestClient) wrapError(errMsg interface{}) error {
	for _, er := range serviceInstance.errs {
		if strings.EqualFold(errMsg.(string), er.Error()) {
			return er
		}
	}
	return fmt.Errorf("%v", errMsg)
}

{{$prefix := .Prefix}}
{{$hasWriter := false}}
{{range .Endpoints}}
func (serviceInstance *{{$serviceName}}RestClient) {{.Name}}({{range $i, $param := .Parameters}}{{if eq $param.ParamType 4}}{{$hasWriter = true}}{{end}}{{if gt $i 0}}, {{end}}{{$param.String}}{{end}}) {{if gt (len .Results) 1}}({{end}}{{$resultJsonTypeName := .ResultJsonTypeName}}{{range $i, $r := .Results}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}int{{end}}{{if eq $r 3}}*dbr.Tx{{end}}{{if eq $r 4}}common.Tx{{end}}{{if eq $r 1}}error{{end}}{{if eq $r 2}}*{{$resultJsonTypeName}}{{end}}{{if eq $r 5}}airline.Config{{end}}{{end}}{{if gt (len .Results) 1}}){{end}} {
	u, err := url.Parse(serviceInstance.address + "{{$prefix}}/{{$serviceName}}/{{.Name}}")
	if err != nil {
		return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
	}{{$hasParam := 0}}{{$hasJsonParam := 0}}{{range .Parameters}}{{if ne .ParamType 3}}{{$hasParam = 1}}{{else}}{{$hasJsonParam = 1}}{{end}}{{end}}{{if $hasParam}}
	query := u.Query(){{end}}{{range .Parameters}}{{if ne .ParamType 3}}{{if ne .ParamType 4}}{{if ne .ParamType 7}}
	query.Set("{{.Name}}", {{if eq .ParamType 0}}strconv.Itoa({{.Name}}){{end}}{{if eq .ParamType 1}}{{.Name}}{{end}}{{if eq .ParamType 5}}{{.Name}}.String(){{end}}{{if eq .ParamType 2}}strconv.FormatBool({{.Name}}){{end}}{{if eq .ParamType 6}}strconv.FormatFloat({{.Name}}, 'f', -1, 64){{end}}){{end}}{{end}}{{end}}{{end}}
	{{if $hasParam}}
	u.RawQuery = query.Encode(){{end}}
	{{$Results := .Results}}{{range .Parameters}}{{if eq .ParamType 3}}
	data, err := json.Marshal({{.Name}})
	if err != nil {
		return{{range $i, $r := $Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
	}{{end}}{{end}}
	
	{{- if .IsPost }}
		{{- if .HasContext }}
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(data))
			if err != nil {
				return{{range $i, $r := $Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
			resp, err := serviceInstance.client.Do(req)
			if err != nil {
				return{{range $i, $r := $Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
		{{- else }}
			resp := serviceInstance.client.Post(u.String(), {{if $hasJsonParam}}{{range .Parameters}}{{if eq .ParamType 3}}data{{end}}{{end}}{{else}}nil{{end}})
			if resp.Err != nil {
				return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}resp.Err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
		{{- end }}
	{{- else }}
		{{- if .HasContext }}
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
			if err != nil {
				return{{range $i, $r := $Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
			resp, err := serviceInstance.client.Do(req)
			if err != nil {
				return{{range $i, $r := $Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
		{{- else }}
			resp := serviceInstance.client.Get(u.String())
			if resp.Err != nil {
				return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}resp.Err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
		{{- end }}
	{{- end }}
	{{$hasBodyJson := 0}}{{$hasRespErr := 0}}{{range .Results}}{{if eq . 2}}{{$hasBodyJson = 1}}{{end}}{{if eq . 1}}{{$hasRespErr = 1}}{{end}}{{end}}{{if $hasRespErr}}
	{{- if .HasContext }}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return{{range $i, $r := $Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
		}
		if len(body) > 0 {
	{{- else }}
		if len(resp.Body) > 0 {
	{{- end }}
		{{- if $hasWriter}}
		if resp.StatusCode == http.StatusOK {
			if _, err := writer.Write(resp.Body); err != nil {
				return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			} else {
				return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}-1{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}nil{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
		}
		{{- end}}
		respErr := map[string]interface{}{}
		{{- if .HasContext }}
			if err := json.Unmarshal(body, &respErr); err != nil {
				return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}resp.StatusCode{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
		{{- else }}
			if err := json.Unmarshal(resp.Body, &respErr); err != nil {
				return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}resp.StatusCode{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
			}
		{{- end }}
		if errorMessage, found := respErr["Error"]; found && len(respErr) == 1 {
			return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}resp.StatusCode{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}serviceInstance.wrapError(errorMessage){{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
		}
	}{{end}}{{if $hasBodyJson}}
	respJson := &{{.ResultJsonTypeName}}{}
	{{- if .HasContext }}
		if len(body) > 0 {
		if err := json.Unmarshal(body, respJson); err != nil {
	{{- else }}
		if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, respJson); err != nil {
	{{- end }}
			return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}resp.StatusCode{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 1}}err{{end}}{{if eq $r 2}}nil{{end}}{{if eq $r 5}}nil{{end}}{{end}}
		}
	}{{end}}
	return{{range $i, $r := .Results}}{{if eq $i 0}} {{end}}{{if gt $i 0}}, {{end}}{{if eq $r 0}}resp.StatusCode{{end}}{{if eq $r 1}}nil{{end}}{{if eq $r 3}}nil{{end}}{{if eq $r 4}}nil{{end}}{{if eq $r 2}}respJson{{end}}{{if eq $r 5}}nil{{end}}{{end}}
	{{- $hasWriter = false }}
}
{{end}}{{end}}`
