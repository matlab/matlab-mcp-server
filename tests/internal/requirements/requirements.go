// Copyright 2026 The MathWorks, Inc.

package requirements

import "context"

type TestingT interface {
	Context() context.Context
	Fatal(args ...any)
	Helper()
	Skip(args ...any)
}
