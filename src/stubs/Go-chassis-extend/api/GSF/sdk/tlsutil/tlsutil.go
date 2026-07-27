package tlsutil

import "crypto/tls"

func GetTLSConfig(serviceName, schema string, role string) (*tls.Config, error) {
	return &tls.Config{InsecureSkipVerify: true}, nil
}