// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"aifunc/cli/internal/fileutil"
	"aifunc/cli/internal/types"
)

var packageNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

func LoadPackageSpec(dir string) (types.PackageSpec, error) {
	data, err := fileutil.ReadJSON(filepath.Join(dir, "package.json"))
	if err != nil {
		return types.PackageSpec{}, err
	}
	var spec types.PackageSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return types.PackageSpec{}, err
	}
	return spec, nil
}

func ValidatePackageDir(dir string) (types.PackageSpec, error) {
	spec, err := LoadPackageSpec(dir)
	if err != nil {
		return types.PackageSpec{}, fmt.Errorf("load package.json: %w", err)
	}
	if strings.TrimSpace(spec.Name) == "" {
		return types.PackageSpec{}, fmt.Errorf("package name is required")
	}
	if !packageNamePattern.MatchString(spec.Name) {
		return types.PackageSpec{}, fmt.Errorf("invalid package name %q", spec.Name)
	}
	if strings.TrimSpace(spec.Version) == "" {
		return types.PackageSpec{}, fmt.Errorf("version is required")
	}
	if strings.TrimSpace(spec.Engine) == "" {
		return types.PackageSpec{}, fmt.Errorf("engine is required (e.g. \"^0.1.0\")")
	}
	if err := requireFile(dir, "api.json"); err != nil {
		return types.PackageSpec{}, err
	}
	if err := requireDir(dir, "prompts"); err != nil {
		return types.PackageSpec{}, err
	}
	return spec, nil
}

func requireFile(root, rel string) error {
	info, err := os.Stat(filepath.Join(root, rel))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("missing required file %s", rel)
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s must be a file", rel)
	}
	return nil
}

func requireDir(root, rel string) error {
	info, err := os.Stat(filepath.Join(root, rel))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("missing required directory %s", rel)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s must be a directory", rel)
	}
	return nil
}
