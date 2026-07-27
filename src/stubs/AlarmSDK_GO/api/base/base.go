package base

type CSPAlarmManager interface {
	SendAlarm(alarm CSPAlarm) bool
	RegisterRsetClear()
}

type CSPAlarm interface {
	AppendParameter(key, value string)
}

type GenerateOrClearType int

const (
	GenerateAlarm GenerateOrClearType = 0
	ClearAlarm    GenerateOrClearType = 1
)