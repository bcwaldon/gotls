package gotls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type TLSInfo struct {
	CertFile string
	KeyFile  string
	CAFile   string

	// parseFunc exists to simplify testing. Typically, parseFunc
	// should be left nil. In that case, tls.X509KeyPair will be used.
	parseFunc func([]byte, []byte) (tls.Certificate, error)
}

func (info TLSInfo) String() string {
	return fmt.Sprintf("cert = %s, key = %s, ca = %s", info.CertFile, info.KeyFile, info.CAFile)
}

func (info TLSInfo) Empty() bool {
	return info.CertFile == "" && info.KeyFile == ""
}

func (info TLSInfo) baseConfig() (*tls.Config, error) {
	if info.KeyFile == "" || info.CertFile == "" {
		return nil, fmt.Errorf("KeyFile and CertFile must both be present[key: %v, cert: %v]", info.KeyFile, info.CertFile)
	}

	cert, err := ioutil.ReadFile(info.CertFile)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(info.KeyFile)
	if err != nil {
		return nil, err
	}

	parseFunc := info.parseFunc
	if parseFunc == nil {
		parseFunc = tls.X509KeyPair
	}

	tlsCert, err := parseFunc(cert, key)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS10,
	}
	return cfg, nil
}

// ServerConfig generates a tls.Config object for use by an HTTP server
func (info TLSInfo) ServerConfig() (*tls.Config, error) {
	cfg, err := info.baseConfig()
	if err != nil {
		return nil, err
	}

	if info.CAFile != "" {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
		cp, err := newCertPool(info.CAFile)
		if err != nil {
			return nil, err
		}
		cfg.ClientCAs = cp
	} else {
		cfg.ClientAuth = tls.NoClientCert
	}

	return cfg, nil
}

// ClientConfig generates a tls.Config object for use by an HTTP client
func (info TLSInfo) ClientConfig() (cfg *tls.Config, err error) {
	if !info.Empty() {
		cfg, err = info.baseConfig()
		if err != nil {
			return nil, err
		}
	} else {
		cfg = &tls.Config{}
	}

	if info.CAFile != "" {
		cfg.RootCAs, err = newCertPool(info.CAFile)
		if err != nil {
			return
		}
	}

	return
}

// newCertPool creates x509 certPool with provided CA file
func newCertPool(CAFile string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	pemByte, err := ioutil.ReadFile(CAFile)
	if err != nil {
		return nil, err
	}

	for {
		var block *pem.Block
		block, pemByte = pem.Decode(pemByte)
		if block == nil {
			return certPool, nil
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certPool.AddCert(cert)
	}
}
