// Copyright 2026 The MathWorks, Inc.

package installmatlab

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

// FileDownloader downloads a file from a URL to a local path.
type FileDownloader interface {
	Download(ctx context.Context, url string, destPath string) error
}

// CommandRunner executes a command and returns its combined output.
type CommandRunner interface {
	Run(ctx context.Context, name string, args []string) (string, error)
}

// OSLayer provides OS-level operations needed by the usecase.
type OSLayer interface {
	GOOS() string
	GOARCH() string
	UserHomeDir() (string, error)
	MkdirAll(path string, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Chmod(name string, mode os.FileMode) error
}

type Args struct {
	Release     string
	Destination string
	Products    []string
}

type ReturnArgs struct {
	Output string
}

type Usecase struct {
	downloader    FileDownloader
	commandRunner CommandRunner
	osLayer       OSLayer
}

func New(
	downloader FileDownloader,
	commandRunner CommandRunner,
	osLayer OSLayer,
) *Usecase {
	return &Usecase{
		downloader:    downloader,
		commandRunner: commandRunner,
		osLayer:       osLayer,
	}
}

func (u *Usecase) Execute(ctx context.Context, logger entities.Logger, args Args) (ReturnArgs, error) {
	logger.Debug("Entering InstallMATLAB Usecase")
	defer logger.Debug("Exiting InstallMATLAB Usecase")

	if args.Release == "" {
		return ReturnArgs{}, fmt.Errorf("release is required (e.g., R2025a)")
	}
	if len(args.Products) == 0 {
		return ReturnArgs{}, fmt.Errorf("at least one product name is required")
	}

	mpmPath, err := u.ensureMpm(ctx, logger)
	if err != nil {
		return ReturnArgs{}, fmt.Errorf("failed to ensure mpm is available: %w", err)
	}

	output, err := u.runMpmInstall(ctx, logger, mpmPath, args)
	if err != nil {
		return ReturnArgs{Output: output}, err
	}

	return ReturnArgs{Output: output}, nil
}

func (u *Usecase) ensureMpm(ctx context.Context, logger entities.Logger) (string, error) {
	mpmPath, err := u.getMpmPath()
	if err != nil {
		return "", err
	}

	if _, err := u.osLayer.Stat(mpmPath); err == nil {
		logger.Info("mpm already present at " + mpmPath)
		return mpmPath, nil
	}

	downloadURL, err := u.getMpmDownloadURL()
	if err != nil {
		return "", err
	}

	logger.With("url", downloadURL).Info("Downloading mpm")

	if err := u.downloader.Download(ctx, downloadURL, mpmPath); err != nil {
		return "", fmt.Errorf("failed to download mpm: %w", err)
	}

	if u.osLayer.GOOS() != "windows" {
		if err := u.osLayer.Chmod(mpmPath, 0755); err != nil {
			return "", fmt.Errorf("failed to make mpm executable: %w", err)
		}
	}

	logger.Info("mpm downloaded to " + mpmPath)
	return mpmPath, nil
}

func (u *Usecase) getMpmPath() (string, error) {
	homeDir, err := u.osLayer.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	mpmDir := filepath.Join(homeDir, ".mathworks", "mpm")
	if err := u.osLayer.MkdirAll(mpmDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create mpm directory: %w", err)
	}

	binaryName := "mpm"
	if u.osLayer.GOOS() == "windows" {
		binaryName = "mpm.exe"
	}

	return filepath.Join(mpmDir, binaryName), nil
}

func (u *Usecase) getMpmDownloadURL() (string, error) {
	goos := u.osLayer.GOOS()
	goarch := u.osLayer.GOARCH()

	switch {
	case goos == "windows":
		return "https://www.mathworks.com/mpm/win64/mpm", nil
	case goos == "darwin" && goarch == "arm64":
		return "https://www.mathworks.com/mpm/maca64/mpm", nil
	case goos == "darwin" && goarch == "amd64":
		return "https://www.mathworks.com/mpm/maci64/mpm", nil
	case goos == "linux":
		return "https://www.mathworks.com/mpm/glnxa64/mpm", nil
	default:
		return "", fmt.Errorf("unsupported platform: %s/%s", goos, goarch)
	}
}

func (u *Usecase) runMpmInstall(ctx context.Context, logger entities.Logger, mpmPath string, args Args) (string, error) {
	cmdArgs := []string{
		"install",
		"--release=" + args.Release,
	}
	if args.Destination != "" {
		cmdArgs = append(cmdArgs, "--destination="+args.Destination)
	}
	cmdArgs = append(cmdArgs, "--products")
	cmdArgs = append(cmdArgs, args.Products...)

	logger.With("command", mpmPath+" "+strings.Join(cmdArgs, " ")).Info("Running mpm install")

	output, err := u.commandRunner.Run(ctx, mpmPath, cmdArgs)
	if err != nil {
		return output, fmt.Errorf("mpm install failed: %w", err)
	}

	return output, nil
}
