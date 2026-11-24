// Copyright 2026 The MathWorks, Inc.

package toolboxcache_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/matlab/matlab-mcp-server/tests/testutils/toolboxcache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClear_DeletesTheToolboxCacheOfTheReleaseUnderTest(t *testing.T) {
	// Arrange
	localAppData := t.TempDir()
	matlabRoot := newMATLABRoot(t, "R2027a")
	cache := writeFile(t, localAppData, "R2027a", "toolbox_cache-27.1.0-1643244820-win64.xml")
	deployedCache := writeFile(t, localAppData, "R2027a", "toolbox_cache-mcr-27.1.0-1643244820-win64.xml")
	settings := writeFile(t, localAppData, "R2027a", "matlab.mlsettings")
	otherReleaseCache := writeFile(t, localAppData, "R2026b", "toolbox_cache-26.2.0-1643244820-win64.xml")

	// Act
	deleted, err := toolboxcache.Clear(localAppData, matlabRoot)

	// Assert
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{cache, deployedCache}, deleted)
	assert.NoFileExists(t, cache)
	assert.NoFileExists(t, deployedCache)
	assert.FileExists(t, settings, "only the toolbox path cache should be deleted")
	assert.FileExists(t, otherReleaseCache, "the cache of another release should stay")
}

func TestClear_DeletesTheCacheWhenTheLocalAppDataPathHasSquareBrackets(t *testing.T) {
	// Arrange
	localAppData := filepath.Join(t.TempDir(), "user[1]")
	matlabRoot := newMATLABRoot(t, "R2027a")
	cache := writeFile(t, localAppData, "R2027a", "toolbox_cache-27.1.0-1643244820-win64.xml")

	// Act
	deleted, err := toolboxcache.Clear(localAppData, matlabRoot)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{cache}, deleted)
	assert.NoFileExists(t, cache)
}

func TestClear_DoesNothingWhenTheReleaseHasNoCache(t *testing.T) {
	// Arrange
	localAppData := t.TempDir()
	matlabRoot := newMATLABRoot(t, "R2027a")

	// Act
	deleted, err := toolboxcache.Clear(localAppData, matlabRoot)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, deleted)
}

func TestClear_ReturnsAnErrorWhenVersionInfoIsMissing(t *testing.T) {
	// Arrange
	localAppData := t.TempDir()
	cache := writeFile(t, localAppData, "R2027a", "toolbox_cache-27.1.0-1643244820-win64.xml")

	// Act
	deleted, err := toolboxcache.Clear(localAppData, t.TempDir())

	// Assert
	require.ErrorIs(t, err, fs.ErrNotExist)
	assert.Empty(t, deleted)
	assert.FileExists(t, cache)
}

func TestClear_ReturnsAnErrorWhenVersionInfoIsNotXML(t *testing.T) {
	// Arrange
	localAppData := t.TempDir()
	matlabRoot := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(matlabRoot, "VersionInfo.xml"), []byte("<MathWorks_version_info>"), 0o600))

	// Act
	deleted, err := toolboxcache.Clear(localAppData, matlabRoot)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "VersionInfo.xml")
	assert.Empty(t, deleted)
}

func TestClear_RejectsAReleaseThatIsNotAMATLABRelease(t *testing.T) {
	for _, release := range []string{"", "../R2027a", "R2027*"} {
		t.Run(release, func(t *testing.T) {
			// Arrange
			localAppData := t.TempDir()
			matlabRoot := newMATLABRoot(t, release)

			// Act
			deleted, err := toolboxcache.Clear(localAppData, matlabRoot)

			// Assert
			require.ErrorIs(t, err, toolboxcache.ErrUnexpectedRelease)
			assert.Empty(t, deleted)
		})
	}
}

func TestClear_ReturnsAnErrorWhenTheLocalAppDataFolderIsNotSet(t *testing.T) {
	// Arrange
	matlabRoot := newMATLABRoot(t, "R2027a")

	// Act
	deleted, err := toolboxcache.Clear("", matlabRoot)

	// Assert
	require.ErrorIs(t, err, toolboxcache.ErrLocalAppDataNotSet)
	assert.Empty(t, deleted)
}

func TestClear_ReportsAFileThatItCannotDelete(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a read-only folder does not stop a delete on Windows")
	}

	// Arrange
	localAppData := t.TempDir()
	matlabRoot := newMATLABRoot(t, "R2027a")
	cache := writeFile(t, localAppData, "R2027a", "toolbox_cache-27.1.0-1643244820-win64.xml")
	cacheFolder := filepath.Dir(cache)
	require.NoError(t, os.Chmod(cacheFolder, 0o500))       //nolint:gosec // A directory needs the execute bit
	t.Cleanup(func() { _ = os.Chmod(cacheFolder, 0o700) }) //nolint:gosec // A directory needs the execute bit

	// Act
	deleted, err := toolboxcache.Clear(localAppData, matlabRoot)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), cache)
	assert.Empty(t, deleted)
}

func newMATLABRoot(t *testing.T, release string) string {
	t.Helper()
	matlabRoot := t.TempDir()
	versionInfo := `<?xml version="1.0" encoding="UTF-8"?>
<MathWorks_version_info>
  <version>27.1.0.3422810</version>
  <release>` + release + `</release>
</MathWorks_version_info>
`
	require.NoError(t, os.WriteFile(filepath.Join(matlabRoot, "VersionInfo.xml"), []byte(versionInfo), 0o600))
	return matlabRoot
}

func writeFile(t *testing.T, localAppData, release, name string) string {
	t.Helper()
	folder := filepath.Join(localAppData, "MathWorks", "MATLAB", release)
	require.NoError(t, os.MkdirAll(folder, 0o700))
	path := filepath.Join(folder, name)
	require.NoError(t, os.WriteFile(path, nil, 0o600))
	return path
}
