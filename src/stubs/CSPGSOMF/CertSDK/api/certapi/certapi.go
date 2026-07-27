package certapi

import "CSPGSOMF/CertSDK/api/base"

func CspCertSDKInit() error { return nil }
func GetExCertManagerInstance() base.CSPExCertManager {
	return &mockCertManager{}
}

type mockCertManager struct{}

func (m *mockCertManager) SubscribeExCert(appName string, scenes []base.CspExSceneInfo, handler func([]*base.CspExCertInfo, int) error, path string) error {
	return nil
}
func (m *mockCertManager) UnsubscribeExCert(appName string, scenes []base.CspExSceneInfo) error {
	return nil
}
func (m *mockCertManager) GetExCertPathInfo(sceneName string) (*base.CspExCertPathInfo, error) {
	return &base.CspExCertPathInfo{}, nil
}
func (m *mockCertManager) GetExCertPrivateKeyPwd(sceneName string) ([]byte, error) {
	return []byte{}, nil
}