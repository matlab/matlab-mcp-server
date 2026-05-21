// Copyright 2026 The MathWorks, Inc.

package mockruntime

import (
	"crypto/tls"
	"os"
)

type Environment interface {
	Getenv(key string) string
}

type FileSystem interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type TLSMaterialProvider interface {
	GeneratePEM() (certPEM []byte, keyPEM []byte, err error)
	TLSConfig(certPEM []byte, keyPEM []byte) (*tls.Config, error)
}

type Runtime struct {
	Env         Environment
	FS          FileSystem
	TLSProvider TLSMaterialProvider
}

func NewRuntime(env Environment, fileSystem FileSystem, tlsProvider TLSMaterialProvider) *Runtime {
	return &Runtime{
		Env:         env,
		FS:          fileSystem,
		TLSProvider: tlsProvider,
	}
}

func NewDefaultTLSMaterialProvider() TLSMaterialProvider {
	return defaultTLSMaterialProvider{}
}
