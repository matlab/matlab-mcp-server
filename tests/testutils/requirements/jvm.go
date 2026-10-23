// Copyright 2026 The MathWorks, Inc.

package requirements

import (
	"context"
	"slices"
	"strings"
)

const (
	jvmForceEnvVar       = "MCP_SYSTEM_TEST_FORCE_JVM"
	jvmCheckCode         = `fprintf("MCP_USEJAVA_JVM=%d\n", usejava('jvm'));`
	jvmUnavailableReason = `usejava('jvm') is false: the test needs Java, and this MATLAB has no Java it can use (from R2026b, MATLAB does not include a JRE)`
)

// JVM skips the test when the MATLAB under test has no JVM.
func JVM(ctx context.Context, t TestingT, session MATLABSession) {
	t.Helper()

	output, err := session.EvaluateCode(ctx, jvmCheckCode)
	if err != nil {
		t.Fatalf(`usejava('jvm') should run in the MATLAB under test: %v`, err)
		return
	}

	fields := strings.Fields(output)
	switch {
	case slices.Contains(fields, "MCP_USEJAVA_JVM=1"):
	case slices.Contains(fields, "MCP_USEJAVA_JVM=0"):
		unavailable(t, jvmForceEnvVar, jvmUnavailableReason)
	default:
		t.Fatalf(`usejava('jvm') should print MCP_USEJAVA_JVM=0 or 1, got %q`, output)
	}
}
