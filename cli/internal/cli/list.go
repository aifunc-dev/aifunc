// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"aifunc/cli/internal/lockfile"
	"aifunc/cli/internal/types"
	"aifunc/cli/internal/workspace"

	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List installed packages",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
	}
}

func runList() error {
	ws, err := workspace.FromCurrentDir()
	if err != nil {
		return err
	}

	if err := ws.RequireConfig(); err != nil {
		return err
	}

	cfg, err := readConfig(ws.ConfigPath())
	if err != nil {
		return fmt.Errorf("reading aifunc.json: %w", err)
	}
	ws.SetInputDir(cfg.GetInputDir())

	if len(cfg.Packages) == 0 {
		fmt.Fprintln(os.Stdout, "No packages declared in this project.")
		fmt.Fprintln(os.Stdout, "Use aifn install <source> to install a package, or add entries to aifunc.json manually.")
		return nil
	}

	lock, hasLock := tryReadLock(ws)

	names := make([]string, 0, len(cfg.Packages))
	for name := range cfg.Packages {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Fprintf(os.Stdout, "%d package(s) declared (language: %s, output: %s)\n\n", len(names), cfg.Language, cfg.GetOutputDir())

	for _, name := range names {
		source := cfg.Packages[name]
		if hasLock {
			if pkgLock, ok := lock.Packages[name]; ok {
				fmt.Fprintf(os.Stdout, "  %s\n", name)
				fmt.Fprintf(os.Stdout, "    version: %s | engine: v%s\n", pkgLock.Version, pkgLock.EngineVersion)
				fmt.Fprintf(os.Stdout, "    source:  %s\n", source)
				continue
			}
		}
		fmt.Fprintf(os.Stdout, "  %s (not installed)\n", name)
		fmt.Fprintf(os.Stdout, "    source:  %s\n", source)
	}

	return nil
}

func tryReadLock(ws *workspace.Workspace) (types.LockFile, bool) {
	if _, err := os.Stat(ws.LockPath()); errors.Is(err, os.ErrNotExist) {
		return types.LockFile{}, false
	}
	lf, err := lockfile.Read(ws.LockPath())
	if err != nil {
		return types.LockFile{}, false
	}
	return lf, true
}
