// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"aifunc/cli/internal/manifest"

	"github.com/spf13/cobra"
)

func newValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <path>",
		Short: "Validate a package directory against the spec",
		Long:  "Check whether the specified directory is a valid AIFunc package: verifies required files, field completeness, and format correctness.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(args[0])
		},
	}
}

func runValidate(rawPath string) error {
	dir, err := filepath.Abs(rawPath)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dir)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dir)
	}

	fmt.Fprintf(os.Stdout, "validating %s ...\n\n", dir)

	spec, err := manifest.Validate(dir)
	if err != nil {
		fmt.Fprintf(os.Stdout, "validation failed: %s\n", err)
		return err
	}

	apiPath := filepath.Join(dir, "api.json")
	apiData, err := os.ReadFile(apiPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "reading api.json: %s\n", err)
		return err
	}
	var apiSpec map[string]any
	if err := json.Unmarshal(apiData, &apiSpec); err != nil {
		fmt.Fprintf(os.Stdout, "api.json parse error: %s\n", err)
		return err
	}

	fmt.Fprintln(os.Stdout, "validation passed.")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintf(os.Stdout, "  name:     %s\n", spec.Name)
	fmt.Fprintf(os.Stdout, "  version:  %s\n", spec.Version)
	fmt.Fprintf(os.Stdout, "  desc:     %s\n", spec.Description)
	fmt.Fprintf(os.Stdout, "  engine:   %s\n", spec.Engine)

	if name, ok := apiSpec["name"].(string); ok {
		fmt.Fprintf(os.Stdout, "  function: %s\n", name)
	}

	return nil
}
