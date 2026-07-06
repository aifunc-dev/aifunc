// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package manifest

import "aifunc/cli/internal/types"

func Validate(dir string) (types.PackageSpec, error) {
	return ValidatePackageDir(dir)
}
