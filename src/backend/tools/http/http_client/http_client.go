package http_client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

type DefaultHTTPClient interface {
	Get(url string) *SimpleResponse
	Post(url string, rawBody interface{}) *SimpleResponse
	SendForm(method string, url string, values url.Values, header http.Header) *SimpleResponse
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	cli           *retryablehttp.Client
	customHeaders CustomHeaders
}

type ClientConf struct {
	//Max retry attempts
	RetryMax int
	//Maximum wait time for retry
	RetryWaitMax time.Duration
	//Minimum wait time for retry
	RetryWaitMin time.Duration
	//Client timeout
	Timeout time.Duration
	//Toggle switch for keep alives
	DisableKeepAlives bool
	//Max idle connections (keep-alive), 0 = no limit
	MaxIdleConns int
	//Timout for keep alive connections (higher = worse load balancing)
	IdleConnTimeout time.Duration
	//Custom headers that will be added to every request made using method other than Do
	CustomHeaders CustomHeaders
	//Status codes ranges for retry (default = retry for everything in <500, 599>
	RetryStatusCodes []*StatusCodeRange
}

type CustomHeaders map[string][]string

type StatusCodeRange struct {
	Min int
	Max int
}

var DefaultConf = ClientConf{
	RetryMax:          3,
	RetryWaitMax:      50 * time.Millisecond,
	RetryWaitMin:      10 * time.Millisecond,
	Timeout:           3 * time.Second,
	DisableKeepAlives: false,
	MaxIdleConns:      0,
	IdleConnTimeout:   15 * time.Second,
	CustomHeaders:     CustomHeaders{},
	RetryStatusCodes:  []*StatusCodeRange{{Min: 500, Max: 599}},
}

var DefaultConfNoRetry = ClientConf{
	RetryMax:          0,
	Timeout:           3 * time.Second,
	DisableKeepAlives: false,
	MaxIdleConns:      0,
	IdleConnTimeout:   15 * time.Second,
	CustomHeaders:     CustomHeaders{},
	RetryStatusCodes:  []*StatusCodeRange{{Min: 500, Max: 599}},
}

func NewDefault(log logrus.FieldLogger) DefaultHTTPClient {
	return New(log, &DefaultConf)
}

// NewTestOnly returns valid DefaultHTTPClient with provided http.Client
// function created in order to allow httpmock lib work with DefaultHTTPClient interface
func NewTestOnly(log logrus.FieldLogger, httpCli *http.Client, t *testing.T) DefaultHTTPClient {
	if t == nil {
		panic("test only function")
	}
	c := NewWithoutRetries(log)
	c.(*client).cli.HTTPClient = httpCli
	return c
}

//Provided config must include all fields (best way => copy default and change desired fields)
func New(log logrus.FieldLogger, config ...*ClientConf) DefaultHTTPClient {
	conf := &DefaultConf
	if len(config) > 0 {
		conf = config[0]
	}

	cli := retryablehttp.NewClient()
	cli.RetryWaitMax = conf.RetryWaitMax
	cli.RetryWaitMin = conf.RetryWaitMin
	cli.RetryMax = conf.RetryMax
	cli.RequestLogHook = nil
	if log != nil {
		cli.Logger = NewRHLog(log)
	} else {
		cli.Logger = nil
	}
	cli.CheckRetry = GetStatusCodeRetryPolicy(conf.RetryStatusCodes, log)
	cli.HTTPClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: conf.DisableKeepAlives,
			MaxIdleConns:      conf.MaxIdleConns,
			IdleConnTimeout:   conf.IdleConnTimeout,
		},
		Timeout: conf.Timeout,
	}

	return &client{cli: cli, customHeaders: conf.CustomHeaders}
}

func NewWithoutRetries(log logrus.FieldLogger, config ...*ClientConf) DefaultHTTPClient {
	conf := &DefaultConf
	if len(config) > 0 {
		conf = config[0]
	}

	cli := retryablehttp.NewClient()
	cli.RetryMax = 0
	cli.RequestLogHook = nil
	if log != nil {
		cli.Logger = NewRHLog(log)
	} else {
		cli.Logger = nil
	}
	cli.CheckRetry = GetStatusCodeRetryPolicy([]*StatusCodeRange{{Min: 501, Max: 599}}, log)
	cli.HTTPClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: conf.DisableKeepAlives,
			MaxIdleConns:      conf.MaxIdleConns,
			IdleConnTimeout:   conf.IdleConnTimeout,
		},
		Timeout: conf.Timeout,
	}

	return &client{cli: cli, customHeaders: conf.CustomHeaders}
}

type SimpleResponse struct {
	Body       []byte
	StatusCode int
	Err        error
}

//Do doesn't use custom headers field from config
func (client *client) Do(req *http.Request) (*http.Response, error) {
	rreq, err := ToRetryableReq(req)
	if err != nil {
		return nil, fmt.Errorf("client can't convert req to retryable req: %w", err)
	}

	return client.cli.Do(rreq)
}

func (client *client) Get(url string) *SimpleResponse {
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return &SimpleResponse{Err: fmt.Errorf("can't create req: %w", err)}
	}

	return client.do(req)
}

//Post can take []byte or reader as rawBody
func (client *client) Post(url string, rawBody interface{}) *SimpleResponse {
	req, err := retryablehttp.NewRequest("POST", url, rawBody)
	if err != nil {
		return &SimpleResponse{Err: fmt.Errorf("can't create req: %w", err)}
	}

	return client.do(req)
}

func (client *client) SendForm(method string, url string, values url.Values, header http.Header) *SimpleResponse {
	req, err := retryablehttp.NewRequest(method, url, strings.NewReader(values.Encode()))
	if err != nil {
		return &SimpleResponse{Err: fmt.Errorf("can't create req: %w", err)}
	}
	if header != nil {
		req.Header = header
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return client.do(req)
}

func (client *client) do(req *retryablehttp.Request) *SimpleResponse {
	client.addCustomHeaders(req)

	resp, err := client.cli.Do(req)
	if err != nil {
		return &SimpleResponse{Err: fmt.Errorf("request do failed: %w", err)}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &SimpleResponse{Err: fmt.Errorf("can't read body %w", err)}
	}

	return &SimpleResponse{Err: nil, Body: body, StatusCode: resp.StatusCode}
}

func (client *client) addCustomHeaders(req *retryablehttp.Request) {
	for key, valFn := range client.customHeaders {
		for _, val := range valFn {
			req.Header.Add(key, val)
		}
	}
}
