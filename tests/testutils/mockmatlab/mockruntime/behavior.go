// Copyright 2026 The MathWorks, Inc.

package mockruntime

func ShouldExitFromEvalCode(code string) bool {
	switch code {
	case "exit()", "exit", "quit", "quit()":
		return true
	default:
		return false
	}
}
