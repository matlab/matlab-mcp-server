// Copyright 2026 The MathWorks, Inc.

package sessionmanager

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/time/retry"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

var ErrFailedToAttachToMATLABSession = errors.New("failed to attach to MATLAB session")

const defaultDiscoveryRetryInterval = 1 * time.Second

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type MATLABManager interface {
	StartMATLABSession(ctx context.Context, sessionLogger entities.Logger, startRequest entities.SessionDetails) (entities.SessionID, error)
	StopMATLABSession(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) error
	GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error)
}

type MATLABRootSelector interface {
	SelectMATLABRoot(ctx context.Context, logger entities.Logger) (string, error)
}

type MATLABStartingDirSelector interface {
	SelectMATLABStartingDir(logger entities.Logger) (string, error)
}

type SessionManager struct {
	matlabManager             MATLABManager
	configFactory             ConfigFactory
	matlabRootSelector        MATLABRootSelector
	matlabStartingDirSelector MATLABStartingDirSelector

	discoveryRetryInterval time.Duration

	initOnce          sync.Once
	initErr           error
	matlabRoot        string
	matlabStartingDir string
}

func New(
	matlabManager MATLABManager,
	configFactory ConfigFactory,
	matlabRootSelector MATLABRootSelector,
	matlabStartingDirSelector MATLABStartingDirSelector,
) *SessionManager {
	return &SessionManager{
		matlabManager:             matlabManager,
		configFactory:             configFactory,
		matlabRootSelector:        matlabRootSelector,
		matlabStartingDirSelector: matlabStartingDirSelector,

		discoveryRetryInterval: defaultDiscoveryRetryInterval,
	}
}

func (s *SessionManager) StartSession(ctx context.Context, logger entities.Logger) (entities.SessionID, error) {
	cfg, messagesErr := s.configFactory.Config()
	if messagesErr != nil {
		return 0, messagesErr
	}

	var sessionID entities.SessionID
	var err error

	switch cfg.MATLABSessionMode() {
	case entities.MATLABSessionModeExisting:
		sessionID, err = s.getSessionFromAttachingToExistingMATLAB(ctx, logger, cfg)
		if err != nil {
			logger.WithError(err).Debug("failed to attach to MATLAB session")
			err = ErrFailedToAttachToMATLABSession
		}
	case entities.MATLABSessionModeAuto:
		sessionID, err = s.getSessionInAutoMode(ctx, logger, cfg)
	default:
		sessionID, err = s.getSessionFromLocalMATLABInstallation(ctx, logger, cfg)
	}

	if err != nil {
		return 0, err
	}

	return sessionID, nil
}

func (s *SessionManager) ShouldRestart() (bool, messages.Error) {
	cfg, err := s.configFactory.Config()
	if err != nil {
		return false, err
	}

	shouldNotRestart := cfg.MATLABSessionMode() == entities.MATLABSessionModeExisting &&
		cfg.MATLABSessionConnectionDetails() != ""

	return !shouldNotRestart, nil
}

func (s *SessionManager) StopMATLABSession(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) error {
	return s.matlabManager.StopMATLABSession(ctx, sessionLogger, sessionID)
}

func (s *SessionManager) GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error) {
	return s.matlabManager.GetMATLABSessionClient(ctx, sessionLogger, sessionID)
}

func (s *SessionManager) initializeStartupConfig(ctx context.Context, logger entities.Logger) error {
	matlabRoot, err := s.matlabRootSelector.SelectMATLABRoot(ctx, logger)
	if err != nil {
		return err
	}

	s.matlabRoot = matlabRoot

	matlabStartingDirectory, err := s.matlabStartingDirSelector.SelectMATLABStartingDir(logger)
	if err != nil {
		logger.WithError(err).Warn("failed to determine MATLAB starting directory, proceeding without one")
		return nil
	}

	s.matlabStartingDir = matlabStartingDirectory
	return nil
}

func (s *SessionManager) getSessionFromLocalMATLABInstallation(ctx context.Context, logger entities.Logger, cfg config.Config) (entities.SessionID, error) {
	s.initOnce.Do(func() {
		s.initErr = s.initializeStartupConfig(ctx, logger)
	})
	if s.initErr != nil {
		return 0, s.initErr
	}

	startRequest := entities.LocalSessionDetails{
		MATLABRoot:             s.matlabRoot,
		IsStartingDirectorySet: s.matlabStartingDir != "",
		StartingDirectory:      s.matlabStartingDir,
		ShowMATLABDesktop:      cfg.ShouldShowMATLABDesktop(),
	}

	return s.matlabManager.StartMATLABSession(ctx, logger, startRequest)
}

func (s *SessionManager) getSessionInAutoMode(ctx context.Context, logger entities.Logger, cfg config.Config) (entities.SessionID, error) {
	sessionID, err := s.getSessionFromAttachingToExistingMATLAB(ctx, logger, cfg)
	if err != nil {
		logger.
			WithError(err).
			Debug("auto mode: no existing MATLAB session found, falling back to launching MATLAB")
		return s.getSessionFromLocalMATLABInstallation(ctx, logger, cfg)
	}

	return sessionID, nil
}

func (s *SessionManager) getSessionFromAttachingToExistingMATLAB(ctx context.Context, logger entities.Logger, cfg config.Config) (entities.SessionID, error) {
	startRequest := entities.AttachToExistingSession{}

	discoveryTimeout := cfg.MATLABSessionDiscoveryTimeout()

	// A zero timeout means "try once, no polling". Use half the retry interval
	// so the context stays alive long enough for a single attempt but expires
	// before the retry strategy fires a second one.
	if discoveryTimeout == 0 {
		discoveryTimeout = s.discoveryRetryInterval / 2
	}

	attachCtx, cancel := context.WithTimeout(ctx, discoveryTimeout)
	defer cancel()

	var lastErr error

	sessionID, err := retry.Retry(attachCtx, func() (entities.SessionID, bool, error) {
		sessionID, err := s.matlabManager.StartMATLABSession(attachCtx, logger, startRequest)
		if err != nil {
			lastErr = err
			return 0, false, nil
		}

		return sessionID, true, nil
	}, retry.NewLinearRetryStrategy(s.discoveryRetryInterval))

	if err != nil {
		if lastErr != nil {
			// Return the last error rather than timeout error
			return 0, lastErr
		}
		return 0, err
	}

	return sessionID, nil
}
