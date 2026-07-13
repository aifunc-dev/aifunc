// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"aifunc/cli/internal/fileutil"
	"aifunc/cli/internal/source"
	"aifunc/cli/internal/workspace"
)

type Manifest struct {
	ManifestVersion int                    `json:"manifestVersion"`
	Latest          map[string]string      `json:"latest"`
	Versions        map[string]VersionInfo `json:"versions"`
}

type VersionInfo struct {
	ReleaseDate          string                  `json:"releaseDate"`
	MinCliVersion        string                  `json:"minCliVersion"`
	PackageSchemaVersion string                  `json:"packageSchemaVersion"`
	Deprecated           bool                    `json:"deprecated"`
	Changelog            string                  `json:"changelog"`
	Languages            map[string]LanguageInfo `json:"languages"`
}

type LanguageInfo struct {
	Status    string  `json:"status"`
	Path      *string `json:"path"`
	Integrity *string `json:"integrity"`
}

func LoadManifest(repoRoot string) (*Manifest, error) {
	manifestPath := filepath.Join(repoRoot, "_engine", "manifest.json")
	data, err := fileutil.ReadJSON(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("_engine/manifest.json not found in repository")
		}
		return nil, fmt.Errorf("reading manifest.json: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parsing manifest.json: %w", err)
	}

	return &manifest, nil
}

// ResolveVersion finds the highest available engine version that satisfies engineRange
// (e.g. "^0.1.0", ">=0.1.0", "0.1.0") for the given language.
func ResolveVersion(manifest *Manifest, engineRange string, language string) (string, error) {
	if engineRange == "" {
		engineRange = "^0.1.0"
	}

	var resolvedVersion string
	for version, info := range manifest.Versions {
		ok, err := matchesRange(version, engineRange)
		if err != nil {
			return "", fmt.Errorf("invalid engine range %q: %w", engineRange, err)
		}
		if !ok {
			continue
		}

		langInfo, exists := info.Languages[language]
		if !exists || langInfo.Status == "unavailable" {
			continue
		}

		if resolvedVersion == "" || compareSemVerStr(version, resolvedVersion) > 0 {
			resolvedVersion = version
		}
	}

	if resolvedVersion == "" {
		return "", fmt.Errorf("no engine version found for %s satisfying %s", language, engineRange)
	}

	return resolvedVersion, nil
}

func Install(repoRoot string, version string, language string, ws *workspace.Workspace) error {
	manifest, err := LoadManifest(repoRoot)
	if err != nil {
		return err
	}

	versionInfo, ok := manifest.Versions[version]
	if !ok {
		return fmt.Errorf("version %s not found in manifest", version)
	}

	langInfo, ok := versionInfo.Languages[language]
	if !ok {
		return fmt.Errorf("version %s does not support language %s", version, language)
	}

	if langInfo.Status == "unavailable" {
		return fmt.Errorf("engine version %s for %s is unavailable", version, language)
	}

	if langInfo.Path == nil {
		return fmt.Errorf("engine version %s for %s is missing the path field", version, language)
	}

	srcPath := filepath.Join(repoRoot, "_engine", *langInfo.Path)
	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("engine source directory does not exist: %s", srcPath)
	}

	versionDir := "v" + version
	if language == "python" || language == "java" || language == "csharp" {
		versionDir = "v" + strings.ReplaceAll(version, ".", "_")
	}
	destPath := filepath.Join(ws.EngineCachePath(), language, versionDir)
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create engine cache directory: %w", err)
	}

	if err := source.CopyLocal(srcPath, destPath); err != nil {
		return fmt.Errorf("failed to copy engine files: %w", err)
	}

	return nil
}

// parseSemVer parses a "X.Y.Z" version string into its three numeric parts.
func parseSemVer(v string) ([3]int, error) {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return [3]int{}, fmt.Errorf("invalid version format: %q", v)
	}
	var nums [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return [3]int{}, fmt.Errorf("invalid version format: %q", v)
		}
		nums[i] = n
	}
	return nums, nil
}

// compareSemVer returns -1, 0, or 1 if a is less than, equal to, or greater than b.
func compareSemVer(a, b [3]int) int {
	for i := 0; i < 3; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func compareSemVerStr(a, b string) int {
	av, err := parseSemVer(a)
	if err != nil {
		return 0
	}
	bv, err := parseSemVer(b)
	if err != nil {
		return 0
	}
	return compareSemVer(av, bv)
}

// matchesRange reports whether candidate satisfies rangeStr.
// Supported operators: ^, >=, >, <=, <, and bare version (exact match).
//
// Caret semantics:
//   - ^X.Y.Z where X > 0  => >=X.Y.Z <(X+1).0.0
//   - ^0.Y.Z              => >=0.Y.Z  <0.(Y+1).0
func matchesRange(candidate, rangeStr string) (bool, error) {
	rangeStr = strings.TrimSpace(rangeStr)

	var op, versionStr string
	for _, prefix := range []string{">=", "<=", "^", ">", "<"} {
		if strings.HasPrefix(rangeStr, prefix) {
			op = prefix
			versionStr = strings.TrimSpace(rangeStr[len(prefix):])
			break
		}
	}
	if op == "" {
		op = "="
		versionStr = rangeStr
	}

	base, err := parseSemVer(versionStr)
	if err != nil {
		return false, fmt.Errorf("invalid version in engine range %q: %w", rangeStr, err)
	}

	v, err := parseSemVer(candidate)
	if err != nil {
		return false, fmt.Errorf("invalid candidate version %q: %w", candidate, err)
	}

	cmp := compareSemVer(v, base)

	switch op {
	case "^":
		if cmp < 0 {
			return false, nil
		}
		if base[0] > 0 {
			return v[0] == base[0], nil
		}
		return v[0] == 0 && v[1] == base[1], nil
	case ">=":
		return cmp >= 0, nil
	case ">":
		return cmp > 0, nil
	case "<=":
		return cmp <= 0, nil
	case "<":
		return cmp < 0, nil
	case "=":
		return cmp == 0, nil
	}
	return false, fmt.Errorf("unknown operator %q", op)
}

// CompareSemVerStr is the exported version of compareSemVerStr.
func CompareSemVerStr(a, b string) int {
	return compareSemVerStr(a, b)
}

// MatchesRange is the exported version of matchesRange.
func MatchesRange(candidate, rangeStr string) (bool, error) {
	return matchesRange(candidate, rangeStr)
}
