package rest

import "context"

type Request struct{}
type Response struct{}

func NewRequest(method, url string, body []byte) (*Request, error) {
	return &Request{}, nil
}

func (r *Request) SetHeader(key, value string) {}
func (r *Request) GetHeader(key string) string { return "" }
func (r *Request) SetBody(body []byte)         {}
func (r *Request) GetBody() []byte             { return nil }
func (r *Request) Close()                      {}

func (r *Response) GetStatusCode() int { return 200 }
func (r *Response) ReadBody() []byte   { return nil }
func (r *Response) Close()             {}

type RestInvoker struct{}

func NewRestInvoker() *RestInvoker { return &RestInvoker{} }

func (i *RestInvoker) Do(ctx context.Context, req *Request) (*Response, error) {
	return &Response{}, nil
}

func (i *RestInvoker) ContextDo(ctx context.Context, req *Request) (*Response, error) {
	return &Response{}, nil
}