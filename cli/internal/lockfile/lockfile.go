// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package lockfile

import (
	"encoding/json"
	"os"

	"aifunc/cli/internal/fileutil"
	"aifunc/cli/internal/types"
)

func Read(path string) (types.LockFile, error) {
	data, err := fileutil.ReadJSON(path)
	if err != nil {
		return types.LockFile{}, err
	}
	var lf types.LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return types.LockFile{}, err
	}
	return lf, nil
}

func Write(path string, lf types.LockFile) error {
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}
