// Copyright 2026 The MathWorks, Inc.

package requirements

const (
	JVMForceEnvVar       = jvmForceEnvVar
	JVMCheckCode         = jvmCheckCode
	JVMUnavailableReason = jvmUnavailableReason
)

func Unavailable(t TestingT, forceEnvVar string, reason string) {
	unavailable(t, forceEnvVar, reason)
}
