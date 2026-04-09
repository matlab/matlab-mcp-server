// Copyright 2025-2026 The MathWorks, Inc.

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
)

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewClientForSelfSignedTLSServer(certificatePEM []byte) (HttpClient, error) {
	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM(certificatePEM); !ok {
		return nil, fmt.Errorf("failed to append certificate to pool")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		},
	}

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &http.Client{
		Transport: transport,
		Jar:       jar,
	}, nil
}

func (f *Factory) NewClientOverUDS(socketPath string) HttpClient {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialSocket(ctx, socketPath)
		},
	}

	return &http.Client{
		Transport: transport,
	}
}

// dialSocket dials either a Unix socket or a Windows named pipe
func dialSocket(ctx context.Context, socketPath string) (net.Conn, error) {
	// Check if this is a Windows named pipe path
	if isNamedPipePath(socketPath) {
		return dialNamedPipe(socketPath)
	}
	// Otherwise use Unix socket
	var d net.Dialer
	return d.DialContext(ctx, "unix", socketPath)
}

// isNamedPipePath checks if the path is a Windows named pipe path
func isNamedPipePath(path string) bool {
	return len(path) >= 9 && path[:9] == `\\.\pipe\`
}
