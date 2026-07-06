// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultCacheDir = ".aifunc"
	PackagesDir     = "packages"
	EngineDir       = "_engine"
	ConfigFile      = "aifunc.json"
	LockFile        = "aifunc-lock.json"
	GitIgnoreFile   = ".gitignore"
)

type Workspace struct {
	Root     string
	InputDir string
}

func FromCurrentDir() (*Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Workspace{Root: cwd, InputDir: DefaultCacheDir}, nil
}

func (w *Workspace) SetInputDir(inputDir string) {
	if inputDir != "" {
		w.InputDir = inputDir
	}
}

// ConfigPath returns the path to aifunc.json at the project root.
func (w *Workspace) ConfigPath() string {
	return filepath.Join(w.Root, ConfigFile)
}

// LockPath returns the path to aifunc-lock.json at the project root.
func (w *Workspace) LockPath() string {
	return filepath.Join(w.Root, LockFile)
}

// GitIgnorePath returns the path to .gitignore at the project root.
func (w *Workspace) GitIgnorePath() string {
	return filepath.Join(w.Root, GitIgnoreFile)
}

// CachePath returns the root of the download cache directory (inputDir).
func (w *Workspace) CachePath() string {
	return filepath.Join(w.Root, w.InputDir)
}

// PackagesPath returns the path where downloaded packages are cached.
func (w *Workspace) PackagesPath() string {
	return filepath.Join(w.CachePath(), PackagesDir)
}

// EngineCachePath returns the path where downloaded engine SDKs are cached.
func (w *Workspace) EngineCachePath() string {
	return filepath.Join(w.CachePath(), EngineDir)
}

// EnsureCache creates the .aifunc/ cache directories if they don't exist.
func (w *Workspace) EnsureCache() error {
	for _, dir := range []string{
		w.CachePath(),
		w.PackagesPath(),
		w.EngineCachePath(),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// RequireConfig checks that aifunc.json exists, returning a helpful error if not.
func (w *Workspace) RequireConfig() error {
	_, err := os.Stat(w.ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("aifunc.json not found, run aifunc init first")
		}
		return err
	}
	return nil
}
