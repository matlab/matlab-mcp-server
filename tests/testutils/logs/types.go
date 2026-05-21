// Copyright 2026 The MathWorks, Inc.

package logs

type DumpPattern struct {
	Glob   string
	Header string
}

type DumpEntry struct {
	Header  string
	File    string
	Content string
}
