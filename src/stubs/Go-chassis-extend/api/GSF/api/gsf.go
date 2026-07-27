package api

import (
	gceBase "Go-chassis-extend/api/base"
	"Go-chassis-extend/api/GSF/api/base"
	"Go-chassis-extend/api/ServiceComb/go-chassis/core/bundle/dependency"
)

type Registry interface {
	GetAllMicroServiceInstanceInfo(selfServiceID string, msKey base.MicroServiceKey) ([]base.MicroServiceInstance, error)
	WatchMicroServiceV1(selfServiceID string, msKeys []base.MicroServiceKey, callback base.WatchServiceNotify, opts ...dependency.WatchOption) error
	UpdateMicroServiceInstanceProperties(selfServiceID, selfInstanceID string, properties map[string]string) error
}

type CspInitOption func()

func WithLocation(loc *gceBase.Location) CspInitOption { return func() {} }

func CspInit(opts ...CspInitOption) error { return nil }
func CspStart(opts ...CspInitOption)       {}
func RegistExitHandler(handler interface{}) {}
func HealthCheckStart(protocol string)     {}

func NewRegistry() Registry { return &mockRegistry{} }

type mockRegistry struct{}

func (r *mockRegistry) GetAllMicroServiceInstanceInfo(selfServiceID string, msKey base.MicroServiceKey) ([]base.MicroServiceInstance, error) {
	return []base.MicroServiceInstance{}, nil
}
func (r *mockRegistry) WatchMicroServiceV1(selfServiceID string, msKeys []base.MicroServiceKey, callback base.WatchServiceNotify, opts ...dependency.WatchOption) error {
	return nil
}
func (r *mockRegistry) UpdateMicroServiceInstanceProperties(selfServiceID, selfInstanceID string, properties map[string]string) error {
	return nil
}