// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package detect

import (
	"os"
	"path/filepath"
	"strings"
)

// Language represents a recommended project language.
type Language string

const (
	TypeScript Language = "typescript"
	Python     Language = "python"
)

// Result describes the outcome of scanning the project environment.
type Result struct {
	// Recommended is the language with the strongest evidence.
	// It is empty when no signals were found.
	Recommended Language
	// TypeScriptScore / PythonScore expose the weighted evidence counts,
	// mainly useful for diagnostics and tests.
	TypeScriptScore int
	PythonScore     int
	// Reasons lists the signal files/markers that drove the recommendation.
	Reasons []string
}

// tsMarkers are strong indicators of a TypeScript/JavaScript project.
var tsMarkers = []struct {
	name   string
	weight int
}{
	{"tsconfig.json", 5},
	{"package.json", 4},
	{"pnpm-lock.yaml", 3},
	{"yarn.lock", 3},
	{"package-lock.json", 3},
	{"bun.lockb", 3},
	{"node_modules", 2},
	{"deno.json", 4},
}

// pyMarkers are strong indicators of a Python project.
var pyMarkers = []struct {
	name   string
	weight int
}{
	{"pyproject.toml", 5},
	{"requirements.txt", 4},
	{"setup.py", 4},
	{"setup.cfg", 3},
	{"Pipfile", 4},
	{"poetry.lock", 3},
	{"uv.lock", 3},
	{".venv", 2},
	{"venv", 2},
}

// Scan inspects the given project root and recommends a language without
// committing to it. The recommendation is advisory only; callers decide
// whether and how to surface or honour it.
func Scan(root string) Result {
	res := Result{}

	for _, m := range tsMarkers {
		if exists(filepath.Join(root, m.name)) {
			res.TypeScriptScore += m.weight
			res.Reasons = append(res.Reasons, m.name)
		}
	}
	for _, m := range pyMarkers {
		if exists(filepath.Join(root, m.name)) {
			res.PythonScore += m.weight
			res.Reasons = append(res.Reasons, m.name)
		}
	}

	// Source-file presence is a weaker, fallback signal: only scan the top
	// level to keep this fast and avoid walking large trees.
	tsFiles, pyFiles := countTopLevelSources(root)
	if tsFiles > 0 {
		res.TypeScriptScore++
		res.Reasons = append(res.Reasons, "*.ts")
	}
	if pyFiles > 0 {
		res.PythonScore++
		res.Reasons = append(res.Reasons, "*.py")
	}

	switch {
	case res.TypeScriptScore > res.PythonScore:
		res.Recommended = TypeScript
	case res.PythonScore > res.TypeScriptScore:
		res.Recommended = Python
	}

	return res
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func countTopLevelSources(root string) (tsCount, pyCount int) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, 0
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.ToLower(e.Name())
		switch {
		case strings.HasSuffix(name, ".ts"), strings.HasSuffix(name, ".tsx"),
			strings.HasSuffix(name, ".mts"), strings.HasSuffix(name, ".cts"):
			tsCount++
		case strings.HasSuffix(name, ".py"), strings.HasSuffix(name, ".pyi"):
			pyCount++
		}
	}
	return tsCount, pyCount
}
