package api

import (
	"net/http"
)

type CspRestInvoker struct{}
type CspResponse struct{}

func NewCspRestInvoker() *CspRestInvoker { return &CspRestInvoker{} }

func (i *CspRestInvoker) Invoke(method, url string, headers http.Header, body []byte) (*CspResponse, error) {
	return &CspResponse{}, nil
}

func (r *CspResponse) GetStatusCode() int { return http.StatusOK }
func (r *CspResponse) ReadBody() []byte   { return nil }
func (r *CspResponse) Close()             {}