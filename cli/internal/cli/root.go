// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/spf13/cobra"
)

var version = "dev"

func Execute(args []string) error {
	cmd := NewRootCommand()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "aifn",
		Short:         "AIFunc CLI",
		Long:          "aifn is the command-line tool for managing AIFunc package dependencies.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
	}
	cmd.AddCommand(
		newInitCommand(),
		newInstallCommand(),
		newBuildCommand(),
		newUninstallCommand(),
		newCreateCommand(),
		newValidateCommand(),
		newListCommand(),
	)
	return cmd
}
