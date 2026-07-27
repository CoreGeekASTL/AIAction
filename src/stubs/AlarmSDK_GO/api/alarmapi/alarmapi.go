package alarmapi

import "AlarmSDK_GO/api/base"

func CSPInitAlarmSDK(appID, serviceName, nodeName, nodeIP string) base.CSPAlarmManager {
	return &mockAlarmManager{}
}
func GetNodeIP() string { return "127.0.0.1" }
func InitCSPAlarm(alarmID string, alarmType base.GenerateOrClearType) base.CSPAlarm {
	return &mockAlarm{params: make(map[string]string)}
}

type mockAlarmManager struct{}
func (m *mockAlarmManager) SendAlarm(alarm base.CSPAlarm) bool { return true }
func (m *mockAlarmManager) RegisterRsetClear()                 {}

type mockAlarm struct {
	params map[string]string
}

func (a *mockAlarm) AppendParameter(key, value string) {
	a.params[key] = value
}