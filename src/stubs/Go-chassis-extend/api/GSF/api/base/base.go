package base

type MicroServiceKey struct {
	AppId       string
	ServiceName string
	Version     string
}

type MicroServiceInstance struct {
	InstanceId    string
	Properties    map[string]string
	HostName      string
	EndPointMap   map[string]string
	Status        string
	EndpointsList []string
}

type MicroServiceInstanceWrapper struct {
	Instance *MicroServiceInstance
}

type MicroServiceInstanceChangedEvent struct {
	Action       string
	Instance     *MicroServiceInstance
	InstanceList []*MicroServiceInstanceWrapper
}

type WatchNotifier interface {
	WatchServiceCallBack(event *MicroServiceInstanceChangedEvent)
}

type MicroService struct {
	ServiceId   string
	AppId       string
	ServiceName string
	Version     string
}

type WatchServiceNotify interface {
	WatchServiceCallBack(event *MicroServiceInstanceChangedEvent)
}

type WatchLabelNotify interface {
	WatchLabelCallBack(event *MicroServiceInstanceChangedEvent)
}

const RestProtocal = "rest"