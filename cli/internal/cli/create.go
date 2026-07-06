// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"aifunc/cli/internal/scaffold"

	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"new"},
		Short:   "Create a new package scaffold",
		Long:    "Create a template directory conforming to the AIFunc package spec in the current directory.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing required argument <name>\n\nUsage: aifn create <name>")
			}
			if len(args) > 1 {
				return fmt.Errorf("expected exactly 1 argument, got %d\n\nUsage: aifn create <name>", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(args[0])
		},
	}
}

func runCreate(name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	if err := scaffold.Create(cwd, name); err != nil {
		return err
	}

	pkgDir := filepath.Join(cwd, name)
	fmt.Fprintf(os.Stdout, "created package scaffold: %s\n", pkgDir)
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "generated files:")
	fmt.Fprintf(os.Stdout, "  %s/package.json       package metadata\n", name)
	fmt.Fprintf(os.Stdout, "  %s/api.json           API interface definition\n", name)
	// fmt.Fprintf(os.Stdout, "  %s/model-params.json  model parameter config\n", name)
	// fmt.Fprintf(os.Stdout, "  %s/mock.json          mock data\n", name)
	fmt.Fprintf(os.Stdout, "  %s/prompts/general.md prompt template\n", name)
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Next: edit the files above to define your AI function.")

	return nil
}
