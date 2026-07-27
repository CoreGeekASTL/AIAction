package core

import (
	"context"
	"Go-chassis-extend/api/ServiceComb/go-chassis/client/rest"
)

type RestInvoker struct{}

func NewRestInvoker() *RestInvoker { return &RestInvoker{} }

func (i *RestInvoker) ContextDo(ctx context.Context, req *rest.Request) (*rest.Response, error) {
	return &rest.Response{}, nil
}