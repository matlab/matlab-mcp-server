// Copyright 2026 The MathWorks, Inc.

package toolboxcache

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

const cacheFilePattern = "toolbox_cache-*.xml"

var releasePattern = regexp.MustCompile(`^R\d{4}[ab]$`)

var (
	ErrLocalAppDataNotSet = errors.New("the local app data folder is not set")
	ErrUnexpectedRelease  = errors.New("unexpected MATLAB release")
)

type versionInfo struct {
	Release string `xml:"release"`
}

func Clear(localAppData, matlabRoot string) ([]string, error) {
	if localAppData == "" {
		return nil, ErrLocalAppDataNotSet
	}

	release, err := readRelease(matlabRoot)
	if err != nil {
		return nil, err
	}

	cacheFiles, err := findCacheFiles(filepath.Join(localAppData, "MathWorks", "MATLAB", release))
	if err != nil {
		return nil, err
	}

	var deleted []string
	var removeErrs []error
	for _, cacheFile := range cacheFiles {
		if err := os.Remove(cacheFile); err != nil {
			removeErrs = append(removeErrs, fmt.Errorf("failed to delete %s: %w", cacheFile, err))
			continue
		}
		deleted = append(deleted, cacheFile)
	}

	return deleted, errors.Join(removeErrs...)
}

func findCacheFiles(releaseFolder string) ([]string, error) {
	entries, err := os.ReadDir(releaseFolder)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find the toolbox path cache: %w", err)
	}

	var cacheFiles []string
	for _, entry := range entries {
		if matched, _ := filepath.Match(cacheFilePattern, entry.Name()); matched && !entry.IsDir() {
			cacheFiles = append(cacheFiles, filepath.Join(releaseFolder, entry.Name()))
		}
	}
	return cacheFiles, nil
}

func readRelease(matlabRoot string) (string, error) {
	versionInfoPath := filepath.Join(matlabRoot, "VersionInfo.xml")
	data, err := os.ReadFile(versionInfoPath) //nolint:gosec // Trusted test path
	if err != nil {
		return "", fmt.Errorf("failed to read the MATLAB release: %w", err)
	}

	var info versionInfo
	if err := xml.Unmarshal(data, &info); err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", versionInfoPath, err)
	}

	if !releasePattern.MatchString(info.Release) {
		return "", fmt.Errorf("%w %q in %s", ErrUnexpectedRelease, info.Release, versionInfoPath)
	}

	return info.Release, nil
}
