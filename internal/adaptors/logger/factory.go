// Copyright 2025-2026 The MathWorks, Inc.

package logger

import (
	"io"
	"log/slog"
	"path/filepath"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultGlobalLogLevel slog.Level = slog.LevelDebug

const (
	logFileName         = "server"
	watchdogLogFileName = "watchdog"
	logFileExt          = ".log"
)

type ConfigFactory interface {
	Config() (config.Config, messages.Error)
}

type DirectoryFactory interface {
	Directory() (directory.Directory, messages.Error)
}

type FilenameFactory interface {
	FilenameWithSuffix(fileName string, ext string, suffix string) string
}

type OSLayer interface {
	Stderr() io.Writer
	Create(name string) (osfacade.File, error)
}

type Factory struct {
	configFactory    ConfigFactory
	directoryFactory DirectoryFactory
	filenameFactory  FilenameFactory
	osLayer          OSLayer

	initOnce              sync.Once
	initError             messages.Error
	parsedLogLevel        slog.Level
	duplicateLogsToStderr bool
	logFile               entities.Writer

	globalLoggerOnce sync.Once
	globalLogger     *slogLogger
}

func NewFactory(
	configFactory ConfigFactory,
	directoryFactory DirectoryFactory,
	filenameFactory FilenameFactory,
	osLayer OSLayer,
) *Factory {
	return &Factory{
		configFactory:    configFactory,
		directoryFactory: directoryFactory,
		filenameFactory:  filenameFactory,
		osLayer:          osLayer,

		parsedLogLevel: defaultGlobalLogLevel,
	}
}

func (f *Factory) NewMCPSessionLogger(session *mcp.ServerSession) (entities.Logger, messages.Error) {
	// In MCP Server development, special care should be given to logging:
	//
	// https://modelcontextprotocol.io/quickstart/server#logging-in-mcp-servers
	//
	// In essence, you can't log to standard `stdout`, and while you may log to `stderr`, you should log to the client:
	//
	// https://modelcontextprotocol.io/specification/2025-06-18/server/utilities/logging
	if err := f.initialize(); err != nil {
		return nil, err
	}

	sessionHandler := mcp.NewLoggingHandler(session, &mcp.LoggingHandlerOptions{})

	handler := slog.NewJSONHandler(f.logFile, &slog.HandlerOptions{
		Level: f.parsedLogLevel,
	})

	return &slogLogger{
		logger: slog.New(NewMultiHandler(sessionHandler, handler)),
	}, nil
}

func (f *Factory) GetGlobalLogger() (entities.Logger, messages.Error) {
	// There are cases where we want to log, but we don't have an MCP session yet.
	// In those cases, we log to the log file and, optionally, to stderr if requested.
	// Logging to stderr is allowed by the MCP spec:
	//
	// https://modelcontextprotocol.io/docs/develop/build-server#best-practices
	if err := f.initialize(); err != nil {
		return nil, err
	}

	f.globalLoggerOnce.Do(func() {
		logWriter := f.logFile
		if f.duplicateLogsToStderr {
			logWriter = io.MultiWriter(f.osLayer.Stderr(), f.logFile)
		}

		handler := slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
			Level: f.parsedLogLevel,
		})
		f.globalLogger = &slogLogger{
			logger: slog.New(handler),
		}
	})

	return f.globalLogger, nil
}

func (f *Factory) initialize() messages.Error {
	f.initOnce.Do(func() {
		config, messagesErr := f.configFactory.Config()
		if messagesErr != nil {
			f.initError = messagesErr
			return
		}

		f.duplicateLogsToStderr = config.DuplicateLogsToStderr()

		var parsedLogLevel slog.Level
		logLevel := config.LogLevel()
		switch logLevel {
		case entities.LogLevelDebug:
			parsedLogLevel = slog.LevelDebug
		case entities.LogLevelInfo:
			parsedLogLevel = slog.LevelInfo
		case entities.LogLevelWarn:
			parsedLogLevel = slog.LevelWarn
		case entities.LogLevelError:
			parsedLogLevel = slog.LevelError
		default:
			f.initError = messages.New_StartupErrors_InvalidLogLevel_Error(string(logLevel))
			return
		}
		f.parsedLogLevel = parsedLogLevel

		dir, messagesErr := f.directoryFactory.Directory()
		if messagesErr != nil {
			f.initError = messagesErr
			return
		}

		baseDir := dir.BaseDir()
		id := dir.ID()

		logFileBase := filepath.Join(baseDir, logFileName)
		if config.WatchdogMode() {
			logFileBase = filepath.Join(baseDir, watchdogLogFileName)
		}

		logFilePath := f.filenameFactory.FilenameWithSuffix(logFileBase, logFileExt, id)

		logFile, err := f.osLayer.Create(logFilePath)
		if err != nil {
			f.initError = messages.New_StartupErrors_FailedToCreateLogFile_Error(logFilePath)
			return
		}

		f.logFile = logFile
	})

	return f.initError
}
