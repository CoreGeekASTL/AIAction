// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."

// Package https contains codes of http client
package https

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"

	"GIDS/common/logger"
)

// ParamType define request input param type
type ParamType string

const (
	// KeyValuePair is key-value pair type, which is stored in map format and converted to JSON format.
	KeyValuePair = "KeyValuePair"
	// Struct is user-defined structure parameter, which is converted to JSON as the parameter body.
	Struct = "Struct"
	// IOReader reads from streams and writes them to the HTTP request body.
	IOReader = "IOReader"
	// DefaultParamType default param type is KeyValuePair
	DefaultParamType = KeyValuePair
)

const (
	// InitialBackOff is the initial backoff time
	InitialBackOff = 2 * time.Second
	// MaxBackoff is the max backoff time
	MaxBackoff = 60 * time.Second
)

var (
	RetryStatusCode = []int{http.StatusTooManyRequests, http.StatusBadGateway, http.StatusGatewayTimeout, http.StatusServiceUnavailable}
)

// BuilderCompleter is an interface that defines the actions allowed after a request is constructed.
type BuilderCompleter interface {
	Do() Response
}

// Builder define a request builder
type Builder interface {
	Method(method string) Builder
	Context(ctx context.Context) Builder
	URL(url string) Builder
	Header(key, value string) Builder
	Headers(map[string]string) Builder
	Param(key string, value interface{}) Builder
	Params(map[string]interface{}) Builder
	ParamType(t ParamType) Builder
	ParamFromInterface(interface{}) Builder
	ParamFromReader(io.Reader) Builder
	WithRetry(times int) Builder
	Complete() BuilderCompleter
}

type request struct {
	doFunc        func(req *http.Request) (*http.Response, error)
	completeFunc  func() BuilderCompleter
	method        string
	url           string
	headers       map[string]string
	paramType     ParamType
	params        map[string]interface{}
	paramReader   io.Reader
	paramStruct   interface{}
	ctx           context.Context
	maxRetryTimes int
}

// HTTPDoer is an interface for the one method of http.Client that is used by Pusher
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewRequest create a http request builder
func NewRequest(client HTTPDoer) Builder {
	return newRequest(client)
}

func newRequest(client HTTPDoer) *request {
	doFunc := func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("client is nil")
	}
	if client != nil {
		doFunc = client.Do
	}
	req := &request{
		doFunc:        doFunc,
		method:        "GET",
		paramType:     DefaultParamType,
		maxRetryTimes: 0,
	}
	req.completeFunc = func() BuilderCompleter {
		return req
	}
	return req
}

// Context implemented the Builder interface
func (r *request) Context(ctx context.Context) Builder {
	r.ctx = ctx
	return r
}

// Method implemented the Builder interface
func (r *request) Method(method string) Builder {
	r.method = method
	return r
}

// URL implemented the Builder interface
func (r *request) URL(url string) Builder {
	r.url = url
	return r
}

// Header implemented the Builder interface
func (r *request) Header(key, value string) Builder {
	if r.headers == nil {
		r.headers = make(map[string]string, 1)
	}
	r.headers[key] = value
	return r
}

// Headers implemented the Builder interface
func (r *request) Headers(m map[string]string) Builder {
	if len(m) == 0 {
		return r
	}
	if r.headers == nil {
		r.headers = make(map[string]string, len(m))
	}
	for k, v := range m {
		r.headers[k] = v
	}
	return r
}

// Param implemented the Builder interface
func (r *request) Param(key string, value interface{}) Builder {
	if r.params == nil {
		r.params = make(map[string]interface{}, 1)
	}
	r.params[key] = value
	return r
}

// Params implemented the Builder interface
func (r *request) Params(m map[string]interface{}) Builder {
	if len(m) == 0 {
		return r
	}
	if r.params == nil {
		r.params = make(map[string]interface{}, len(m))
	}
	for k, v := range m {
		r.params[k] = v
	}
	return r
}

// ParamType implemented the Builder interface
func (r *request) ParamType(t ParamType) Builder {
	r.paramType = t
	return r
}

// ParamFromInterface implemented the Builder interface
func (r *request) ParamFromInterface(i interface{}) Builder {
	if i == nil {
		logger.Warnf("[HTTP] http request builder, function ParamFromInterface's param is nil")
		return r
	}
	r.paramStruct = i
	r.paramType = Struct
	return r
}

// ParamFromReader implemented the Builder interface
func (r *request) ParamFromReader(reader io.Reader) Builder {
	if reader == nil {
		logger.Warnf("[HTTP] http request builder, function ParamFromReader's param is nil")
		return r
	}
	r.paramReader = reader
	r.paramType = IOReader
	return r
}

// WithRetry 重试次数
func (r *request) WithRetry(times int) Builder {
	r.maxRetryTimes = times
	return r
}

// Complete implemented the Builder interface
func (r *request) Complete() BuilderCompleter {
	return r.completeFunc()
}

// Do implement the Builder interface
func (r *request) Do() Response {
	httpRequest, err := r.buildHTTPRequest()
	if err != nil {
		return newResponse(&http.Response{}, err)
	}
	var httpResponse *http.Response
	for i := 0; i <= r.maxRetryTimes; i++ {
		httpResponse, err = r.doFunc(httpRequest)
		if err != nil {
			logger.Debugf("http request doFunc returned error: %v", err)
		}
		if !r.shouldRetry(httpResponse, err) {
			break
		}
		backoff := calBackoff(i)
		resetBody(httpRequest)
		if httpResponse != nil && httpResponse.Body != nil {
			if closeErr := httpResponse.Body.Close(); closeErr != nil {
				logger.Errorf("close response body failed: %v", closeErr)
			}
		}
		time.Sleep(backoff)
	}
	if httpResponse == nil {
		return newResponse(&http.Response{}, err)
	}
	return newResponse(httpResponse, err)
}

func resetBody(httpRequest *http.Request) {
	if httpRequest.GetBody != nil {
		body, err := httpRequest.GetBody()
		if err == nil {
			httpRequest.Body = body
		}
	}
}

func calBackoff(attempt int) time.Duration {
	backoff := InitialBackOff * time.Duration(1<<uint(attempt))
	if backoff > MaxBackoff {
		backoff = MaxBackoff
	}
	return backoff
}

func (r *request) shouldRetry(resp *http.Response, err error) bool {
	if err != nil && isTransientNetErr(err) {
		return true
	}
	if resp != nil {
		return isRetryableStatusCode(resp.StatusCode)
	}
	return false
}

func isTransientNetErr(err error) bool {
	var netErr net.Error
	// timeout
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return true
		}
	}
	// connection refused
	if errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	if errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}

func isRetryableStatusCode(code int) bool {
	for _, c := range RetryStatusCode {
		if code == c {
			return true
		}
	}
	return false
}

func (r *request) buildHTTPRequest() (*http.Request, error) {
	var (
		requestURL string
		body       io.Reader
		err        error
	)
	requestURL, err = r.buildURL()
	if err != nil {
		return nil, err
	}

	body, err = r.buildBody()
	if err != nil {
		return nil, err
	}

	if r.ctx == nil {
		r.ctx = context.TODO()
	}

	httpRequest, err := http.NewRequestWithContext(r.ctx, r.method, requestURL, body)
	if err != nil {
		return nil, err
	}

	r.buildHeader(httpRequest)
	httpRequest.Close = true
	return httpRequest, nil
}

func (r *request) buildURL() (string, error) {
	u, err := url.Parse(r.url)
	if err != nil {
		return "", err
	}
	if r.method == "GET" {
		query := u.Query()
		for k, v := range r.params {
			vString, ok := v.(string)
			if !ok {
				logger.Errorf("[HTTP] Get Request' param is not string, key: %s, value: %v", k, v)
			}
			query.Add(k, vString)
		}
		u.RawQuery = query.Encode()
	}
	return u.String(), nil
}

func (r *request) buildBody() (io.Reader, error) {
	if r.method == "GET" {
		return nil, nil
	}
	var body io.Reader
	switch r.paramType {
	case KeyValuePair:
		bodyByte, err := json.Marshal(r.params)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bodyByte)
	case Struct:
		bodyByte, err := json.Marshal(r.paramStruct)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bodyByte)
	case IOReader:
		body = r.paramReader
	default:
		return nil, fmt.Errorf("unexpected param type: %s", r.paramType)
	}
	return body, nil
}

func (r *request) buildHeader(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "cos/v0.0.0")
	if len(r.headers) != 0 {
		for k, v := range r.headers {
			request.Header.Set(k, v)
		}
	}
}

func newResponse(resp *http.Response, err error) Response {
	return &response{response: resp, err: err}
}

// Response an interface that indicates the response result of an HTTP request.
type Response interface {
	StatusCode() int
	Error() error
	ResponseBody() io.ReadCloser
	ResponseToStruct(interface{}) error
	ResponseToWriter(writer io.Writer) (int64, error)
	Response() *http.Response
	IsSuccessCode() bool
}

type response struct {
	response *http.Response
	err      error
}

// StatusCode implement the response interface.
func (r *response) StatusCode() int {
	return r.response.StatusCode
}

// Error implement the response interface.
func (r *response) Error() error {
	return r.err
}

// ResponseBody implement the response interface.
func (r *response) ResponseBody() io.ReadCloser {
	return r.response.Body
}

// ResponseToStruct implement the response interface.
func (r *response) ResponseToStruct(i interface{}) error {
	defer CloseResponseBody(r)
	respByte, err := io.ReadAll(r.response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(respByte, i)
}

// ResponseToWriter implement the response interface.
func (r *response) ResponseToWriter(writer io.Writer) (int64, error) {
	defer CloseResponseBody(r)
	return io.Copy(writer, r.ResponseBody())
}

// Response implement the response interface.
func (r *response) Response() *http.Response {
	return r.response
}

func (r *response) IsSuccessCode() bool {
	return r.StatusCode() >= http.StatusOK && r.StatusCode() < http.StatusMultipleChoices
}

// CloseResponseBody close body
func CloseResponseBody(resp Response) {
	if resp == nil || resp.Response() == nil || resp.Response().Body == nil {
		return
	}
	if err := resp.Response().Body.Close(); err != nil {
		logger.Errorf("close body failed! err: %v", err)
	}
	return
}
