// Copyright 2025-2026 The MathWorks, Inc.

package globalmatlab_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab"
	"github.com/stretchr/testify/assert"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockMATLABManagerAdaptor := &mocks.MockMATLABManagerAdaptor{}
	defer mockMATLABManagerAdaptor.AssertExpectations(t)

	// Act
	globalMATLAB := globalmatlab.New(mockMATLABManagerAdaptor)

	// Assert
	assert.NotNil(t, globalMATLAB)
}
