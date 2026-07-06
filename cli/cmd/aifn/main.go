// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"aifunc/cli/internal/cli"
)

func main() {
	if err := cli.Execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
