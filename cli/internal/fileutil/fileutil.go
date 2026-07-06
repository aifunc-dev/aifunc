// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package fileutil

import (
	"bytes"
	"os"
	"strings"
)

var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// StripBOM removes a leading UTF-8 BOM from data if present.
func StripBOM(data []byte) []byte {
	return bytes.TrimPrefix(data, utf8BOM)
}

// NormalizeNewlines replaces \r\n with \n for cross-platform consistency.
func NormalizeNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// ReadJSON reads a file and strips UTF-8 BOM, returning bytes ready for json.Unmarshal.
func ReadJSON(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return StripBOM(data), nil
}

// ReadText reads a file, strips UTF-8 BOM, and normalizes line endings.
func ReadText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return NormalizeNewlines(string(StripBOM(data))), nil
}
