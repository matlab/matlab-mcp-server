// Copyright 2026 The MathWorks, Inc.

package fileuri_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-server/tests/testutils/fileuri"
	"github.com/stretchr/testify/assert"
)

func TestToFileURI_AbsolutePath(t *testing.T) {
	// Act
	result := fileuri.ToFileURI("/home/user/data")

	// Assert
	assert.Equal(t, "file:///home/user/data", result)
}

func TestToFileURI_AddsLeadingSlashForDriveLetterPath(t *testing.T) {
	// Act
	result := fileuri.ToFileURI("C:/Users/test")

	// Assert
	assert.Equal(t, "file:///C:/Users/test", result)
}
