package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
)

func LoadTLSCreds(hostCrtPath string, privateKeyPath string, caCertPath string) (tls.Certificate, *x509.CertPool, error) {
	emptyCert := tls.Certificate{}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return emptyCert, nil, fmt.Errorf("failed to get current path")
	}

	hostCrt := path.Join(path.Dir(file), hostCrtPath)
	privateKey := path.Join(path.Dir(file), privateKeyPath)
	caCert := path.Join(path.Dir(file), caCertPath)

	certificate, err := tls.LoadX509KeyPair(hostCrt, privateKey)
	if err != nil {
		return emptyCert, nil, err
	}

	certPool := x509.NewCertPool()
	cert, err := ioutil.ReadFile(caCert)
	if err != nil {
		return emptyCert, nil, err
	}

	ok = certPool.AppendCertsFromPEM(cert)
	if !ok {
		return emptyCert, nil, fmt.Errorf("failed to append ca certificate")
	}

	return certificate, certPool, nil
}
