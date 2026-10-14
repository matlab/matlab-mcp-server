// Copyright 2026 The MathWorks, Inc.

package fileuri

import (
	"path/filepath"
	"strings"
)

func ToFileURI(dir string) string {
	slashPath := filepath.ToSlash(dir)
	if !strings.HasPrefix(slashPath, "/") {
		slashPath = "/" + slashPath
	}
	return "file://" + slashPath
}
