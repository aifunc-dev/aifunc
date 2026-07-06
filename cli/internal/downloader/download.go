// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package downloader

import (
	"fmt"
	"os"
	"path/filepath"

	"aifunc/cli/internal/manifest"
	"aifunc/cli/internal/source"
	"aifunc/cli/internal/workspace"
)

type Result struct {
	Name string
	// RepoRoot is the root of the cloned repository (for git sources).
	// It is valid only until Cleanup is called.
	RepoRoot string
	Cleanup  func()
}

// Download fetches a remote package into the workspace cache.
// For local (file:) sources, use ValidateLocal instead -- local packages are
// never copied; they are referenced directly from their source path.
func Download(raw string, ws *workspace.Workspace) (Result, error) {
	if err := os.MkdirAll(ws.PackagesPath(), 0755); err != nil {
		return Result{}, err
	}

	repoRoot, packageDir, cleanup, err := fetchToStaging(raw)
	if err != nil {
		return Result{}, err
	}

	spec, err := manifest.Validate(packageDir)
	if err != nil {
		cleanup()
		return Result{}, err
	}
	dest := filepath.Join(ws.PackagesPath(), spec.Name)
	if err := source.CopyLocal(packageDir, dest); err != nil {
		cleanup()
		return Result{}, err
	}
	return Result{Name: spec.Name, RepoRoot: repoRoot, Cleanup: cleanup}, nil
}

// ValidateLocal validates a local package directory and returns its name and
// absolute path. Nothing is copied into the workspace cache.
func ValidateLocal(rawPath string) (name string, absPath string, err error) {
	abs, err := filepath.Abs(rawPath)
	if err != nil {
		return "", "", fmt.Errorf("resolve path %q: %w", rawPath, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("local package path does not exist: %s", abs)
		}
		return "", "", err
	}
	if !info.IsDir() {
		return "", "", fmt.Errorf("local package path is not a directory: %s", abs)
	}
	spec, err := manifest.Validate(abs)
	if err != nil {
		return "", "", err
	}
	return spec.Name, abs, nil
}

// fetchToStaging clones the git source into a temp directory and returns:
//   - repoRoot: the repository root
//   - packageDir: the package subdirectory within the repo
//   - cleanup: function to remove the temp directory
func fetchToStaging(raw string) (repoRoot string, packageDir string, cleanup func(), err error) {
	noop := func() {}

	src, parseErr := source.Parse(raw)
	if parseErr != nil {
		return "", "", noop, parseErr
	}

	if src.Kind == source.KindLocal {
		return "", "", noop, fmt.Errorf("local sources must use ValidateLocal, not Download")
	}

	tmp, mkErr := os.MkdirTemp("", "aifn-download-*")
	if mkErr != nil {
		return "", "", noop, mkErr
	}
	cleanupFn := func() { _ = os.RemoveAll(tmp) }

	switch src.Kind {
	case source.KindGit:
		clonedRepo := filepath.Join(tmp, "repo")
		if gitErr := source.CloneRepo(src, clonedRepo); gitErr != nil {
			cleanupFn()
			return "", "", noop, gitErr
		}
		pkgDir := clonedRepo
		if src.SubPath != "" {
			pkgDir = filepath.Join(clonedRepo, filepath.FromSlash(src.SubPath))
		}
		return clonedRepo, pkgDir, cleanupFn, nil

	default:
		cleanupFn()
		return "", "", noop, fmt.Errorf("unsupported source kind %q", src.Kind)
	}
}
