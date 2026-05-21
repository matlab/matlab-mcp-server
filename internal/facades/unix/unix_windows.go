// Copyright 2026 The MathWorks, Inc.
//go:build windows

package unix

// Rlimit struct mirrors unix.Rlimit for cross-platform compatibility.

type Rlimit struct {
	Cur uint64
	Max uint64
}

func (uf *UnixFacade) Getrlimit(_ int, _ *Rlimit) error {
	return nil
}

func (uf *UnixFacade) Setrlimit(_ int, _ *Rlimit) error {
	return nil
}
