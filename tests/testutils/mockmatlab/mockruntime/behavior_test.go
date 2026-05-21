// Copyright 2026 The MathWorks, Inc.

package mockruntime_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
	"github.com/stretchr/testify/assert"
)

func TestShouldExitFromEvalCode(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{name: "exit command", code: "exit", want: true},
		{name: "exit function", code: "exit()", want: true},
		{name: "quit command", code: "quit", want: true},
		{name: "quit function", code: "quit()", want: true},
		{name: "non-exit code", code: "disp('hello')", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := mockruntime.ShouldExitFromEvalCode(tc.code)

			// Assert
			assert.Equal(t, tc.want, got)
		})
	}
}
