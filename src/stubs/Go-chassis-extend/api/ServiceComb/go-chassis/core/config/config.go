package config

type GlobalDefinition struct {
	AppID string
}

func GetSelfServiceID() string                     { return "" }
func GetSelfInstanceID() string                    { return "" }
func GetSelfServiceName() string                   { return "" }
func GetGlobalDefinition() *GlobalDefinition       { return &GlobalDefinition{} }