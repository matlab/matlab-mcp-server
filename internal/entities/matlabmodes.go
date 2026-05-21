// Copyright 2026 The MathWorks, Inc.

package entities

type DisplayMode string

const (
	DisplayModeDesktop   DisplayMode = "desktop"
	DisplayModeNoDesktop DisplayMode = "nodesktop"
)

type MATLABSessionMode string

const (
	MATLABSessionModeNew      MATLABSessionMode = "new"
	MATLABSessionModeExisting MATLABSessionMode = "existing"
)
