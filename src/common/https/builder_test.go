// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."
package https

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
* 测试用例描述：Test_request
* 预置条件：
* 操作步骤：
*     1. 创建一个GET Request请求，预置url、header、param
*     2. 创建一个POST请求，使用结构体作为预置请求体
*     3. 创建一个POST请求，使用Reader作为请求体输入源
*     4. 创建一个POST请求，使用KeyValuePair作为请求体输入源
* 预期结果：
*     1. 返回error为nil
*     2. 返回error为nil
*     3. 返回error为nil
*     4. 返回error为nil
* 修改历史：
*     1. 2025-8-13 新建测试用例
*     2. 2025-8-13 补充context的测试
 */
func TestRequest(t *testing.T) {
	ctx := context.TODO()
	tests := []struct {
		name    string
		client  HTTPDoer
		preExec func(r Builder)
		want    error
	}{
		{
			name: "test get request do success",
			client: &fakeClient{
				response: nil,
				err:      nil,
			},
			preExec: func(r Builder) {
				r.Method("GET").
					URL("xx").
					ParamType(DefaultParamType).
					Headers(map[string]string{"a": "b"}).
					Header("xx", "xx").
					Params(map[string]interface{}{"b": ""}).
					Param("xx", "xx")
			},
			want: nil,
		},
		{
			name: "test post request do success",
			client: &fakeClient{
				response: nil,
				err:      nil,
			},
			preExec: func(r Builder) {
				r.Method("POST").
					URL("xx").
					ParamType(Struct).
					ParamFromInterface("")
			},
			want: nil,
		},
		{
			name: "test IOReader param do success",
			client: &fakeClient{
				response: nil,
				err:      nil,
			},
			preExec: func(r Builder) {
				r.Method("POST").
					URL("xx").
					Header("xx", "xx").
					Param("x", "x").
					ParamType(IOReader).
					ParamFromReader(bytes.NewReader([]byte{}))
			},
			want: nil,
		},
		{
			name: "test KeyValuePair param do success",
			client: &fakeClient{
				response: nil,
				err:      nil,
			},
			preExec: func(r Builder) {
				r.Method("POST").
					URL("xx").
					Header("xx", "xx").
					Param("x", "x").
					ParamType(KeyValuePair).
					Context(ctx)
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest(tt.client)
			tt.preExec(req)
			if got := req.Complete().Do(); tt.want != got.Error() {
				t.Errorf("Do() = %v, want %v", got.Error(), tt.want)
			}
		})
	}
}

type fakeClient struct {
	response *http.Response
	err      error
}

func (f *fakeClient) Do(req *http.Request) (*http.Response, error) {
	return f.response, f.err
}

type fakeReader struct {
	byteReader *bytes.Reader
	readFailed bool
}

func (f *fakeReader) Read(p []byte) (int, error) {
	if f.readFailed {
		return 0, errors.New("test")
	}
	return f.byteReader.Read(p)
}

func (f *fakeReader) Close() error {
	return nil
}

/*
* 测试用例描述：Test_response_ResponseToStruct
* 预置条件：
* 操作步骤：
*     1. 构造一个response，无法解析
*     2. 构造无法读取的文件
* 预期结果：
*     1. 返回解析报错
*     2. 返回读取报错
* 修改历史：
*     1. 2025-8-13 新建测试用例
 */
func TestResponseResponseToStruct(t *testing.T) {
	type fields struct {
		response *http.Response
		err      error
	}
	type args struct {
		i interface{}
	}
	var x string
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "parse failed",
			fields: fields{
				response: &http.Response{Body: &fakeReader{bytes.NewReader([]byte("haha")), false}},
			},
			args: args{i: x},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err)
			},
		},
		{
			name: "read failed",
			fields: fields{
				response: &http.Response{Body: &fakeReader{bytes.NewReader([]byte("haha")), true}},
			},
			args: args{i: x},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &response{
				response: tt.fields.response,
				err:      tt.fields.err,
			}
			err := r.ResponseToStruct(tt.args.i)
			tt.wantErr(t, err)
		})
	}
}

/*
* 测试用例描述：Test_response_ResponseToStruct
* 预置条件：
* 操作步骤：
*     1. 构造一个response
* 预期结果：
*     1. 返回解析结果
* 修改历史：
*     1. 2025-8-13 新建测试用例
 */
func TestResponseResponseToWriter(t *testing.T) {
	type fields struct {
		response *http.Response
		err      error
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		want       int64
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "parse failed",
			fields: fields{
				response: &http.Response{Body: &fakeReader{bytes.NewReader([]byte("haha")), false}},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
			wantWriter: "haha",
			want:       4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &response{
				response: tt.fields.response,
				err:      tt.fields.err,
			}
			writer := &bytes.Buffer{}
			got, err := r.ResponseToWriter(writer)
			if !tt.wantErr(t, err, fmt.Sprintf("ResponseToWriter(%v)", writer)) {
				return
			}
			assert.Equalf(t, tt.wantWriter, writer.String(), "ResponseToWriter(%v)", writer)
			assert.Equalf(t, tt.want, got, "ResponseToWriter(%v)", writer)
		})
	}
}

/*
* 测试用例描述：Test_request_shouldRetry
* 预置条件：
* 操作步骤：
*     1. 请求正常响应
*     2. 请求返回429
*     3. 请求返回502
*     4. 请求返回503
*     5. 请求返回504
*     6. connection refused错误
*     7. connection reset错误
*     8. io.eof 错误
*     9. timeout错误
* 预期结果：
*     1. shouldRetry返回false
*     2. shouldRetry返回true
*     3. shouldRetry返回true
*     4. shouldRetry返回true
*     5. shouldRetry返回true
*     6. shouldRetry返回true
*     7. shouldRetry返回true
*     8. shouldRetry返回true
*     9. shouldRetry返回true
* 修改历史：
*     1. 2025-12-13 新建测试用例
 */
func TestRequestShouldRetry(t *testing.T) {

	tests := []struct {
		name string
		resp *http.Response
		err  error
		want bool
	}{
		{name: "200", resp: &http.Response{StatusCode: http.StatusOK}, err: nil, want: false},
		{name: "429", resp: &http.Response{StatusCode: http.StatusTooManyRequests}, err: nil, want: true},
		{name: "502", resp: &http.Response{StatusCode: http.StatusBadGateway}, err: nil, want: true},
		{name: "503", resp: &http.Response{StatusCode: http.StatusServiceUnavailable}, err: nil, want: true},
		{name: "504", resp: &http.Response{StatusCode: http.StatusGatewayTimeout}, err: nil, want: true},
		{name: "connection refused",
			resp: &http.Response{StatusCode: http.StatusOK}, err: syscall.ECONNREFUSED, want: true},
		{name: "connection reset",
			resp: &http.Response{StatusCode: http.StatusOK}, err: syscall.ECONNRESET, want: true},
		{name: "io.EOF error",
			resp: &http.Response{StatusCode: http.StatusOK}, err: io.EOF, want: true},
		{name: "netErr.Timeout error",
			resp: &http.Response{StatusCode: http.StatusOK}, err: syscall.ETIMEDOUT, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &request{}
			assert.Equalf(t, tt.want, r.shouldRetry(tt.resp, tt.err), "shouldRetry(%v, %v)", tt.resp, tt.err)
		})
	}
}
