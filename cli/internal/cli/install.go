// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aifunc/cli/internal/downloader"
	"aifunc/cli/internal/engine"
	"aifunc/cli/internal/fileutil"
	"aifunc/cli/internal/lockfile"
	"aifunc/cli/internal/manifest"
	"aifunc/cli/internal/source"
	"aifunc/cli/internal/types"
	"aifunc/cli/internal/workspace"

	"github.com/spf13/cobra"
)

type engineInstallInfo struct {
	version     string
	engineRange string
	repoRoot    string
	source      string
	usedBy      []string
	cleanup     func()
}

func newInstallCommand() *cobra.Command {
	var lang string
	var outputDir string
	cmd := &cobra.Command{
		Use:     "install [source...]",
		Aliases: []string{"i"},
		Short:   "Install packages from aifunc.json or specified sources",
		Long: `Install packages in two modes:

  No arguments: reads aifunc.json and installs all declared packages.
  With arguments: downloads the specified packages and adds them to aifunc.json.

Supported source formats:
  gitee:owner/repo/path
  github:owner/repo/path
  https://gitee.com/owner/repo/tree/ref/path
  https://github.com/owner/repo/tree/ref/path`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(args, lang, outputDir)
		},
	}
	cmd.Flags().StringVarP(&lang, "lang", "l", "", "target language (overrides aifunc.json)")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (overrides aifunc.json)")
	return cmd
}

func runInstall(args []string, langOverride string, outputOverride string) error {
	ws, err := workspace.FromCurrentDir()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return runInstallFromConfig(ws, langOverride, outputOverride)
	}

	return runInstallSources(ws, args, langOverride, outputOverride)
}

func runInstallFromConfig(ws *workspace.Workspace, langOverride string, outputOverride string) error {
	if _, err := os.Stat(ws.ConfigPath()); errors.Is(err, os.ErrNotExist) {
		if langOverride != "" {
			return autoCreateConfig(ws, langOverride, outputOverride)
		}
		fmt.Fprintln(os.Stdout, "aifunc.json not found")
		fmt.Fprint(os.Stdout, "Initialize now? (Y/n) ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "y" || line == "" {
			if err := doInit(ws); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, "Add packages to aifunc.json and run aifunc install again.")
			return nil
		}
		fmt.Fprintln(os.Stdout, "Cancelled. Run aifunc init first.")
		return errors.New("aifunc.json not found")
	} else if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, "Reading aifunc.json ...")
	cfg, err := readConfig(ws.ConfigPath())
	if err != nil {
		return fmt.Errorf("aifunc.json parse error: %w", err)
	}

	if langOverride != "" {
		cfg.Language = langOverride
	}

	if outputOverride != "" {
		cfg.OutputDir = outputOverride
	}

	ws.SetInputDir(cfg.GetInputDir())

	if len(cfg.Packages) == 0 {
		fmt.Fprintln(os.Stdout, "No packages declared in aifunc.json.")
		fmt.Fprintln(os.Stdout, `Add entries to "packages", for example:`)
		fmt.Fprintln(os.Stdout, `  "packages": { "short-summary": "gitee:aifunc-dev/aifunc-packages/short-summary" }`)
		return nil
	}

	lock := readOrNewLock(ws)

	if err := ws.EnsureCache(); err != nil {
		return err
	}

	engines := make(map[string]*engineInstallInfo)

	for name, src := range cfg.Packages {
		if err := downloadOne(ws, cfg, name, src, lock, engines); err != nil {
			return err
		}
	}

	return finishInstall(ws, cfg, lock, engines, nil)
}

func runInstallSources(ws *workspace.Workspace, rawSources []string, langOverride string, outputOverride string) error {
	if _, err := os.Stat(ws.ConfigPath()); errors.Is(err, os.ErrNotExist) {
		if langOverride != "" {
			if err := autoCreateConfig(ws, langOverride, outputOverride); err != nil {
				return err
			}
		} else {
			fmt.Fprintln(os.Stdout, "aifunc.json not found, initializing first...")
			if err := doInit(ws); err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	}

	cfg, err := readConfig(ws.ConfigPath())
	if err != nil {
		return fmt.Errorf("aifunc.json parse error: %w", err)
	}

	originalLanguage := cfg.Language
	originalOutputDir := cfg.OutputDir
	if langOverride != "" {
		cfg.Language = langOverride
	}

	if outputOverride != "" {
		cfg.OutputDir = outputOverride
	}

	ws.SetInputDir(cfg.GetInputDir())

	lock := readOrNewLock(ws)

	if err := ws.EnsureCache(); err != nil {
		return err
	}

	engines := make(map[string]*engineInstallInfo)
	var installed []string

	for _, raw := range rawSources {
		src, err := source.Parse(raw)
		if err != nil {
			return fmt.Errorf("invalid package source %q: %w", raw, err)
		}

		if src.Kind == source.KindLocal {
			pkgName, absPath, err := downloader.ValidateLocal(src.Path)
			if err != nil {
				return fmt.Errorf("local package validation failed: %w", err)
			}
			fmt.Printf("linking %s ... ", raw)

			spec, err := manifest.LoadPackageSpec(absPath)
			if err != nil {
				fmt.Println("failed")
				return fmt.Errorf("reading package.json for %s: %w", pkgName, err)
			}

			engineVer, err := resolveEngineFromCache(ws, spec, cfg.Language)
			if err != nil {
				fmt.Println("failed")
				return fmt.Errorf("resolving engine version for %s: %w", pkgName, err)
			}

			relSource := localPathToFileSource(ws.Root, absPath)
			cfg.Packages[pkgName] = relSource
			lock.Packages[pkgName] = types.PackageLock{
				Source:        relSource,
				Version:       spec.Version,
				EngineVersion: engineVer,
				ResolvedAt:    time.Now().UTC().Format(time.RFC3339),
			}
			installed = append(installed, pkgName)
			fmt.Println("done")
			continue
		}

		fmt.Printf("downloading %s ... ", raw)
		result, err := downloader.Download(raw, ws)
		if err != nil {
			fmt.Println("failed")
			return fmt.Errorf("fetch failed: %w", err)
		}
		fmt.Println("done")

		name := result.Name
		normalised := normaliseSource(src, raw)
		cfg.Packages[name] = normalised

		pkgPath := filepath.Join(ws.PackagesPath(), name)
		spec, err := manifest.LoadPackageSpec(pkgPath)
		if err != nil {
			result.Cleanup()
			return fmt.Errorf("reading package.json for %s: %w", name, err)
		}

		engineVer, err := resolveEngineVersion(result.RepoRoot, spec, cfg.Language)
		if err != nil {
			result.Cleanup()
			return fmt.Errorf("resolving engine version for %s: %w", name, err)
		}

		lock.Packages[name] = types.PackageLock{
			Source:        normalised,
			Version:       spec.Version,
			EngineVersion: engineVer,
			ResolvedAt:    time.Now().UTC().Format(time.RFC3339),
		}

		engineRange := spec.Engine
		if engineRange == "" {
			engineRange = "^0.1.0"
		}

		if engines[engineVer] == nil {
			engines[engineVer] = &engineInstallInfo{
				version:     engineVer,
				engineRange: engineRange,
				repoRoot:    result.RepoRoot,
				source:      repoSource(src),
				usedBy:      []string{name},
				cleanup:     result.Cleanup,
			}
		} else {
			engines[engineVer].usedBy = append(engines[engineVer].usedBy, name)
			result.Cleanup()
		}

		installed = append(installed, name)
	}

	cfgToWrite := cfg
	cfgToWrite.Language = originalLanguage
	cfgToWrite.OutputDir = originalOutputDir
	cfgToWrite.Packages = make(map[string]string, len(cfg.Packages))
	for k, v := range cfg.Packages {
		cfgToWrite.Packages[k] = v
	}
	if err := writeConfig(ws.ConfigPath(), cfgToWrite); err != nil {
		return fmt.Errorf("failed to write aifunc.json: %w", err)
	}
	fmt.Printf("wrote %d package(s) to aifunc.json\n", len(installed))

	return finishInstall(ws, cfg, lock, engines, installed)
}

func downloadOne(ws *workspace.Workspace, cfg types.AifuncConfig, name, src string, lock types.LockFile, engines map[string]*engineInstallInfo) error {
	parsedSrc, parseErr := source.Parse(src)
	if parseErr != nil {
		return fmt.Errorf("invalid package source %q in aifunc.json: %w", src, parseErr)
	}

	if parsedSrc.Kind == source.KindLocal {
		return installLocalOne(ws, cfg, name, src, parsedSrc.Path, lock)
	}

	fmt.Printf("downloading %s ... ", name)
	result, err := downloader.Download(src, ws)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("fetch %s: %w", name, err)
	}
	fmt.Println("done")

	pkgPath := filepath.Join(ws.PackagesPath(), result.Name)
	spec, err := manifest.LoadPackageSpec(pkgPath)
	if err != nil {
		result.Cleanup()
		return fmt.Errorf("reading package.json for %s: %w", name, err)
	}

	engineVer, err := resolveEngineVersion(result.RepoRoot, spec, cfg.Language)
	if err != nil {
		result.Cleanup()
		return fmt.Errorf("resolving engine version for %s: %w", name, err)
	}

	lock.Packages[name] = types.PackageLock{
		Source:        src,
		Version:       spec.Version,
		EngineVersion: engineVer,
		ResolvedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	engineRange := spec.Engine
	if engineRange == "" {
		engineRange = "^0.1.0"
	}

	if engines[engineVer] == nil {
		engines[engineVer] = &engineInstallInfo{
			version:     engineVer,
			engineRange: engineRange,
			repoRoot:    result.RepoRoot,
			source:      repoSource(parsedSrc),
			usedBy:      []string{name},
			cleanup:     result.Cleanup,
		}
	} else {
		engines[engineVer].usedBy = append(engines[engineVer].usedBy, name)
		result.Cleanup()
	}

	return nil
}

// installLocalOne validates a local package and writes its lock entry.
// The package is never copied -- it is referenced by a relative path when on the
// same drive as the project root, or by an absolute path across drives (Windows).
func installLocalOne(ws *workspace.Workspace, cfg types.AifuncConfig, name, rawSrc, rawPath string, lock types.LockFile) error {
	fmt.Printf("linking %s ... ", rawSrc)

	pkgName, absPath, err := downloader.ValidateLocal(rawPath)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("local package validation failed: %w", err)
	}

	if name == "" {
		name = pkgName
	}

	spec, err := manifest.LoadPackageSpec(absPath)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("reading package.json for %s: %w", name, err)
	}

	engineVer, err := resolveEngineFromCache(ws, spec, cfg.Language)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("resolving engine version for %s: %w", name, err)
	}

	relSource := localPathToFileSource(ws.Root, absPath)
	lock.Packages[name] = types.PackageLock{
		Source:        relSource,
		Version:       spec.Version,
		EngineVersion: engineVer,
		ResolvedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	fmt.Println("done")
	return nil
}

func resolveEngineVersion(repoRoot string, spec types.PackageSpec, language string) (string, error) {
	manifest, err := engine.LoadManifest(repoRoot)
	if err != nil {
		return "", err
	}

	engineRange := spec.Engine
	if engineRange == "" {
		engineRange = "^0.1.0"
	}

	version, err := engine.ResolveVersion(manifest, engineRange, language)
	if err != nil {
		return "", err
	}

	return version, nil
}

// resolveEngineFromCache finds a cached engine version in .aifunc/_engine/ that satisfies the range.
// If no cached engine is found, it automatically fetches from the default remote engine repository.
func resolveEngineFromCache(ws *workspace.Workspace, spec types.PackageSpec, language string) (string, error) {
	engineRange := spec.Engine
	if engineRange == "" {
		engineRange = "^0.1.0"
	}

	engineDir := filepath.Join(ws.EngineCachePath(), language)
	entries, _ := os.ReadDir(engineDir)

	var resolved string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "v") {
			continue
		}
		version := strings.TrimPrefix(name, "v")
		if language == "python" {
			version = strings.ReplaceAll(version, "_", ".")
		}
		ok, err := engine.MatchesRange(version, engineRange)
		if err != nil || !ok {
			continue
		}
		if resolved == "" || engine.CompareSemVerStr(version, resolved) > 0 {
			resolved = version
		}
	}

	if resolved != "" {
		return resolved, nil
	}

	fmt.Println("no cached engine found, fetching from remote...")
	ver, err := fetchEngineFromRemote(ws, engineRange, language)
	if err != nil {
		return "", fmt.Errorf("auto-fetch engine failed: %w", err)
	}
	return ver, nil
}

const defaultEngineRepo = "github:aifunc-dev/aifunc-packages"

// fetchEngineFromRemote clones the default engine repo, resolves the engine version, and installs it.
func fetchEngineFromRemote(ws *workspace.Workspace, engineRange string, language string) (string, error) {
	src, err := source.Parse(defaultEngineRepo)
	if err != nil {
		return "", err
	}

	tmp, err := os.MkdirTemp("", "aifn-engine-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)

	repoDir := filepath.Join(tmp, "repo")
	if err := source.CloneRepo(src, repoDir); err != nil {
		return "", fmt.Errorf("cloning engine repository: %w", err)
	}

	manifest, err := engine.LoadManifest(repoDir)
	if err != nil {
		return "", err
	}

	version, err := engine.ResolveVersion(manifest, engineRange, language)
	if err != nil {
		return "", err
	}

	if err := engine.Install(repoDir, version, language, ws); err != nil {
		return "", fmt.Errorf("installing engine v%s: %w", version, err)
	}

	fmt.Printf("engine v%s ... done\n", version)
	return version, nil
}

func finishInstall(ws *workspace.Workspace, cfg types.AifuncConfig, lock types.LockFile, engines map[string]*engineInstallInfo, installed []string) error {
	for version, info := range engines {
		fmt.Printf("engine v%s ... ", version)

		if err := engine.Install(info.repoRoot, version, cfg.Language, ws); err != nil {
			info.cleanup()
			fmt.Println("failed")
			return fmt.Errorf("installing engine v%s: %w", version, err)
		}
		info.cleanup()

		lock.Engines["v"+version] = types.EngineLock{
			Source:       info.source,
			Language:     cfg.Language,
			Path:         engineCachePath(cfg.Language, version),
			ResolvedFrom: info.engineRange,
			ResolvedAt:   time.Now().UTC().Format(time.RFC3339),
			UsedBy:       mergeUsedBy(lock.Engines["v"+version].UsedBy, info.usedBy),
		}
		fmt.Println("done")
	}

	fmt.Printf("\nbuilding output to %s/ ...\n", cfg.GetOutputDir())

	if err := os.MkdirAll(cfg.GetOutputDir(), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	buildList := installed
	if len(buildList) == 0 {
		buildList = make([]string, 0, len(lock.Packages))
		for name := range lock.Packages {
			buildList = append(buildList, name)
		}
	}

	for _, name := range buildList {
		if err := buildOneWithEngine(ws, cfg, name, lock); err != nil {
			return fmt.Errorf("build %s: %w", name, err)
		}
	}

	if cfg.Language == "python" {
		ensurePythonRootMarkers(cfg.GetOutputDir())
	}

	if err := lockfile.Write(ws.LockPath(), lock); err != nil {
		return fmt.Errorf("failed to write aifunc-lock.json: %w", err)
	}

	if len(installed) > 0 {
		fmt.Printf("\ninstalled %d package(s).\n", len(installed))
		for _, name := range installed {
			if pkgLock, ok := lock.Packages[name]; ok {
				fmt.Printf("  %s (v%s, engine v%s)\n", name, pkgLock.Version, pkgLock.EngineVersion)
			}
		}
	} else {
		fmt.Printf("\ninstalled %d package(s).\n", len(lock.Packages))
		for name, pkgLock := range lock.Packages {
			fmt.Printf("  %s (v%s, engine v%s)\n", name, pkgLock.Version, pkgLock.EngineVersion)
		}
	}
	return nil
}

func readOrNewLock(ws *workspace.Workspace) types.LockFile {
	if _, err := os.Stat(ws.LockPath()); err == nil {
		if lf, err := lockfile.Read(ws.LockPath()); err == nil {
			return lf
		}
	}
	return types.LockFile{
		LockVersion: 1,
		Packages:    map[string]types.PackageLock{},
		Engines:     map[string]types.EngineLock{},
	}
}

// localPathToFileSource converts an absolute local path to a "file:" source string.
// If the package path is on the same drive/volume as the project root (basePath),
// it uses a relative path. On different drives (Windows only), it falls back to
// the absolute path since relative paths cannot cross drive boundaries.
func localPathToFileSource(basePath, absPath string) string {
	baseVol := filepath.VolumeName(basePath)
	pkgVol := filepath.VolumeName(absPath)

	if baseVol != "" && pkgVol != "" && !strings.EqualFold(baseVol, pkgVol) {
		return "file:" + filepath.ToSlash(absPath)
	}

	rel, err := filepath.Rel(basePath, absPath)
	if err != nil {
		return "file:" + filepath.ToSlash(absPath)
	}

	rel = filepath.ToSlash(rel)

	if !strings.HasPrefix(rel, ".") {
		rel = "./" + rel
	}
	return "file:" + rel
}

func normaliseSource(src source.Source, raw string) string {
	if src.Kind == source.KindLocal {
		if strings.HasPrefix(raw, "file:") {
			return raw
		}
		return "file:" + raw
	}
	if src.Kind != source.KindGit {
		return raw
	}
	path := src.Owner + "/" + src.Repo
	if src.SubPath != "" {
		path += "/" + src.SubPath
	}
	return string(src.Provider) + ":" + path
}

func repoSource(src source.Source) string {
	if src.Kind != source.KindGit {
		return src.Raw
	}
	return string(src.Provider) + ":" + src.Owner + "/" + src.Repo
}

func mergeUsedBy(existing []string, newNames []string) []string {
	seen := make(map[string]bool)
	for _, name := range existing {
		seen[name] = true
	}
	merged := append([]string{}, existing...)
	for _, name := range newNames {
		if !seen[name] {
			merged = append(merged, name)
		}
	}
	return merged
}

func engineCachePath(language string, version string) string {
	versionDir := "v" + version
	if language == "python" {
		versionDir = "v" + strings.ReplaceAll(version, ".", "_")
	}
	return language + "/" + versionDir
}

func readConfig(path string) (types.AifuncConfig, error) {
	data, err := fileutil.ReadJSON(path)
	if err != nil {
		return types.AifuncConfig{}, err
	}
	var cfg types.AifuncConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return types.AifuncConfig{}, err
	}
	return cfg, nil
}

func autoCreateConfig(ws *workspace.Workspace, lang string, outputDir string) error {
	if outputDir == "" {
		if strings.EqualFold(lang, "typescript") {
			outputDir = "src/aifunc"
		} else {
			outputDir = "aifunc"
		}
	}

	cfg := types.AifuncConfig{
		ConfigVersion: 1,
		Language:      strings.ToLower(lang),
		OutputDir:     outputDir,
		Packages:      map[string]string{},
	}

	if err := writeConfig(ws.ConfigPath(), cfg); err != nil {
		return fmt.Errorf("failed to write aifunc.json: %w", err)
	}
	fmt.Fprintf(os.Stdout, "created aifunc.json (language: %s, output: %s)\n", cfg.Language, outputDir)
	return nil
}

