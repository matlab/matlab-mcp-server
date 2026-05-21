// Copyright 2026 The MathWorks, Inc.
//go:build !windows

package resourcelimit_test

import (
	"errors"
	"testing"

	"golang.org/x/sys/unix"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/resourcelimit"
	unixfacade "github.com/matlab/matlab-mcp-core-server/internal/facades/unix"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	resourcelimitmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/resourcelimit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Arrange
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	// Act
	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Assert
	require.NotNil(t, manager)
}

func TestCapOpenFilesLimit_GetGlobalLoggerFails(t *testing.T) {
	// Arrange
	const rlimDesired uint64 = 1024

	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(nil, messages.AnError).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)

	// Assert
	require.Error(t, err)
	assert.Nil(t, reset)
	require.ErrorIs(t, err, messages.AnError)
}

func TestCapOpenFilesLimit_CurrentBelowLimit_NoOp(t *testing.T) {
	// Arrange
	const rlimCur uint64 = 512
	const rlimMax uint64 = 1024
	const rlimDesired uint64 = 1024

	logger := testutils.NewInspectableLogger()
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockSyscall.EXPECT().
		Getrlimit(unix.RLIMIT_NOFILE, mock.AnythingOfType("*unix.Rlimit")).
		RunAndReturn(func(_ int, rlim *unixfacade.Rlimit) error {
			rlim.Cur = rlimCur
			rlim.Max = rlimMax
			return nil
		}).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, reset)
	require.NoError(t, reset())
}

func TestCapOpenFilesLimit_CurrentEqualToLimit_NoOp(t *testing.T) {
	// Arrange
	const rlimCur uint64 = 1024
	const rlimMax uint64 = 4096
	const rlimDesired uint64 = 1024

	logger := testutils.NewInspectableLogger()
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockSyscall.EXPECT().
		Getrlimit(unix.RLIMIT_NOFILE, mock.AnythingOfType("*unix.Rlimit")).
		RunAndReturn(func(_ int, rlim *unixfacade.Rlimit) error {
			rlim.Cur = rlimCur
			rlim.Max = rlimMax
			return nil
		}).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, reset)
	require.NoError(t, reset())
}

func TestCapOpenFilesLimit_CurrentAboveLimit_LowersLimit(t *testing.T) {
	// Arrange
	const rlimCur uint64 = 65535
	const rlimMax uint64 = 65535
	const rlimDesired uint64 = 1024

	logger := testutils.NewInspectableLogger()
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockSyscall.EXPECT().
		Getrlimit(unix.RLIMIT_NOFILE, mock.AnythingOfType("*unix.Rlimit")).
		RunAndReturn(func(_ int, rlim *unixfacade.Rlimit) error {
			rlim.Cur = rlimCur
			rlim.Max = rlimMax
			return nil
		}).
		Once()

	mockSyscall.EXPECT().
		Setrlimit(unix.RLIMIT_NOFILE, mock.MatchedBy(func(rlim *unixfacade.Rlimit) bool {
			return rlim.Cur == rlimDesired && rlim.Max == rlimMax
		})).
		Return(nil).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, reset)
}

func TestCapOpenFilesLimit_ResetRestoresOriginal(t *testing.T) {
	// Arrange
	const rlimCur uint64 = 65535
	const rlimMax uint64 = 65535
	const rlimDesired uint64 = 1024

	logger := testutils.NewInspectableLogger()
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockSyscall.EXPECT().
		Getrlimit(unix.RLIMIT_NOFILE, mock.AnythingOfType("*unix.Rlimit")).
		RunAndReturn(func(_ int, rlim *unixfacade.Rlimit) error {
			rlim.Cur = rlimCur
			rlim.Max = rlimMax
			return nil
		}).
		Once()

	// First Setrlimit call: lowering
	mockSyscall.EXPECT().
		Setrlimit(unix.RLIMIT_NOFILE, mock.MatchedBy(func(rlim *unixfacade.Rlimit) bool {
			return rlim.Cur == rlimDesired && rlim.Max == rlimMax
		})).
		Return(nil).
		Once()

	// Second Setrlimit call: restoring
	mockSyscall.EXPECT().
		Setrlimit(unix.RLIMIT_NOFILE, mock.MatchedBy(func(rlim *unixfacade.Rlimit) bool {
			return rlim.Cur == rlimCur && rlim.Max == rlimMax
		})).
		Return(nil).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)
	require.NoError(t, err)

	err = reset()

	// Assert
	require.NoError(t, err)
}

func TestCapOpenFilesLimit_GetrlimitFails(t *testing.T) {
	// Arrange
	const rlimDesired uint64 = 1024

	logger := testutils.NewInspectableLogger()
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	expectedErr := errors.New("getrlimit error")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockSyscall.EXPECT().
		Getrlimit(unix.RLIMIT_NOFILE, mock.AnythingOfType("*unix.Rlimit")).
		Return(expectedErr).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)

	// Assert
	require.Error(t, err)
	assert.Nil(t, reset)
	require.ErrorIs(t, err, expectedErr)
	assert.NotEmpty(t, logger.ErrorLogs())
}

func TestCapOpenFilesLimit_SetrlimitFails(t *testing.T) {
	// Arrange
	const rlimCur uint64 = 65535
	const rlimMax uint64 = 65535
	const rlimDesired uint64 = 1024

	logger := testutils.NewInspectableLogger()
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	expectedErr := errors.New("setrlimit error")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockSyscall.EXPECT().
		Getrlimit(unix.RLIMIT_NOFILE, mock.AnythingOfType("*unix.Rlimit")).
		RunAndReturn(func(_ int, rlim *unixfacade.Rlimit) error {
			rlim.Cur = rlimCur
			rlim.Max = rlimMax
			return nil
		}).
		Once()

	mockSyscall.EXPECT().
		Setrlimit(unix.RLIMIT_NOFILE, mock.MatchedBy(func(rlim *unixfacade.Rlimit) bool {
			return rlim.Cur == rlimDesired && rlim.Max == rlimMax
		})).
		Return(expectedErr).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)

	// Assert
	require.Error(t, err)
	assert.Nil(t, reset)
	require.ErrorIs(t, err, expectedErr)
	assert.NotEmpty(t, logger.ErrorLogs())
}

func TestCapOpenFilesLimit_ResetFails(t *testing.T) {
	// Arrange
	const rlimCur uint64 = 65535
	const rlimMax uint64 = 65535
	const rlimDesired uint64 = 1024

	logger := testutils.NewInspectableLogger()
	mockLoggerFactory := resourcelimitmocks.NewMockLoggerFactory(t)
	mockSyscall := resourcelimitmocks.NewMockSyscallLayer(t)

	expectedErr := errors.New("restore error")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(logger, nil).
		Once()

	mockSyscall.EXPECT().
		Getrlimit(unix.RLIMIT_NOFILE, mock.AnythingOfType("*unix.Rlimit")).
		RunAndReturn(func(_ int, rlim *unixfacade.Rlimit) error {
			rlim.Cur = rlimCur
			rlim.Max = rlimMax
			return nil
		}).
		Once()

	// Lowering succeeds
	mockSyscall.EXPECT().
		Setrlimit(unix.RLIMIT_NOFILE, mock.MatchedBy(func(rlim *unixfacade.Rlimit) bool {
			return rlim.Cur == rlimDesired && rlim.Max == rlimMax
		})).
		Return(nil).
		Once()

	// Restoring fails
	mockSyscall.EXPECT().
		Setrlimit(unix.RLIMIT_NOFILE, mock.MatchedBy(func(rlim *unixfacade.Rlimit) bool {
			return rlim.Cur == rlimCur && rlim.Max == rlimMax
		})).
		Return(expectedErr).
		Once()

	manager := resourcelimit.New(mockLoggerFactory, mockSyscall)

	// Act
	reset, err := manager.CapOpenFilesLimit(rlimDesired)
	require.NoError(t, err)

	err = reset()

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}
