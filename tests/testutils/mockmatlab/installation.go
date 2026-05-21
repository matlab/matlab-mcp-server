// Copyright 2026 The MathWorks, Inc.

package mockmatlab

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/tests/testutils/facades/filefacade"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/installer"
	"github.com/stretchr/testify/require"
)

type Installation struct {
	MATLABRoot string
}

func BuildAndInstall(t *testing.T) *Installation {
	t.Helper()
	matlabRoot := t.TempDir()
	inst := installer.New(
		filefacade.RealFileSystem{},
		installer.GoListModuleRootFinder{},
		installer.BinaryBuilderFunc(buildPlatformSpecificBinaries),
	)
	require.NoError(t, inst.BuildAndInstall(matlabRoot))

	return &Installation{MATLABRoot: matlabRoot}
}
