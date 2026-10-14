// Copyright 2026 The MathWorks, Inc.

package system_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-server/internal/adaptors/globalmatlab/sessionmanager"
	"github.com/matlab/matlab-mcp-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabversion"
	"github.com/matlab/matlab-mcp-server/internal/facades/iofacade"
	"github.com/matlab/matlab-mcp-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-server/tests/system/testdata"
	"github.com/matlab/matlab-mcp-server/tests/testutils/externalmatlab"
	"github.com/matlab/matlab-mcp-server/tests/testutils/fileuri"
	"github.com/matlab/matlab-mcp-server/tests/testutils/matlabtempdir"
	"github.com/matlab/matlab-mcp-server/tests/testutils/mcpclient"
	"github.com/matlab/matlab-mcp-server/tests/testutils/pathcontrol"
	"github.com/matlab/matlab-mcp-server/tests/testutils/sessiondetails"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/suite"
)

const (
	sessionAValue = 42
	sessionBValue = 99

	fallbackLogMessage = "falling back to launching MATLAB"

	failFastDiscovery = "--matlab-session-discovery-timeout=2s"

	// Must stay equal to MinimumMatlabRelease in buildfile.m.
	minimumToolboxRelease = "R2023a"
)

type ExistingSessionTestSuite struct {
	SystemTestSuite
	testHome string
}

func (s *ExistingSessionTestSuite) SetupSuite() {
	s.SystemTestSuite.SetupSuite()
	s.skipIfToolboxUnsupported()
	s.testHome = matlabtempdir.New(s.T(), "existing-session-home-")
	s.Require().NoError(os.MkdirAll(s.testPrefDir(), 0o700))
	// Keeps a userpath warning out of the tool output.
	s.Require().NoError(os.MkdirAll(filepath.Join(s.testHome, "Documents", "MATLAB"), 0o700))
	s.installToolbox()
}

func (s *ExistingSessionTestSuite) SetupTest() {
	s.SystemTestSuite.SetupTest()
	s.defaultEnv = s.isolatedEnv(s.defaultEnv)
	s.removeSharedSessionFile()
}

func (s *ExistingSessionTestSuite) TestExistingSessionConnectionLifecycle() {
	ctx := s.T().Context()

	// Step 1: attach in existing mode while no MATLAB is shared.
	s.inSession(ctx, allowErrorLogs, []string{existingSessionMode, failFastDiscovery}, func(session *SystemSession) {
		_, err := session.EvaluateCode(ctx, `disp('no session yet')`)
		s.Require().ErrorContains(err, sessionmanager.ErrFailedToAttachToMATLABSession.Error(),
			"existing mode without a shared session should fail to attach")
	})

	// Step 2: the user shares MATLAB A. Attach in existing mode and run each single-session tool.
	matlabA := s.startExternalMATLAB(sessionAValue)
	matlabA.Share(s.T())
	s.inSession(ctx, assertNoErrorLogs, []string{existingSessionMode}, func(session *SystemSession) {
		s.assertReadsSessionValue(ctx, session, sessionAValue)
		s.exerciseSingleSessionTools(ctx, session)
	})
	matlabA.AssertAlive(s.T())

	// Step 3: attach in auto mode while MATLAB A is shared.
	s.assertAttachesTo(ctx, autoSessionMode, sessionAValue)
	matlabA.AssertAlive(s.T())

	// Step 4: the user starts MATLAB B and does not share it.
	matlabB := s.startExternalMATLAB(sessionBValue)
	s.assertAttachesTo(ctx, existingSessionMode, sessionAValue)

	// Step 5: the user shares MATLAB B, which overwrites the published session details.
	matlabB.Share(s.T())
	s.assertAttachesTo(ctx, existingSessionMode, sessionBValue)

	// Step 6: confirm both external MATLABs are still running.
	matlabA.AssertAlive(s.T())
	matlabB.AssertAlive(s.T())
}

func (s *ExistingSessionTestSuite) TestAutoModeFallbackLaunchesLocalMATLAB() {
	ctx := s.T().Context()

	session := s.createSession(ctx, autoSessionMode)
	defer s.CleanupSession(session, assertNoErrorLogs)

	output, err := session.EvaluateCode(ctx, `disp('auto fallback')`)
	s.Require().NoError(err, "auto mode should launch a local MATLAB when no session is shared")
	s.Contains(output, "auto fallback")
	s.assertServerLogContains(session, fallbackLogMessage)
}

func TestExistingSessionSuite(t *testing.T) {
	suite.Run(t, new(ExistingSessionTestSuite))
}

func (s *ExistingSessionTestSuite) skipIfToolboxUnsupported() {
	version, err := matlabversion.New(osfacade.New(), iofacade.New()).Get(s.matlabRoot())
	s.Require().NoError(err, "should read the MATLAB release under test")

	if releaseNumber(version.ReleaseFamily) < releaseNumber(minimumToolboxRelease) {
		s.T().Skipf("the MATLAB MCP Server Toolbox needs %s or later, and the MATLAB under test is %s",
			minimumToolboxRelease, version.ReleaseFamily)
	}
}

func releaseNumber(release string) string {
	return strings.TrimPrefix(release, "R")
}

func (s *ExistingSessionTestSuite) installToolbox() {
	_, envWithMATLAB := s.matlabEnv()
	installCmd := exec.CommandContext(s.T().Context(), s.mcpServerPath, "--setup-matlab") //nolint:gosec // Trusted path in tests
	installCmd.Env = s.isolatedEnv(envWithMATLAB)
	out, err := installCmd.CombinedOutput()
	s.Require().NoError(err, "setup-matlab should succeed:\n%s", out)
}

func (s *ExistingSessionTestSuite) removeSharedSessionFile() {
	if err := sessiondetails.Remove(s.testHome); !errors.Is(err, fs.ErrNotExist) {
		s.Require().NoError(err, "should remove the shared session details file, or the test attaches to a stale session")
	}
}

func (s *ExistingSessionTestSuite) isolatedEnv(env []string) []string {
	env = pathcontrol.UpdateEnvEntry(env, "HOME", s.testHome)
	env = pathcontrol.UpdateEnvEntry(env, "APPDATA", filepath.Join(s.testHome, "AppData", "Roaming"))
	return pathcontrol.UpdateEnvEntry(env, "MATLAB_PREFDIR", s.testPrefDir())
}

func (s *ExistingSessionTestSuite) testPrefDir() string {
	return filepath.Join(s.testHome, "MathWorks", "MATLAB", "prefdir")
}

func (s *ExistingSessionTestSuite) assertReadsSessionValue(ctx context.Context, session *SystemSession, want int) {
	s.T().Helper()
	output, err := session.EvaluateCode(ctx, `fprintf('val=%d\n', sessionValue);`)
	s.Require().NoError(err, "should read the attached session's workspace variable")
	s.Contains(strings.Split(output, "\n"), fmt.Sprintf("val=%d", want),
		"should attach to the session holding sessionValue=%d", want)
}

func (s *ExistingSessionTestSuite) exerciseSingleSessionTools(ctx context.Context, session *SystemSession) {
	s.T().Helper()

	guidelines, err := session.ReadResource(ctx, "guidelines://coding")
	s.Require().NoError(err, "should read coding guidelines resource")
	s.Contains(guidelines, "MATLAB Coding Standards", "should contain coding standards title")

	info, err := session.DetectToolboxes(ctx)
	s.Require().NoError(err, "should detect toolboxes")
	s.Contains(info, "MATLAB Version:", "should discover MATLAB version")

	output, err := session.EvaluateCode(ctx, `a = int32(2); b = int32(3);`)
	s.Require().NoError(err)
	s.Empty(output, "semicolon-terminated statements should produce no output")

	output, err = session.EvaluateCode(ctx, `c = a + b; fprintf('c=%d\n', c);`)
	s.Require().NoError(err)
	s.Contains(strings.Split(output, "\n"), "c=5", "should compute 2 + 3 = 5")

	issues, err := session.CheckCode(ctx, s.problematicCodePath())
	s.Require().NoError(err, "should check code without error")
	testdata.AssertProblematicCodeIssues(s.T(), issues)

	scriptOutput, err := session.RunFile(ctx, s.testScriptPath())
	s.Require().NoError(err, "should execute script file without error")
	testdata.TestScript.Assert(s.T(), scriptOutput)

	testOutput, err := session.RunTestFile(ctx, s.testMathFunctionsPath())
	s.Require().NoError(err, "should execute test suite without error")
	testdata.TestMathFunctions.Assert(s.T(), testOutput)
}

func (s *ExistingSessionTestSuite) assertServerLogContains(session *SystemSession, substr string) {
	s.T().Helper()
	s.Require().Eventually(func() bool {
		logs, err := session.ReadServerLogs()
		return err == nil && strings.Contains(logs, substr)
	}, 10*time.Second, 100*time.Millisecond, "server log should contain %q", substr)
}

func (s *ExistingSessionTestSuite) createSession(ctx context.Context, args ...string) *SystemSession {
	s.T().Helper()
	sessionRootedAtTestData := []mcpclient.CreateSessionOption{
		mcpclient.WithRoots(&mcp.Root{URI: fileuri.ToFileURI(s.testDataDir), Name: "test-data"}),
	}
	return s.CreateMCPSession(ctx, nil, sessionRootedAtTestData, args...)
}

func (s *ExistingSessionTestSuite) inSession(ctx context.Context, errorLogs errorLogPolicy, args []string, use func(*SystemSession)) {
	s.T().Helper()
	session := s.createSession(ctx, args...)
	defer s.CleanupSession(session, errorLogs)
	use(session)
}

func (s *ExistingSessionTestSuite) assertAttachesTo(ctx context.Context, sessionMode string, want int) {
	s.T().Helper()
	s.inSession(ctx, assertNoErrorLogs, []string{sessionMode}, func(session *SystemSession) {
		s.assertReadsSessionValue(ctx, session, want)
	})
}

func (s *ExistingSessionTestSuite) startExternalMATLAB(sessionValue int) *externalmatlab.Instance {
	s.T().Helper()
	prepCode := fmt.Sprintf("sessionValue = int32(%d);", sessionValue)
	return externalmatlab.Start(s.T(), s.matlabExe(), s.defaultEnv, prepCode)
}
