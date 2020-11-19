package http_client

import (
	"net/http"
	"net/url"
)

var defaultSimpleResponse = &SimpleResponse{
	Body:       nil,
	StatusCode: http.StatusOK,
	Err:        nil,
}

type SendFormRequest struct {
	Id     int
	Method string
	Url    string
	Values url.Values
	Header http.Header
}

type MockClient struct {
	lastId              int
	cannedResponses     []*SimpleResponse
	cannedResponsesSend int
	SendFormRequests    []*SendFormRequest
}

func NewMockEmpty(cannedResponses []*SimpleResponse) *MockClient {
	return &MockClient{
		cannedResponses:     cannedResponses,
		cannedResponsesSend: 0,
		lastId:              0}
}

func (m *MockClient) getSimpleResponse() *SimpleResponse {
	if m.cannedResponsesSend < len(m.cannedResponses) {
		m.cannedResponsesSend++
		return m.cannedResponses[m.cannedResponsesSend-1]
	}
	return defaultSimpleResponse
}

func (m *MockClient) Get(url string) *SimpleResponse {
	return m.getSimpleResponse()
}

func (m *MockClient) Post(url string, rawBody interface{}) *SimpleResponse {
	return m.getSimpleResponse()
}

func (m *MockClient) SendForm(method string, url string, values url.Values, header http.Header) *SimpleResponse {
	m.SendFormRequests = append(m.SendFormRequests, &SendFormRequest{
		Id:     m.lastId,
		Method: method,
		Url:    url,
		Values: values,
		Header: header,
	})
	m.lastId++
	return m.getSimpleResponse()
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:           "OK",
		StatusCode:       http.StatusOK,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}, nil
}

func (m *MockClient) AuthorizeClientFromEnv() error {
	return nil
}
