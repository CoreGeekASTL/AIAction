/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2026. All rights reserved.
 */

package https

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"GIDS/common/logger"
)

const (
	ServerType = "server"
	ClientType = "client"
)

type CertInfo struct {
	KeyFile  string
	CertFile string
	CaFile   string
	KeyPwd   []byte
}

func GetTLS(info CertInfo, tlsType string) *tls.Config {
	tlsConfig := &tls.Config{
		NextProtos: []string{"h2", "http/1.1"},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	caPool, certs, err := getCert(info)
	if err != nil {
		logger.Errorf("get cert info err, err is %v", err)
	}
	if certs != nil {
		tlsConfig.Certificates = []tls.Certificate{*certs}
	}

	switch tlsType {
	case ServerType:
		tlsConfig.ClientCAs = caPool
	case ClientType:
		v := verify{
			rootCA: caPool,
		}
		tlsConfig.RootCAs = caPool
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.VerifyConnection = v.verifyConnection
	default:
		logger.Infof("unknown tls type: %s", tlsType)
	}
	return tlsConfig
}

func getCert(info CertInfo) (*x509.CertPool, *tls.Certificate, error) {
	tlsCrt, err := defaultLoadFile(info.CertFile)
	if err != nil {
		logger.Errorf("[tlsConfig] failed to read cert file, error: %v", err)
		return nil, nil, err
	}
	tlsKey, err := defaultLoadFile(info.KeyFile)
	if err != nil {
		logger.Errorf("[tlsConfig] failed to read key file, error: %v", err)
		return nil, nil, err
	}
	if len(info.KeyPwd) != 0 {
		privBlock, _ := pem.Decode(tlsKey)
		if privBlock == nil {
			logger.Errorf("[tlsConfig] failed to parse private key")
		}
		decryptedBlock, err := x509.DecryptPEMBlock(privBlock, info.KeyPwd)
		if err != nil {
			logger.Errorf("[tlsConfig] decrypt private key failed, err is : %v", err)
		}
		tlsKey = pem.EncodeToMemory(&pem.Block{
			Type:  privBlock.Type,
			Bytes: decryptedBlock,
		})
	}

	cert, err := tls.X509KeyPair(tlsCrt, tlsKey)
	if err != nil {
		return nil, nil, fmt.Errorf("[tlsConfig] generate certificate failed: %v", err)
	}

	caCrt, err := defaultLoadFile(info.CaFile)
	if err != nil {
		logger.Errorf("[tlsConfig] failed to read ca file, error: %v", err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)

	return pool, &cert, nil
}

func defaultLoadFile(filePath string) ([]byte, error) {
	relPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(relPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type verify struct {
	rootCA *x509.CertPool
}

func (v *verify) verifyConnection(cs tls.ConnectionState) error {
	if len(cs.PeerCertificates) == 0 {
		return fmt.Errorf("no peer certificates")
	}

	// 获取服务器证书
	cert := cs.PeerCertificates[0]

	if err := v.checkValidity(cert); err != nil {
		return err
	}
	if err := v.checkBasicConstraints(cert); err != nil {
		return err
	}
	if err := v.checkSignatureAlgorithm(cert); err != nil {
		return err
	}
	if err := v.checkKeyUsage(cert); err != nil {
		return err
	}

	if len(cert.DNSNames) > 0 {
		if err := v.verifyHostname(cert, cs.ServerName); err != nil {
			return err
		}
	}

	return nil
}

// checkValidity 验证证书有效期
func (v *verify) checkValidity(cert *x509.Certificate) error {
	now := time.Now()

	if now.Before(cert.NotBefore) {
		return fmt.Errorf("certificate has not yet taken effect (effective date: %s)", cert.NotBefore)
	}

	if now.After(cert.NotAfter) {
		return fmt.Errorf("certificate has expired (expiration date: %s)", cert.NotAfter)
	}

	return nil
}

// checkKeyUsage 验证密钥用法
func (v *verify) checkKeyUsage(cert *x509.Certificate) error {
	// 检查不允许的用法
	if cert.KeyUsage&x509.KeyUsageCertSign != 0 {
		// 服务器证书不应该允许证书签名
		return fmt.Errorf("cert should not allow certificate signing")
	}

	if cert.KeyUsage&x509.KeyUsageCRLSign != 0 {
		// 服务器证书不应该允许 CRL 签名
		return fmt.Errorf("certificate should not allow CRL signature")
	}
	return nil
}

// checkBasicConstraints 验证基本约束
func (v *verify) checkBasicConstraints(cert *x509.Certificate) error {
	// 服务器证书不应该是 CA 证书
	if cert.IsCA {
		return fmt.Errorf("server certificate cannot be a CA certificate")
	}

	return nil
}

// checkSignatureAlgorithm 验证签名算法
func (v *verify) checkSignatureAlgorithm(cert *x509.Certificate) error {
	algorithm := cert.SignatureAlgorithm.String()

	// 检查弱签名算法
	weakAlgorithms := map[x509.SignatureAlgorithm]bool{
		x509.MD2WithRSA:    true,
		x509.MD5WithRSA:    true,
		x509.DSAWithSHA1:   true,
		x509.ECDSAWithSHA1: true,
		x509.SHA1WithRSA:   true,
	}

	if weakAlgorithms[cert.SignatureAlgorithm] {
		return fmt.Errorf("signature algorithm not allowed: %s", algorithm)
	}

	return nil
}

// verifyCertificateChain 验证证书链
func (v *verify) verifyCertificateChain(cert *x509.Certificate) error {
	opts := x509.VerifyOptions{
		Roots:       v.rootCA,
		CurrentTime: time.Now(),
		KeyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	if _, err := cert.Verify(opts); err != nil {
		return err
	}

	return nil
}

// verifyHostname 验证主机名
func (v *verify) verifyHostname(cert *x509.Certificate, hostname string) error {
	// 1. 首先检查 SAN
	for _, dnsName := range cert.DNSNames {
		if matchHostname(dnsName, hostname) {
			logger.Infof("cert match hostname %v successful", hostname)
			return nil
		}
	}

	// 2. 检查 IP 地址
	if ip := net.ParseIP(hostname); ip != nil {
		for _, certIP := range cert.IPAddresses {
			if certIP.Equal(ip) {
				logger.Infof("cert match hostip %v successful", certIP)
				return nil
			}
		}
	}

	return fmt.Errorf("host %s none of the names in the certificate match", hostname)
}

// matchHostname 主机名匹配函数
func matchHostname(pattern, hostname string) bool {
	pattern = strings.ToLower(pattern)
	hostname = strings.ToLower(strings.TrimSuffix(hostname, "."))

	if len(pattern) == 0 || len(hostname) == 0 {
		return false
	}

	patternParts := strings.Split(pattern, ".")
	hostParts := strings.Split(hostname, ".")

	if len(patternParts) != len(hostParts) {
		return false
	}

	for i, patternPart := range patternParts {
		if patternPart == "*" {
			continue
		}
		if patternPart != hostParts[i] {
			return false
		}
	}
	return true
}
