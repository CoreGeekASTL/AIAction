package base

type CSPExCertManager interface {
	SubscribeExCert(appName string, scenes []CspExSceneInfo, handler func(certInfo []*CspExCertInfo, notifyType int) error, path string) error
	UnsubscribeExCert(appName string, scenes []CspExSceneInfo) error
	GetExCertPathInfo(sceneName string) (*CspExCertPathInfo, error)
	GetExCertPrivateKeyPwd(sceneName string) ([]byte, error)
}

type CspExSceneInfo struct {
	SceneName   string
	SceneDescCN string
	SceneDescEN string
	SceneType   int
	Feature     int
}

type CspExCertPathInfo struct {
	ExCaFilePath         string
	ExDeviceFilePath     string
	ExPrivateKeyFilePath string
}

type CspExCertInfo struct {
	SceneName string
}