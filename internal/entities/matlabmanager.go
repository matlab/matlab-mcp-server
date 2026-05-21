// Copyright 2025-2026 The MathWorks, Inc.

package entities

import "context"

type MATLABSessionClient interface {
	Eval(ctx context.Context, sessionLogger Logger, request EvalRequest) (EvalResponse, error)
	EvalWithCapture(ctx context.Context, logger Logger, input EvalRequest) (EvalResponse, error)
	FEval(ctx context.Context, sessionLogger Logger, request FEvalRequest) (FEvalResponse, error)
	Ping(ctx context.Context, sessionLogger Logger) PingResponse
}

type MATLABManager interface {
	ListEnvironments(ctx context.Context, sessionLogger Logger) []EnvironmentInfo
	StartMATLABSession(ctx context.Context, sessionLogger Logger, startRequest SessionDetails) (SessionID, error)
	StopMATLABSession(ctx context.Context, sessionLogger Logger, sessionID SessionID) error
	GetMATLABSessionClient(ctx context.Context, sessionLogger Logger, sessionID SessionID) (MATLABSessionClient, error)
}

type EnvironmentInfo struct {
	MATLABRoot string
	Version    string
}

type SessionID int

// SessionDetails is an interface to disambiguate which type of MATLAB session to start.
type SessionDetails interface {
	interfacelock()
}

type LocalSessionDetails struct {
	MATLABRoot             string
	IsStartingDirectorySet bool
	StartingDirectory      string
	ShowMATLABDesktop      bool
}

func (l LocalSessionDetails) interfacelock() {}

type AttachToExistingSession struct{}

func (a AttachToExistingSession) interfacelock() {}

type EvalRequest struct {
	Code string
}

type EvalResponse struct {
	ConsoleOutput string
	Images        [][]byte
}

type FEvalRequest struct {
	Function   string
	Arguments  []string
	NumOutputs int
}

type FEvalResponse struct {
	Outputs []any
}

type PingResponse struct {
	IsAlive bool
}
