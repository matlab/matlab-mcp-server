// Copyright 2026 The MathWorks, Inc.

package defaultparameters

import (
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

const envVarNamePrefix = "MW_MCP_SERVER_"

func HelpMode() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "HelpMode",
		/* flagName */ "help",
		/* hiddenFlag */ false,
		/* envVarName */ "",
		/* descriptionKey */ messages.CLIMessages_HelpDescription,
		/* defaultValue */ false,
		/* recordToLog */ false,
		/* piiSafe */ true,
	)
}

func VersionMode() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "VersionMode",
		/* flagName */ "version",
		/* hiddenFlag */ false,
		/* envVarName */ "",
		/* descriptionKey */ messages.CLIMessages_VersionDescription,
		/* defaultValue */ false,
		/* recordToLog */ false,
		/* piiSafe */ true,
	)
}

func SetupMATLABMode() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "SetupMATLABMode",
		/* flagName */ "setup-matlab",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"SERVER_SETUP_MATLAB",
		/* descriptionKey */ messages.CLIMessages_SetupMATLABDescription,
		/* defaultValue */ false,
		/* recordToLog */ false,
		/* piiSafe */ true,
	)
}

func PreferredLocalMATLABRoot() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "PreferredLocalMATLABRoot",
		/* flagName */ "matlab-root",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"MATLAB_ROOT",
		/* descriptionKey */ messages.CLIMessages_PreferredLocalMATLABRootDescription,
		/* defaultValue */ "",
		/* recordToLog */ true,
		/* piiSafe */ false,
	)
}

func PreferredMATLABStartingDirectory() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "PreferredMATLABStartingDirectory",
		/* flagName */ "initial-working-folder",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"INITIAL_WORKING_FOLDER",
		/* descriptionKey */ messages.CLIMessages_PreferredMATLABStartingDirectoryDescription,
		/* defaultValue */ "",
		/* recordToLog */ true,
		/* piiSafe */ false,
	)
}

func BaseDir() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "BaseDir",
		/* flagName */ "log-folder",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"LOG_FOLDER",
		/* descriptionKey */ messages.CLIMessages_BaseDirDescription,
		/* defaultValue */ "",
		/* recordToLog */ false,
		/* piiSafe */ false,
	)
}

func LogLevel() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "LogLevel",
		/* flagName */ "log-level",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"LOG_LEVEL",
		/* descriptionKey */ messages.CLIMessages_LogLevelDescription,
		/* defaultValue */ "info",
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func DuplicateLogsToStderr() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "DuplicateLogsToStderr",
		/* flagName */ "duplicate-logs-to-stderr",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"DUPLICATE_LOGS_TO_STDERR",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ false,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func InitializeMATLABOnStartup() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "InitializeMATLABOnStartup",
		/* flagName */ "initialize-matlab-on-startup",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"INITIALIZE_MATLAB_ON_STARTUP",
		/* descriptionKey */ messages.CLIMessages_InitializeMATLABOnStartupDescription,
		/* defaultValue */ false,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func MATLABDisplayMode() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "MATLABDisplayMode",
		/* flagName */ "matlab-display-mode",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"MATLAB_DISPLAY_MODE",
		/* descriptionKey */ messages.CLIMessages_DisplayModeDescription,
		/* defaultValue */ "desktop",
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func UseSingleMATLABSession() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "UseSingleMATLABSession",
		/* flagName */ "use-single-matlab-session",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"USE_SINGLE_MATLAB_SESSION",
		/* descriptionKey */ messages.CLIMessages_UseSingleMATLABSessionDescription,
		/* defaultValue */ true,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func WatchdogMode() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "WatchdogMode",
		/* flagName */ "watchdog",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"WATCHDOG_MODE",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ false,
		/* recordToLog */ false,
		/* piiSafe */ true,
	)
}

func ServerInstanceID() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "ServerInstanceID",
		/* flagName */ "server-instance-id",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"SERVER_INSTANCE_ID",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ "",
		/* recordToLog */ false,
		/* piiSafe */ false,
	)
}

func MATLABSessionMode() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "MATLABSessionMode",
		/* flagName */ "matlab-session-mode",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"MATLAB_SESSION_MODE",
		/* descriptionKey */ messages.CLIMessages_MATLABSessionModeDescription,
		/* defaultValue */ string(entities.MATLABSessionModeNew),
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func MATLABSessionConnectionDetails() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "MATLABSessionConnectionDetails",
		/* flagName */ "matlab-session-connection-details",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"MATLAB_SESSION_CONNECTION_DETAILS",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ "",
		/* recordToLog */ false,
		/* piiSafe */ false,
	)
}

func MATLABSessionConnectionTimeout() *parameter.Parameter[time.Duration] {
	return parameter.NewParameter(
		/* id */ "MATLABSessionConnectionTimeout",
		/* flagName */ "matlab-session-connection-timeout",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"MATLAB_SESSION_CONNECTION_TIMEOUT",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ 5*time.Second,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func MATLABSessionDiscoveryTimeout() *parameter.Parameter[time.Duration] {
	return parameter.NewParameter(
		/* id */ "MATLABSessionDiscoveryTimeout",
		/* flagName */ "matlab-session-discovery-timeout",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"MATLAB_SESSION_DISCOVERY_TIMEOUT",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ 30*time.Second,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func DisableTelemetry() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "DisableTelemetry",
		/* flagName */ "disable-telemetry",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"DISABLE_TELEMETRY",
		/* descriptionKey */ messages.CLIMessages_DisableTelemetryDescription,
		/* defaultValue */ false,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func TelemetryCollectorEndpoint() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "TelemetryCollectorEndpoint",
		/* flagName */ "telemetry-collector-endpoint",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"TELEMETRY_COLLECTOR_ENDPOINT",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ "",
		/* recordToLog */ true,
		/* piiSafe */ false,
	)
}

func TelemetryCollectionInterval() *parameter.Parameter[time.Duration] {
	return parameter.NewParameter(
		/* id */ "TelemetryCollectionInterval",
		/* flagName */ "telemetry-collection-interval",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"TELEMETRY_COLLECTION_INTERVAL",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ time.Minute,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func TelemetryCollectorEndpointInsecure() *parameter.Parameter[bool] {
	return parameter.NewParameter(
		/* id */ "TelemetryCollectorEndpointInsecure",
		/* flagName */ "telemetry-collector-endpoint-insecure",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"TELEMETRY_COLLECTOR_ENDPOINT_INSECURE",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ false,
		/* recordToLog */ true,
		/* piiSafe */ true,
	)
}

func EmbeddedConnectorDetailsTimeout() *parameter.Parameter[time.Duration] {
	return parameter.NewParameter(
		/* id */ "EmbeddedConnectorDetailsTimeout",
		/* flagName */ "",
		/* hiddenFlag */ true,
		/* envVarName */ envVarNamePrefix+"EMBEDDED_CONNECTOR_DETAILS_TIMEOUT",
		/* descriptionKey */ messages.CLIMessages_InternalUseDescription,
		/* defaultValue */ 10*time.Minute,
		/* recordToLog */ false,
		/* piiSafe */ true,
	)
}

func ExtensionFile() *parameter.Parameter[string] {
	return parameter.NewParameter(
		/* id */ "ExtensionFile",
		/* flagName */ "extension-file",
		/* hiddenFlag */ false,
		/* envVarName */ envVarNamePrefix+"EXTENSION_FILE",
		/* descriptionKey */ messages.CLIMessages_ExtensionFileDescription,
		/* defaultValue */ "",
		/* recordToLog */ true,
		/* piiSafe */ false,
	)
}
