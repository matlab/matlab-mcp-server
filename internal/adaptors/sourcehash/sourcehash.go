// Copyright 2026 The MathWorks, Inc.

package sourcehash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// MismatchError is returned by Check when the stored hash does not match the computed hash.
type MismatchError struct {
	Stored  string
	Current string
}

func (e *MismatchError) Error() string {
	return fmt.Sprintf("hash mismatch: stored %s, current %s", e.Stored, e.Current)
}

// ComputeHash produces a deterministic SHA-256 hash over all git-tracked files
// in sourceDir. Each file contributes its relative path (forward-slash normalized)
// and its content to the hash, ensuring renames and content changes are detected.
func ComputeHash(sourceDir string) (string, error) {
	files, err := gitTrackedFiles(sourceDir)
	if err != nil {
		return "", fmt.Errorf("listing tracked files: %w", err)
	}

	hasher := sha256.New()
	for _, relPath := range files {
		if _, err := fmt.Fprintf(hasher, "file:%s\n", filepath.ToSlash(relPath)); err != nil {
			return "", fmt.Errorf("hashing path %s: %w", relPath, err)
		}

		fullPath := filepath.Join(sourceDir, filepath.FromSlash(relPath))
		f, err := os.Open(fullPath) //nolint:gosec // file paths come from git ls-files, not user input
		if err != nil {
			return "", fmt.Errorf("opening %s: %w", relPath, err)
		}
		if _, err := io.Copy(hasher, f); err != nil {
			_ = f.Close()
			return "", fmt.Errorf("reading %s: %w", relPath, err)
		}
		_ = f.Close()
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Write computes the hash of sourceDir and writes it to hashFile.
func Write(hashFile, sourceDir string) error {
	hash, err := ComputeHash(sourceDir)
	if err != nil {
		return err
	}

	if err := os.WriteFile(hashFile, []byte(hash+"\n"), 0o600); err != nil {
		return fmt.Errorf("failed to write hash file: %w", err)
	}

	return nil
}

// Check computes the hash of sourceDir and compares it to the hash stored in hashFile.
// Returns a *MismatchError if the hashes differ.
func Check(hashFile, sourceDir string) error {
	hash, err := ComputeHash(sourceDir)
	if err != nil {
		return err
	}

	stored, err := os.ReadFile(hashFile) //nolint:gosec // hash file path comes from trusted caller
	if err != nil {
		return fmt.Errorf("failed to read hash file %s: %w", hashFile, err)
	}

	if strings.TrimSpace(string(stored)) != hash {
		return &MismatchError{
			Stored:  strings.TrimSpace(string(stored)),
			Current: hash,
		}
	}

	return nil
}

func gitTrackedFiles(dir string) ([]string, error) {
	cmd := exec.Command("git", "-C", dir, "ls-files", ".")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w", err)
	}

	var files []string
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		files = append(files, line)
	}
	sort.Strings(files)
	return files, nil
}
