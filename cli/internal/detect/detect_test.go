// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(""), 0644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestScanRecommendsTypeScript(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "tsconfig.json")
	writeFile(t, dir, "package.json")

	res := Scan(dir)
	if res.Recommended != TypeScript {
		t.Fatalf("expected TypeScript, got %q (ts=%d py=%d)", res.Recommended, res.TypeScriptScore, res.PythonScore)
	}
}

func TestScanRecommendsPython(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "pyproject.toml")
	writeFile(t, dir, "main.py")

	res := Scan(dir)
	if res.Recommended != Python {
		t.Fatalf("expected Python, got %q (ts=%d py=%d)", res.Recommended, res.TypeScriptScore, res.PythonScore)
	}
}

func TestScanNoSignalsLeavesRecommendationEmpty(t *testing.T) {
	dir := t.TempDir()

	res := Scan(dir)
	if res.Recommended != "" {
		t.Fatalf("expected empty recommendation, got %q", res.Recommended)
	}
}

func TestScanStrongerSignalWins(t *testing.T) {
	dir := t.TempDir()
	// One weak python source file vs strong TS markers.
	writeFile(t, dir, "script.py")
	writeFile(t, dir, "tsconfig.json")

	res := Scan(dir)
	if res.Recommended != TypeScript {
		t.Fatalf("expected TypeScript to win, got %q (ts=%d py=%d)", res.Recommended, res.TypeScriptScore, res.PythonScore)
	}
}
