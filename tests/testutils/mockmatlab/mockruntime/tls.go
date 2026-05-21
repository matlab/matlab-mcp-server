// Copyright 2026 The MathWorks, Inc.

package mockruntime

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

func (r *Runtime) GenerateAndWriteCert(certFile, keyFile string) (*tls.Config, error) {
	certPEM, keyPEM, err := r.TLSProvider.GeneratePEM()
	if err != nil {
		return nil, err
	}

	if err := r.FS.WriteFile(certFile, certPEM, 0o600); err != nil {
		return nil, fmt.Errorf("failed to write cert file: %w", err)
	}
	if err := r.FS.WriteFile(keyFile, keyPEM, 0o600); err != nil {
		return nil, fmt.Errorf("failed to write key file: %w", err)
	}

	tlsConfig, err := r.TLSProvider.TLSConfig(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return tlsConfig, nil
}

type defaultTLSMaterialProvider struct{}

func (defaultTLSMaterialProvider) GeneratePEM() ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
		DNSNames:              []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return certPEM, keyPEM, nil
}

func (defaultTLSMaterialProvider) TLSConfig(certPEM []byte, keyPEM []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load X509 key pair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
