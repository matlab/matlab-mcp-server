// Copyright 2026 The MathWorks, Inc.

// sourcehash computes a deterministic SHA-256 hash of git-tracked source files
// in a directory, and can write or verify that hash against a stored hash file.
//
// Usage:
//
//	sourcehash write <hash-file> <source-dir>
//	sourcehash check <hash-file> <source-dir>
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/sourcehash"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: sourcehash <write|check> <hash-file> <source-dir>\n")
		os.Exit(2)
	}

	command := os.Args[1]
	hashFile := os.Args[2]
	sourceDir := os.Args[3]

	switch command {
	case "write":
		if err := sourcehash.Write(hashFile, sourceDir); err != nil {
			fmt.Fprintf(os.Stderr, "sourcehash: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Wrote source hash to %s\n", hashFile)

	case "check":
		if err := sourcehash.Check(hashFile, sourceDir); err != nil {
			var mismatch *sourcehash.MismatchError
			if errors.As(err, &mismatch) {
				fmt.Fprintf(os.Stderr, "sourcehash: embedded MATLAB addon is stale.\n")
				fmt.Fprintf(os.Stderr, "  stored:  %s\n", mismatch.Stored)
				fmt.Fprintf(os.Stderr, "  current: %s\n", mismatch.Current)
			} else {
				fmt.Fprintf(os.Stderr, "sourcehash: %v\n", err)
			}
			fmt.Fprintf(os.Stderr, "Run 'make update-embedded-matlab-addon' to update it.\n")
			os.Exit(1)
		}
		fmt.Println("Source hash OK")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s (expected 'write' or 'check')\n", command)
		os.Exit(2)
	}
}
