// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aifunc/cli/internal/lockfile"
	"aifunc/cli/internal/source"
	"aifunc/cli/internal/types"
	"aifunc/cli/internal/workspace"

	"github.com/spf13/cobra"
)

// newUninstallCommand creates the "uninstall" CLI command.
func newUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "uninstall <name...>",
		Aliases: []string{"rm"},
		Short:   "Remove installed packages",
		Long:    "Remove specified packages from the project: clear cached sources, generated output, and entries from aifunc.json and the lock file.",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall(args)
		},
	}
}

// --------------------------------------------------------------------
// Core uninstall logic
// --------------------------------------------------------------------

func runUninstall(names []string) error {
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

	lock, hasLock := tryReadLock(ws)

	removed := removePackages(ws, cfg, &lock, hasLock, names)
	if len(removed) == 0 {
		return nil
	}

	if err := persistState(ws, cfg, &lock, hasLock); err != nil {
		return err
	}

	cleanEmptyDirs(ws, cfg)
	fmt.Fprintf(os.Stdout, "\nremoved %d package(s).\n", len(removed))
	return nil
}

// removePackages iterates over requested names, resolves each one, and removes
// its cached sources + generated output. Returns the list of actually removed names.
func removePackages(ws *workspace.Workspace, cfg types.AifuncConfig, lock *types.LockFile, hasLock bool, names []string) []string {
	var removed []string

	for _, name := range names {
		resolved := resolveToConfigName(name, cfg)
		if resolved == "" {
			fmt.Fprintf(os.Stdout, "package %s is not declared in aifunc.json, skipping.\n", name)
			continue
		}

		fmt.Fprintf(os.Stdout, "removing %s ... ", resolved)
		removeSinglePackage(ws, cfg, lock, hasLock, resolved)
		fmt.Fprintln(os.Stdout, "done")

		removed = append(removed, resolved)
	}

	if hasLock && len(removed) > 0 {
		cleanOrphanedEngines(lock, ws, cfg)
	}

	return removed
}

// resolveToConfigName returns the canonical package name in cfg.Packages,
// or "" if it cannot be found.
func resolveToConfigName(input string, cfg types.AifuncConfig) string {
	if _, ok := cfg.Packages[input]; ok {
		return input
	}
	return resolvePackageName(input, cfg)
}

// removeSinglePackage deletes cache, output, and config/lock entries for one package.
func removeSinglePackage(ws *workspace.Workspace, cfg types.AifuncConfig, lock *types.LockFile, hasLock bool, name string) {
	// Remove cached source.
	os.RemoveAll(filepath.Join(ws.PackagesPath(), name))

	// Remove generated output.
	os.RemoveAll(filepath.Join(ws.Root, cfg.GetOutputDir(), outputName(name, cfg.Language)))

	// Remove from config and lock.
	delete(cfg.Packages, name)
	if hasLock {
		delete(lock.Packages, name)
	}
}

// outputName converts a package name to its output directory name,
// handling language-specific conventions (e.g. Python uses underscores).
func outputName(name, language string) string {
	if language == "python" {
		return strings.ReplaceAll(name, "-", "_")
	}
	return name
}

// --------------------------------------------------------------------
// State persistence
// --------------------------------------------------------------------

// persistState writes the updated config and lock file back to disk.
func persistState(ws *workspace.Workspace, cfg types.AifuncConfig, lock *types.LockFile, hasLock bool) error {
	if hasLock {
		if err := writeLock(ws, lock); err != nil {
			return err
		}
	}
	if err := writeConfig(ws.ConfigPath(), cfg); err != nil {
		return fmt.Errorf("failed to write aifunc.json: %w", err)
	}
	return nil
}

// writeLock persists the lock file, or removes it entirely when empty.
func writeLock(ws *workspace.Workspace, lock *types.LockFile) error {
	if len(lock.Packages) == 0 && len(lock.Engines) == 0 {
		os.Remove(ws.LockPath())
		return nil
	}
	if err := lockfile.Write(ws.LockPath(), *lock); err != nil {
		return fmt.Errorf("failed to write aifunc-lock.json: %w", err)
	}
	return nil
}

// --------------------------------------------------------------------
// Orphaned engine cleanup
// --------------------------------------------------------------------

// cleanOrphanedEngines removes engine entries from the lock that are no longer
// referenced by any remaining package, and cleans up their files on disk.
func cleanOrphanedEngines(lock *types.LockFile, ws *workspace.Workspace, cfg types.AifuncConfig) {
	usedVersions := collectUsedEngineVersions(lock)
	remainingPkgs := collectRemainingPackageNames(lock)

	for key, eng := range lock.Engines {
		if !usedVersions[key] {
			removeEngineFiles(ws, cfg, eng, key)
			delete(lock.Engines, key)
			continue
		}
		// Prune stale UsedBy references.
		lock.Engines[key] = pruneUsedBy(eng, remainingPkgs)
	}

	removeIfEmptyDir(filepath.Join(ws.EngineCachePath(), cfg.Language))
}

func collectUsedEngineVersions(lock *types.LockFile) map[string]bool {
	m := make(map[string]bool, len(lock.Packages))
	for _, pkg := range lock.Packages {
		m["v"+pkg.EngineVersion] = true
	}
	return m
}

func collectRemainingPackageNames(lock *types.LockFile) map[string]bool {
	m := make(map[string]bool, len(lock.Packages))
	for name := range lock.Packages {
		m[name] = true
	}
	return m
}

func pruneUsedBy(eng types.EngineLock, remaining map[string]bool) types.EngineLock {
	var filtered []string
	for _, name := range eng.UsedBy {
		if remaining[name] {
			filtered = append(filtered, name)
		}
	}
	eng.UsedBy = filtered
	return eng
}

// removeEngineFiles removes cached and generated engine files for the given key.
func removeEngineFiles(ws *workspace.Workspace, cfg types.AifuncConfig, eng types.EngineLock, key string) {
	removed := false
	for _, rel := range engineRelPaths(eng, key, cfg.Language) {
		if removeAllIfExists(filepath.Join(ws.EngineCachePath(), rel)) {
			removed = true
		}
		if removeAllIfExists(filepath.Join(ws.Root, cfg.GetOutputDir(), "_engine", rel)) {
			removed = true
		}
	}
	if removed {
		fmt.Fprintf(os.Stdout, "cleaning up orphaned engine %s\n", key)
	}
}

// engineRelPaths returns all candidate relative paths for an engine version,
// accounting for language-specific naming conventions.
func engineRelPaths(eng types.EngineLock, key, defaultLang string) []string {
	lang := eng.Language
	if lang == "" {
		lang = defaultLang
	}

	version := strings.TrimPrefix(key, "v")
	if version == "" {
		if eng.Path != "" {
			return []string{eng.Path}
		}
		return nil
	}

	paths := newUniqueStrings()
	paths.Add(eng.Path)

	if lang != "" {
		if lang == "python" {
			paths.Add(filepath.Join(lang, "v"+strings.ReplaceAll(version, ".", "_")))
		}
		paths.Add(filepath.Join(lang, "v"+version))
	}

	return paths.Slice()
}

// --------------------------------------------------------------------
// Directory cleanup helpers
// --------------------------------------------------------------------

func cleanEmptyDirs(ws *workspace.Workspace, cfg types.AifuncConfig) {
	outputDir := filepath.Join(ws.Root, cfg.GetOutputDir())

	// Clean _engine subdirectory.
	removeEmptyTree(filepath.Join(outputDir, "_engine"), markerFiles(cfg.Language, true))

	// Clean output directory itself.
	removeEmptyTree(outputDir, markerFiles(cfg.Language, false))

	// Clean cache directories.
	removeIfEmptyDir(ws.PackagesPath())
	removeEmptyTreeStrict(ws.EngineCachePath())
	removeIfEmptyDir(ws.CachePath())
}

// markerFiles returns the set of "marker" files that are allowed to exist in
// an otherwise-empty directory (and should be removed along with it).
func markerFiles(language string, engineDir bool) map[string]bool {
	m := make(map[string]bool)
	if language == "python" {
		m["__init__.py"] = true
		if !engineDir {
			m["py.typed"] = true
		}
	}
	return m
}

// removeIfEmptyDir removes a directory only if it is completely empty.
func removeIfEmptyDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		return
	}
	os.Remove(dir)
}

// removeEmptyTreeStrict recursively removes empty directories without
// considering any marker files.
func removeEmptyTreeStrict(dir string) {
	removeEmptyTree(dir, nil)
}

// removeEmptyTree recursively removes directories that contain only marker
// files (and no other content). Marker files are removed together with the dir.
func removeEmptyTree(dir string, markers map[string]bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	// Recurse into subdirectories first.
	for _, e := range entries {
		if e.IsDir() {
			removeEmptyTree(filepath.Join(dir, e.Name()), markers)
		}
	}

	// Re-read after children may have been removed.
	entries, err = os.ReadDir(dir)
	if err != nil {
		return
	}

	// Check if all remaining entries are removable markers.
	if !allMarkers(entries, markers) {
		return
	}

	// Remove marker files, then the directory itself.
	for _, e := range entries {
		os.Remove(filepath.Join(dir, e.Name()))
	}
	os.Remove(dir)
}

// allMarkers returns true if every entry is a file listed in markers.
func allMarkers(entries []os.DirEntry, markers map[string]bool) bool {
	for _, e := range entries {
		if e.IsDir() || !markers[e.Name()] {
			return false
		}
	}
	return true
}

// removeAllIfExists removes a path if it exists, returning true on success.
func removeAllIfExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return os.RemoveAll(path) == nil
}

// --------------------------------------------------------------------
// Package name resolution
// --------------------------------------------------------------------

// resolvePackageName attempts to find the package name in cfg.Packages when
// the user passes a source URL or shorthand instead of a plain name.
func resolvePackageName(input string, cfg types.AifuncConfig) string {
	src, err := source.Parse(input)
	if err != nil || src.Kind == source.KindLocal {
		return ""
	}

	normalised := normaliseSource(src, input)

	// Match by normalised source value.
	for name, val := range cfg.Packages {
		if val == normalised {
			return name
		}
		if valSrc, err := source.Parse(val); err == nil {
			if normaliseSource(valSrc, val) == normalised {
				return name
			}
		}
	}

	// Fall back: use the last segment of the sub-path as a candidate name.
	if src.SubPath != "" {
		candidate := filepath.Base(src.SubPath)
		if _, ok := cfg.Packages[candidate]; ok {
			return candidate
		}
	}

	return ""
}

// --------------------------------------------------------------------
// Utility: uniqueStrings
// --------------------------------------------------------------------

// uniqueStrings is a small helper to collect unique non-empty strings in order.
type uniqueStrings struct {
	seen  map[string]bool
	items []string
}

func newUniqueStrings() *uniqueStrings {
	return &uniqueStrings{seen: make(map[string]bool)}
}

func (u *uniqueStrings) Add(s string) {
	s = strings.TrimSpace(s)
	if s == "" || u.seen[s] {
		return
	}
	u.seen[s] = true
	u.items = append(u.items, s)
}

func (u *uniqueStrings) Slice() []string {
	return u.items
}
