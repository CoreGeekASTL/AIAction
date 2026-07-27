package monitor

type MonSdk struct{}

func (m *MonSdk) InitMonitor(appID int, serviceName, instanceName string) error { return nil }
func (m *MonSdk) RegisterBasicInfo(monitorJSON string) error                     { return nil }
func (m *MonSdk) SetMetric(metricID int, obj string, value float64) error        { return nil }
func (m *MonSdk) ObjChange(mocID uint32, changeType int, moiID string) error     { return nil }

var MonSdkInstance = &MonSdk{}