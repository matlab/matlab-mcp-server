// Copyright 2026 The MathWorks, Inc.
//go:build !windows

package unix

import "golang.org/x/sys/unix"

type Rlimit = unix.Rlimit

func (uf *UnixFacade) Getrlimit(resource int, rlim *Rlimit) error {
	return unix.Getrlimit(resource, rlim)
}

func (uf *UnixFacade) Setrlimit(resource int, rlim *Rlimit) error {
	return unix.Setrlimit(resource, rlim)
}
