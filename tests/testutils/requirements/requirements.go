// Copyright 2026 The MathWorks, Inc.

package requirements

import (
	"cmp"
	"context"
	"os"
	"strconv"
)

type TestingT interface {
	Helper()

	Logf(format string, args ...any)
	Skipf(format string, args ...any)
	Fatalf(format string, args ...any)
}

type MATLABSession interface {
	EvaluateCode(ctx context.Context, code string, projectPath ...string) (string, error)
}

// unavailable skips the test, unless forceEnvVar is 1, which runs the test anyway so QE can still see how it fails.
func unavailable(t TestingT, forceEnvVar string, reason string) {
	t.Helper()

	value := cmp.Or(os.Getenv(forceEnvVar), "0")
	forced, err := strconv.ParseBool(value)
	if err != nil {
		t.Fatalf("%s should be 1 or 0, got %q", forceEnvVar, value)
		return
	}

	if forced {
		t.Logf("%s. Running anyway because %s is set", reason, forceEnvVar)
		return
	}

	t.Skipf("%s. Set %s=1 to run the test anyway", reason, forceEnvVar)
}
