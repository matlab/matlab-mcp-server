// Copyright 2025-2026 The MathWorks, Inc.

//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/definition"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/lifecyclesignaler"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector/modes/setupmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/orchestrator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/defaultparameters/selector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/parameter/parser"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/buildinfo"
	files "github.com/matlab/matlab-mcp-core-server/internal/adaptors/filesystem/files"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager/matlabstartingdirselector"
	httpclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/client"
	httpserver "github.com/matlab/matlab-mcp-core-server/internal/adaptors/http/server"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/logger"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlab/codeanalyzer"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlab/matlabrootselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/addonmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/addonmanager/installationsteps"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession"
	localmatlabsessiondirectory "github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directory/matlabfiles"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processdetails"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processlauncher"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabroot"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabversion"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/sessionselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/sessionselector/sessiondiscovery"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/sessionselector/sessiondiscovery/appdatadir"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/baseresource"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/codingguidelines"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/resources/plaintextlivecodegeneration"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/configurator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/rootpathresolver"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/rootstore"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/sdk"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	evalmatlabcodemultisessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	listavailablematlabstool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	startmatlabsessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	stopmatlabsessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	checkmatlabcodesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom"
	customloader "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/loader"
	customvalidator "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/custom/loader/validator"
	detectmatlabtoolboxessinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabcodesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	runmatlabfilesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	runmatlabtestfilesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/messagecatalog"
	osadaptor "github.com/matlab/matlab-mcp-core-server/internal/adaptors/os"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/telemetry"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/telemetry/otel/instruments"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/telemetry/otel/meter/exporter"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/telemetry/otel/meter/provider"
	watchdogclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog/process"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/filefacade"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/iofacade"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/registryfacade"
	unixfacade "github.com/matlab/matlab-mcp-core-server/internal/facades/unix"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/resourcelimit"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/detectmatlabtoolboxes"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalcustomtool/functioncall"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/utils/pathvalidator"
	watchdogprocess "github.com/matlab/matlab-mcp-core-server/internal/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/processhandler"
	transportclient "github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/client"
	transportserver "github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/server/handler"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport/socket"
)

type Application struct {
	ModeSelector        *modeselector.ModeSelector
	MessageCatalog      *messagecatalog.MessageCatalog
	MATLABClientFactory *matlabsessionclient.Factory
	HTTPClientFactory   *httpclient.Factory
	HTTPServerFactory   *httpserver.Factory
	LoggerFactory       *logger.Factory
	// Exposed for integration testing
	LocalMATLABSessionStarter *localmatlabsession.Starter
}

type ApplicationDefinition interface {
	Name() string
	Title() string
	Instructions() string
	Features() definition.Features
	Parameters() []entities.Parameter
	Dependencies(resources definition.DependenciesProviderResources) (any, error)
	Tools(resources definition.ToolsProviderResources) []tools.Tool
}

func Initialize(serverDefinition ApplicationDefinition) *Application {
	wire.Build(
		// Application
		wire.Struct(new(Application), "*"),

		// Mode Selector
		modeselector.New,
		wire.Bind(new(modeselector.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(modeselector.Parser), new(*parser.Parser)),
		wire.Bind(new(modeselector.TelemetryFactory), new(*telemetry.Factory)),
		wire.Bind(new(modeselector.WatchdogProcess), new(*watchdogprocess.Watchdog)),
		wire.Bind(new(modeselector.Orchestrator), new(*orchestrator.Orchestrator)),
		wire.Bind(new(modeselector.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(modeselector.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),
		wire.Bind(new(modeselector.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(modeselector.SetupMATLAB), new(*setupmatlab.Mode)),

		// Setup MATLAB
		setupmatlab.New,
		wire.Bind(new(setupmatlab.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(setupmatlab.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(setupmatlab.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(setupmatlab.MessageCatalog), new(*messagecatalog.MessageCatalog)),
		wire.Bind(new(setupmatlab.WatchdogClient), new(*watchdogclient.Watchdog)),
		wire.Bind(new(setupmatlab.GlobalMATLAB), new(*globalmatlab.GlobalMATLAB)),
		wire.Bind(new(setupmatlab.AddonManager), new(*addonmanager.AddonManager)),

		// Add-On Manager
		addonmanager.New,
		wire.Bind(new(addonmanager.InstallationSteps), new(*installationsteps.InstallationSteps)),

		// Add-On Manager Installation Steps
		installationsteps.New,

		// Telemetry
		telemetry.NewFactory,
		wire.Bind(new(telemetry.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(telemetry.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(telemetry.ExporterFactory), new(*exporter.Factory)),
		wire.Bind(new(telemetry.MeterProviderFactory), new(*provider.Factory)),
		wire.Bind(new(telemetry.InstrumentFactory), new(*instruments.Factory)),
		wire.Bind(new(telemetry.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(telemetry.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(telemetry.OSVersionProvider), new(osadaptor.OS)),
		wire.Bind(new(telemetry.Definition), new(ApplicationDefinition)),

		// Telemetry Instruments
		instruments.NewFactory,

		// Telemetry Exporter
		exporter.NewFactory,
		wire.Bind(new(exporter.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(exporter.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(exporter.OSLayer), new(*osfacade.OsFacade)),

		// Telemetry Provider
		provider.NewFactory,
		wire.Bind(new(provider.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(provider.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(provider.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),

		// Watchdog Process
		watchdogprocess.New,
		wire.Bind(new(watchdogprocess.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(watchdogprocess.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(watchdogprocess.ProcessHandler), new(*processhandler.ProcessHandler)),
		wire.Bind(new(watchdogprocess.OSSignaler), new(*osadaptor.ProcessManager)),
		wire.Bind(new(watchdogprocess.ServerHandlerFactory), new(*handler.Factory)),
		wire.Bind(new(watchdogprocess.ServerFactory), new(*transportserver.Factory)),
		wire.Bind(new(watchdogprocess.SocketFactory), new(*socket.Factory)),

		// Watchdog Transport Server Factory
		transportserver.NewFactory,
		wire.Bind(new(transportserver.HTTPServerFactory), new(*httpserver.Factory)),
		wire.Bind(new(transportserver.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(transportserver.HandlerFactory), new(*handler.Factory)),

		// HTTP Server Factory
		httpserver.NewFactory,
		wire.Bind(new(httpserver.OSLayer), new(*osfacade.OsFacade)),

		// Orchestrator
		orchestrator.New,
		wire.Bind(new(orchestrator.MessageCatalog), new(*messagecatalog.MessageCatalog)),
		wire.Bind(new(orchestrator.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),
		wire.Bind(new(orchestrator.ApplicationDefinition), new(ApplicationDefinition)),
		wire.Bind(new(orchestrator.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(orchestrator.Server), new(*server.Server)),
		wire.Bind(new(orchestrator.WatchdogClient), new(*watchdogclient.Watchdog)),
		wire.Bind(new(orchestrator.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(orchestrator.OSSignaler), new(*osadaptor.ProcessManager)),
		wire.Bind(new(orchestrator.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(orchestrator.ResourceLimitManager), new(*resourcelimit.Manager)),

		// MCP Server
		server.New,
		wire.Bind(new(server.MCPSDKServerFactory), new(*sdk.Factory)),
		wire.Bind(new(server.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(server.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),
		wire.Bind(new(server.MCPServerConfigurator), new(*configurator.Configurator)),

		// RootStore
		rootstore.New,

		// Root Path Resolver
		rootpathresolver.New,
		wire.Bind(new(rootpathresolver.OSLayer), new(*osfacade.OsFacade)),

		// MCP Server (SDK)
		sdk.NewFactory,
		wire.Bind(new(sdk.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(sdk.Definition), new(ApplicationDefinition)),
		wire.Bind(new(sdk.RootStore), new(*rootstore.RootStore)),
		wire.Bind(new(sdk.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(sdk.GlobalMATLAB), new(*globalmatlab.GlobalMATLAB)),

		// MCP Server Configurator
		configurator.New,
		wire.Bind(new(configurator.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(configurator.ApplicationDefinition), new(ApplicationDefinition)),
		wire.Bind(new(configurator.CustomToolFactory), new(*custom.Factory)),

		// Tools
		wire.Bind(new(basetool.LoggerFactory), new(*logger.Factory)),

		listavailablematlabstool.New,
		wire.Bind(new(listavailablematlabstool.Usecase), new(*listavailablematlabs.Usecase)),

		listavailablematlabs.New,

		startmatlabsessiontool.New,
		wire.Bind(new(startmatlabsessiontool.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(startmatlabsessiontool.Usecase), new(*startmatlabsession.Usecase)),

		startmatlabsession.New,

		stopmatlabsessiontool.New,
		wire.Bind(new(stopmatlabsessiontool.Usecase), new(*stopmatlabsession.Usecase)),

		stopmatlabsession.New,

		evalmatlabcodemultisessiontool.New,
		wire.Bind(new(evalmatlabcodemultisessiontool.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(evalmatlabcodemultisessiontool.Usecase), new(*evalmatlabcode.Usecase)),

		evalmatlabcodesinglesessiontool.New,
		wire.Bind(new(evalmatlabcodesinglesessiontool.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(evalmatlabcodesinglesessiontool.Usecase), new(*evalmatlabcode.Usecase)),

		evalmatlabcode.New,
		wire.Bind(new(evalmatlabcode.PathValidator), new(*pathvalidator.PathValidator)),

		checkmatlabcodesinglesessiontool.New,
		wire.Bind(new(checkmatlabcodesinglesessiontool.Usecase), new(*checkmatlabcode.Usecase)),

		checkmatlabcode.New,
		wire.Bind(new(checkmatlabcode.PathValidator), new(*pathvalidator.PathValidator)),
		wire.Bind(new(checkmatlabcode.CodeAnalyzer), new(*codeanalyzer.Analyzer)),

		codeanalyzer.New,

		detectmatlabtoolboxessinglesessiontool.New,
		wire.Bind(new(detectmatlabtoolboxessinglesessiontool.Usecase), new(*detectmatlabtoolboxes.Usecase)),

		detectmatlabtoolboxes.New,

		runmatlabfilesinglesessiontool.New,
		wire.Bind(new(runmatlabfilesinglesessiontool.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(runmatlabfilesinglesessiontool.Usecase), new(*runmatlabfile.Usecase)),

		runmatlabfile.New,
		wire.Bind(new(runmatlabfile.PathValidator), new(*pathvalidator.PathValidator)),

		runmatlabtestfilesinglesessiontool.New,
		wire.Bind(new(runmatlabtestfilesinglesessiontool.Usecase), new(*runmatlabtestfile.Usecase)),

		runmatlabtestfile.New,
		wire.Bind(new(runmatlabtestfile.PathValidator), new(*pathvalidator.PathValidator)),

		// Custom Tool Factory
		custom.NewFactory,
		wire.Bind(new(custom.Loader), new(*customloader.Loader)),
		wire.Bind(new(custom.ConfigFactory), new(*config.Factory)),

		// Custom Tool Loader
		customloader.NewLoader,
		wire.Bind(new(customloader.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(customloader.LoggerFactory), new(*logger.Factory)),
		customvalidator.NewValidator,
		wire.Bind(new(customloader.ToolValidator), new(*customvalidator.Validator)),

		// EvalCustomTool Use Case
		evalcustomtool.New,
		wire.Bind(new(evalcustomtool.FunctionCallAssembler), new(*functioncall.Assembler)),
		functioncall.NewAssembler,
		wire.Bind(new(custom.Usecase), new(*evalcustomtool.Usecase)),

		// Resources
		wire.Bind(new(baseresource.LoggerFactory), new(*logger.Factory)),

		codingguidelines.New,
		plaintextlivecodegeneration.New,

		// Watchdog Client
		watchdogclient.New,
		wire.Bind(new(watchdogclient.WatchdogProcess), new(*process.Factory)),
		wire.Bind(new(watchdogclient.ClientFactory), new(*transportclient.Factory)),
		wire.Bind(new(watchdogclient.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(watchdogclient.SocketFactory), new(*socket.Factory)),

		// Watchdog Process Handler for Watchdog Client
		process.New,
		wire.Bind(new(process.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(process.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(process.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(process.ConfigFactory), new(*config.Factory)),

		// Watchdog Transport Client Factory
		transportclient.NewFactory,
		wire.Bind(new(transportclient.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(transportclient.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(transportclient.HTTPClientFactory), new(*httpclient.Factory)),

		// Global MATLAB
		globalmatlab.New,
		wire.Bind(new(globalmatlab.MATLABManagerAdaptor), new(*sessionmanager.SessionManager)),

		// Session Manager
		sessionmanager.New,
		wire.Bind(new(sessionmanager.MATLABManager), new(*matlabmanager.MATLABManager)),
		wire.Bind(new(sessionmanager.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(sessionmanager.MATLABRootSelector), new(*matlabrootselector.MATLABRootSelector)),
		wire.Bind(new(sessionmanager.MATLABStartingDirSelector), new(*matlabstartingdirselector.MATLABStartingDirSelector)),

		// MATLAB Root Selector
		matlabrootselector.New,
		wire.Bind(new(matlabrootselector.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(matlabrootselector.MATLABManager), new(*matlabmanager.MATLABManager)),

		// MATLAB Starting Dir Selector
		matlabstartingdirselector.New,
		wire.Bind(new(matlabstartingdirselector.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(matlabstartingdirselector.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(matlabstartingdirselector.RootStore), new(*rootstore.RootStore)),
		wire.Bind(new(matlabstartingdirselector.RootPathResolver), new(*rootpathresolver.RootPathResolver)),

		// Entities
		wire.Bind(new(entities.GlobalMATLAB), new(*globalmatlab.GlobalMATLAB)),
		wire.Bind(new(entities.MATLABManager), new(*matlabmanager.MATLABManager)),

		// MATLAB Manager
		matlabmanager.New,
		wire.Bind(new(matlabmanager.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(matlabmanager.MATLABServices), new(*matlabservices.MATLABServices)),
		wire.Bind(new(matlabmanager.MATLABSessionStore), new(*matlabsessionstore.Store)),
		wire.Bind(new(matlabmanager.MATLABSessionClientFactory), new(*matlabsessionclient.Factory)),
		wire.Bind(new(matlabmanager.SessionSelector), new(*sessionselector.SessionSelector)),

		// Session Selector
		sessionselector.New,
		wire.Bind(new(sessionselector.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(sessionselector.SessionDiscoverer), new(*sessiondiscovery.SessionDiscoverer)),

		// MATLAB Services
		matlabservices.New,
		wire.Bind(new(matlabservices.MATLABLocator), new(*matlablocator.MATLABLocator)),
		wire.Bind(new(matlabservices.LocalMATLABSessionLauncher), new(*localmatlabsession.Starter)),

		// MATLAB Locator
		matlablocator.New,
		wire.Bind(new(matlablocator.MATLABRootGetter), new(*matlabroot.Getter)),
		wire.Bind(new(matlablocator.MATLABVersionGetter), new(*matlabversion.Getter)),

		// MATLAB Root Getter
		matlabroot.New,
		wire.Bind(new(matlabroot.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(matlabroot.FileLayer), new(*filefacade.FileFacade)),

		// MATLAB Version Getter
		matlabversion.New,
		wire.Bind(new(matlabversion.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(matlabversion.IOLayer), new(*iofacade.IoFacade)),

		// Local MATLAB Session
		localmatlabsession.NewStarter,
		wire.Bind(new(localmatlabsession.SessionDirectoryFactory), new(*localmatlabsessiondirectory.Factory)),
		wire.Bind(new(localmatlabsession.ProcessDetails), new(*processdetails.ProcessDetails)),
		wire.Bind(new(localmatlabsession.MATLABProcessLauncher), new(*processlauncher.MATLABProcessLauncher)),
		wire.Bind(new(localmatlabsession.Watchdog), new(*watchdogclient.Watchdog)),

		// Local MATLAB Session Directory
		localmatlabsessiondirectory.NewFactory,
		wire.Bind(new(localmatlabsessiondirectory.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(localmatlabsessiondirectory.ApplicationDirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(localmatlabsessiondirectory.MATLABFiles), new(matlabfiles.MATLABFiles)),
		wire.Bind(new(localmatlabsessiondirectory.ConfigFactory), new(*config.Factory)),

		// MATLAB Files Provider
		matlabfiles.New,

		// Local MATLAB Session Process Details
		processdetails.New,
		wire.Bind(new(processdetails.OSLayer), new(*osfacade.OsFacade)),

		// Local MATLAB Process Launcher
		processlauncher.New,

		// MATLAB Session Store
		matlabsessionstore.New,
		wire.Bind(new(matlabsessionstore.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(matlabsessionstore.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),

		// MATLAB Session Client Factory
		matlabsessionclient.NewFactory,
		wire.Bind(new(matlabsessionclient.HttpClientFactory), new(*httpclient.Factory)),

		// Session Discovery
		sessiondiscovery.New,
		wire.Bind(new(sessiondiscovery.AppDataDirGetter), new(*appdatadir.Getter)),
		wire.Bind(new(sessiondiscovery.OSLayer), new(*osfacade.OsFacade)),

		// App Data Dir
		appdatadir.New,
		wire.Bind(new(appdatadir.OSLayer), new(*osfacade.OsFacade)),

		// Shared Dependencies

		// Path Validator
		pathvalidator.New,
		wire.Bind(new(pathvalidator.OSLayer), new(*osfacade.OsFacade)),

		// Process Handler
		processhandler.New,
		wire.Bind(new(processhandler.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(processhandler.OSWrapper), new(*osadaptor.ProcessManager)),

		// HTTP Server Handler Factory
		handler.NewFactory,
		wire.Bind(new(handler.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(handler.ProcessHandler), new(*processhandler.ProcessHandler)),

		// Socket Factory
		socket.NewFactory,
		wire.Bind(new(socket.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(socket.OSLayer), new(*osfacade.OsFacade)),

		// Logger Factory
		logger.NewFactory,
		wire.Bind(new(logger.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(logger.DirectoryFactory), new(*directory.Factory)),
		wire.Bind(new(logger.FilenameFactory), new(*files.Factory)),
		wire.Bind(new(logger.OSLayer), new(*osfacade.OsFacade)),

		// Directory Factory
		directory.NewFactory,
		wire.Bind(new(directory.ConfigFactory), new(*config.Factory)),
		wire.Bind(new(directory.FilenameFactory), new(*files.Factory)),
		wire.Bind(new(directory.OSLayer), new(*osfacade.OsFacade)),

		// Lifecycle Signaler
		lifecyclesignaler.New,

		// Config Factory
		config.NewFactory,
		wire.Bind(new(config.Parser), new(*parser.Parser)),
		wire.Bind(new(config.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(config.BuildInfo), new(*buildinfo.BuildInfo)),

		// BuildInfo
		buildinfo.New,
		wire.Bind(new(buildinfo.OSLayer), new(*osfacade.OsFacade)),

		// Parser
		parser.New,
		wire.Bind(new(parser.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(parser.DefaultParameterFactory), new(*selector.Selector)),
		wire.Bind(new(parser.ParameterFactory), new(ApplicationDefinition)),

		// Default Parameters Selector
		selector.New,
		wire.Bind(new(selector.ApplicationDefinition), new(ApplicationDefinition)),
		wire.Bind(new(selector.MessageCatalog), new(*messagecatalog.MessageCatalog)),

		// Message Catalog
		messagecatalog.New,

		// Files Factory
		files.NewFactory,
		wire.Bind(new(files.OSLayer), new(*osfacade.OsFacade)),

		// HTTP Client Factory
		httpclient.NewFactory,

		// Process Manager
		osadaptor.NewProcessManager,
		wire.Bind(new(osadaptor.OSLayer), new(*osfacade.OsFacade)),

		// OS Adaptor
		osadaptor.New,
		wire.Bind(new(osadaptor.VersionOSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(osadaptor.RegistryLayer), new(*registryfacade.RegistryFacade)),

		// Facades
		osfacade.New,
		iofacade.New,
		filefacade.New,
		registryfacade.New,
		unixfacade.New,

		// Resource Limit
		resourcelimit.New,
		wire.Bind(new(resourcelimit.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(resourcelimit.SyscallLayer), new(*unixfacade.UnixFacade)),
	)

	return nil
}
