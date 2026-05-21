// Copyright 2026 The MathWorks, Inc.

package custom_test

import (
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	evalcustomtoolusecase "github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	basetoolmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	custommocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/custom"
	definitionmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/singlesession/custom/definition"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_HappyPath(t *testing.T) {
	tests := []struct {
		name                    string
		shouldShowMATLABDesktop bool
		expectedResponse        entities.EvalResponse
		expectedContentCount    int
	}{
		{
			"nodesktop",
			false,
			entities.EvalResponse{
				ConsoleOutput: "    17    24     1     8    15",
				Images:        [][]byte{[]byte("image1")},
			},
			2,
		},
		{
			"desktop",
			true,
			entities.EvalResponse{
				ConsoleOutput: "    17    24     1     8    15",
			},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
			defer mockLoggerFactory.AssertExpectations(t)

			mockConfigFactory := &custommocks.MockConfigFactory{}
			defer mockConfigFactory.AssertExpectations(t)

			mockConfig := &configmocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockUsecase := &custommocks.MockUsecase{}
			defer mockUsecase.AssertExpectations(t)

			mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
			defer mockGlobalMATLAB.AssertExpectations(t)

			mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
			defer mockMATLABSessionClient.AssertExpectations(t)

			mockSessionLogger := testutils.NewInspectableLogger()
			ctx := t.Context()
			expectedSession := &mcp.ServerSession{}
			mockValidatedTool := &definitionmocks.MockValidatedTool{}
			defer mockValidatedTool.AssertExpectations(t)

			expectedDefinition := definition.Tool{
				Name:        "generate_magic_square",
				Title:       "Generate Magic Square",
				Description: "Generates a magic square",
				InputSchema: &jsonschema.Schema{
					Type: "object",
					Properties: map[string]*jsonschema.Schema{
						"n": {Type: "number", Description: "Size"},
					},
					Required: []string{"n"},
				},
			}
			expectedSignature := definition.Signature{
				Function: "magic",
				Input:    definition.SignatureInput{Order: []string{"n"}},
			}
			args := map[string]any{"n": float64(5)}
			req := &mcp.CallToolRequest{
				Session: expectedSession,
			}

			mockValidatedTool.EXPECT().
				Definition().
				Return(expectedDefinition).
				Once()
			mockValidatedTool.EXPECT().
				Signature().
				Return(expectedSignature).
				Once()

			mockLoggerFactory.EXPECT().
				NewMCPSessionLogger(expectedSession).
				Return(mockSessionLogger, nil).
				Once()

			mockConfigFactory.EXPECT().
				Config().
				Return(mockConfig, nil).
				Once()

			mockConfig.EXPECT().
				ShouldShowMATLABDesktop().
				Return(tt.shouldShowMATLABDesktop).
				Once()

			mockGlobalMATLAB.EXPECT().
				Client(ctx, mockSessionLogger.AsMockArg()).
				Return(mockMATLABSessionClient, nil).
				Once()

			mockUsecase.EXPECT().
				Execute(
					ctx,
					mockSessionLogger.AsMockArg(),
					mockMATLABSessionClient,
					evalcustomtoolusecase.Args{
						Function:      "magic",
						Order:         []string{"n"},
						ArgumentTypes: map[string]string{"n": "number"},
						Arguments:     args,
						CaptureOutput: !tt.shouldShowMATLABDesktop,
					},
				).
				Return(tt.expectedResponse, nil).
				Once()

			handler := custom.Handler(mockValidatedTool, mockLoggerFactory, mockConfigFactory, mockUsecase, mockGlobalMATLAB)

			// Act
			result, _, err := handler(ctx, req, args)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Len(t, result.Content, tt.expectedContentCount)

			textContent, ok := result.Content[0].(*mcp.TextContent)
			require.True(t, ok)
			assert.Equal(t, tt.expectedResponse.ConsoleOutput, textContent.Text)

			for i, img := range tt.expectedResponse.Images {
				imageContent, ok := result.Content[i+1].(*mcp.ImageContent)
				require.True(t, ok)
				assert.Equal(t, "image/png", imageContent.MIMEType)
				assert.Equal(t, img, imageContent.Data)
			}
		})
	}
}

func TestHandler_NoArguments_HappyPath(t *testing.T) {
	tests := []struct {
		name                  string
		inputSchema           *jsonschema.Schema
		expectedArgumentTypes map[string]string
	}{
		{
			"empty_schema",
			&jsonschema.Schema{Type: "object"},
			map[string]string{},
		},
		{
			"nil_schema",
			nil,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
			defer mockLoggerFactory.AssertExpectations(t)

			mockConfigFactory := &custommocks.MockConfigFactory{}
			defer mockConfigFactory.AssertExpectations(t)

			mockConfig := &configmocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockUsecase := &custommocks.MockUsecase{}
			defer mockUsecase.AssertExpectations(t)

			mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
			defer mockGlobalMATLAB.AssertExpectations(t)

			mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
			defer mockMATLABSessionClient.AssertExpectations(t)

			mockSessionLogger := testutils.NewInspectableLogger()
			ctx := t.Context()
			expectedSession := &mcp.ServerSession{}
			mockValidatedTool := &definitionmocks.MockValidatedTool{}
			defer mockValidatedTool.AssertExpectations(t)

			expectedDefinition := definition.Tool{
				Name:        "get_version",
				Title:       "Get Version",
				Description: "Returns the MATLAB version",
				InputSchema: tt.inputSchema,
			}
			expectedSignature := definition.Signature{
				Function: "version",
				Input:    definition.SignatureInput{Order: []string{}},
			}
			shouldShowMATLABDesktop := false
			args := map[string]any{}
			expectedResponse := entities.EvalResponse{
				ConsoleOutput: "9.14.0",
			}
			req := &mcp.CallToolRequest{
				Session: expectedSession,
			}

			mockValidatedTool.EXPECT().
				Definition().
				Return(expectedDefinition).
				Once()
			mockValidatedTool.EXPECT().
				Signature().
				Return(expectedSignature).
				Once()

			mockLoggerFactory.EXPECT().
				NewMCPSessionLogger(expectedSession).
				Return(mockSessionLogger, nil).
				Once()

			mockConfigFactory.EXPECT().
				Config().
				Return(mockConfig, nil).
				Once()

			mockConfig.EXPECT().
				ShouldShowMATLABDesktop().
				Return(shouldShowMATLABDesktop).
				Once()

			mockGlobalMATLAB.EXPECT().
				Client(ctx, mockSessionLogger.AsMockArg()).
				Return(mockMATLABSessionClient, nil).
				Once()

			mockUsecase.EXPECT().
				Execute(
					ctx,
					mockSessionLogger.AsMockArg(),
					mockMATLABSessionClient,
					evalcustomtoolusecase.Args{
						Function:      "version",
						Order:         []string{},
						ArgumentTypes: tt.expectedArgumentTypes,
						Arguments:     args,
						CaptureOutput: !shouldShowMATLABDesktop,
					},
				).
				Return(expectedResponse, nil).
				Once()

			handler := custom.Handler(mockValidatedTool, mockLoggerFactory, mockConfigFactory, mockUsecase, mockGlobalMATLAB)

			// Act
			result, _, err := handler(ctx, req, args)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Len(t, result.Content, 1)

			textContent, ok := result.Content[0].(*mcp.TextContent)
			require.True(t, ok)
			assert.Equal(t, expectedResponse.ConsoleOutput, textContent.Text)
		})
	}
}

func TestHandler_ConfigError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &custommocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockUsecase := &custommocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSessionLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedSession := &mcp.ServerSession{}
	mockValidatedTool := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool.AssertExpectations(t)

	expectedDefinition := definition.Tool{
		Name:        "generate_magic_square",
		Title:       "Generate Magic Square",
		Description: "Generates a magic square",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"n": {Type: "number", Description: "Size"},
			},
			Required: []string{"n"},
		},
	}
	args := map[string]any{"n": float64(5)}
	expectedError := messages.AnError
	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	mockValidatedTool.EXPECT().
		Definition().
		Return(expectedDefinition).
		Once()
	mockValidatedTool.EXPECT().
		Signature().
		Return(definition.Signature{}).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(nil, expectedError).
		Once()

	handler := custom.Handler(mockValidatedTool, mockLoggerFactory, mockConfigFactory, mockUsecase, mockGlobalMATLAB)

	// Act
	_, _, err := handler(ctx, req, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestHandler_ClientError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &custommocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockUsecase := &custommocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockSessionLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedSession := &mcp.ServerSession{}
	mockValidatedTool := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool.AssertExpectations(t)

	expectedDefinition := definition.Tool{
		Name:        "generate_magic_square",
		Title:       "Generate Magic Square",
		Description: "Generates a magic square",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"n": {Type: "number", Description: "Size"},
			},
			Required: []string{"n"},
		},
	}
	args := map[string]any{"n": float64(5)}
	expectedError := assert.AnError
	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	mockValidatedTool.EXPECT().
		Definition().
		Return(expectedDefinition).
		Once()
	mockValidatedTool.EXPECT().
		Signature().
		Return(definition.Signature{}).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockSessionLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	handler := custom.Handler(mockValidatedTool, mockLoggerFactory, mockConfigFactory, mockUsecase, mockGlobalMATLAB)

	// Act
	_, _, err := handler(ctx, req, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestHandler_UsecaseError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &custommocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockUsecase := &custommocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockMATLABSessionClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockMATLABSessionClient.AssertExpectations(t)

	mockSessionLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedSession := &mcp.ServerSession{}
	mockValidatedTool := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool.AssertExpectations(t)

	expectedDefinition := definition.Tool{
		Name:        "generate_magic_square",
		Title:       "Generate Magic Square",
		Description: "Generates a magic square",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"n": {Type: "number", Description: "Size"},
			},
			Required: []string{"n"},
		},
	}
	expectedSignature := definition.Signature{
		Function: "magic",
		Input:    definition.SignatureInput{Order: []string{"n"}},
	}
	shouldShowMATLABDesktop := false
	args := map[string]any{"n": float64(5)}
	expectedError := assert.AnError
	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	mockValidatedTool.EXPECT().
		Definition().
		Return(expectedDefinition).
		Once()
	mockValidatedTool.EXPECT().
		Signature().
		Return(expectedSignature).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		ShouldShowMATLABDesktop().
		Return(shouldShowMATLABDesktop).
		Once()

	mockGlobalMATLAB.EXPECT().
		Client(ctx, mockSessionLogger.AsMockArg()).
		Return(mockMATLABSessionClient, nil).
		Once()

	mockUsecase.EXPECT().
		Execute(
			ctx,
			mockSessionLogger.AsMockArg(),
			mockMATLABSessionClient,
			evalcustomtoolusecase.Args{
				Function:      "magic",
				Order:         []string{"n"},
				ArgumentTypes: map[string]string{"n": "number"},
				Arguments:     args,
				CaptureOutput: !shouldShowMATLABDesktop,
			},
		).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	handler := custom.Handler(mockValidatedTool, mockLoggerFactory, mockConfigFactory, mockUsecase, mockGlobalMATLAB)

	// Act
	_, _, err := handler(ctx, req, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
}

func TestHandler_LoggerFactoryError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &basetoolmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockConfigFactory := &custommocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockUsecase := &custommocks.MockUsecase{}
	defer mockUsecase.AssertExpectations(t)

	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	ctx := t.Context()
	expectedSession := &mcp.ServerSession{}
	mockValidatedTool := &definitionmocks.MockValidatedTool{}
	defer mockValidatedTool.AssertExpectations(t)

	args := map[string]any{"n": float64(5)}
	expectedError := messages.AnError
	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	mockValidatedTool.EXPECT().
		Definition().
		Return(definition.Tool{}).
		Once()
	mockValidatedTool.EXPECT().
		Signature().
		Return(definition.Signature{}).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(nil, expectedError).
		Once()

	handler := custom.Handler(mockValidatedTool, mockLoggerFactory, mockConfigFactory, mockUsecase, mockGlobalMATLAB)

	// Act
	_, _, err := handler(ctx, req, args)

	// Assert
	require.ErrorIs(t, err, expectedError)
}
