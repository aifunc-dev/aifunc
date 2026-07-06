// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package source

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyLocal(srcPath, dstPath string) error {
	absSrc, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}
	info, err := os.Stat(absSrc)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source %s is not a directory", srcPath)
	}
	if err := os.RemoveAll(dstPath); err != nil {
		return err
	}
	return copyDir(absSrc, dstPath)
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if shouldSkip(name) {
			continue
		}
		srcFile := filepath.Join(src, name)
		dstFile := filepath.Join(dst, name)
		entryInfo, err := entry.Info()
		if err != nil {
			return err
		}
		if entryInfo.IsDir() {
			if err := copyDir(srcFile, dstFile); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(srcFile, dstFile, entryInfo.Mode()); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func shouldSkip(name string) bool {
	return name == ".git" || name == "node_modules" || strings.HasSuffix(name, ".tmp")
}
