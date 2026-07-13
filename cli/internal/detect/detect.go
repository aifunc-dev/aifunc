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
	Go         Language = "go"
	Java       Language = "java"
	CSharp     Language = "csharp"
)

// Result describes the outcome of scanning the project environment.
type Result struct {
	// Recommended is the language with the strongest evidence.
	// It is empty when no signals were found.
	Recommended Language
	// TypeScriptScore / PythonScore / GoScore / JavaScore / CSharpScore expose the weighted
	// evidence counts, mainly useful for diagnostics and tests.
	TypeScriptScore int
	PythonScore     int
	GoScore         int
	JavaScore       int
	CSharpScore     int
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

// goMarkers are strong indicators of a Go project.
var goMarkers = []struct {
	name   string
	weight int
}{
	{"go.mod", 5},
	{"go.sum", 3},
	{"go.work", 4},
}

// javaMarkers are strong indicators of a Java project.
var javaMarkers = []struct {
	name   string
	weight int
}{
	{"pom.xml", 5},
	{"build.gradle", 4},
	{"build.gradle.kts", 4},
	{"settings.gradle", 3},
	{"settings.gradle.kts", 3},
	{".mvn", 3},
	{"gradlew", 2},
	{"mvnw", 2},
}

// csharpMarkers are strong indicators of a C# / .NET project.
var csharpMarkers = []struct {
	name   string
	weight int
}{
	{"global.json", 4},
	{"nuget.config", 3},
	{"NuGet.Config", 3},
	{"Directory.Build.props", 4},
	{"Directory.Packages.props", 4},
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
	for _, m := range goMarkers {
		if exists(filepath.Join(root, m.name)) {
			res.GoScore += m.weight
			res.Reasons = append(res.Reasons, m.name)
		}
	}
	for _, m := range javaMarkers {
		if exists(filepath.Join(root, m.name)) {
			res.JavaScore += m.weight
			res.Reasons = append(res.Reasons, m.name)
		}
	}
	for _, m := range csharpMarkers {
		if exists(filepath.Join(root, m.name)) {
			res.CSharpScore += m.weight
			res.Reasons = append(res.Reasons, m.name)
		}
	}

	// Source-file presence is a weaker, fallback signal: only scan the top
	// level to keep this fast and avoid walking large trees.
	tsFiles, pyFiles, goFiles, javaFiles, csFiles := countTopLevelSources(root)
	if tsFiles > 0 {
		res.TypeScriptScore++
		res.Reasons = append(res.Reasons, "*.ts")
	}
	if pyFiles > 0 {
		res.PythonScore++
		res.Reasons = append(res.Reasons, "*.py")
	}
	if goFiles > 0 {
		res.GoScore++
		res.Reasons = append(res.Reasons, "*.go")
	}
	if javaFiles > 0 {
		res.JavaScore++
		res.Reasons = append(res.Reasons, "*.java")
	}
	if csFiles > 0 {
		res.CSharpScore++
		res.Reasons = append(res.Reasons, "*.cs/*.csproj/*.sln")
	}

	maxScore := res.TypeScriptScore
	res.Recommended = TypeScript
	if res.PythonScore > maxScore {
		maxScore = res.PythonScore
		res.Recommended = Python
	}
	if res.GoScore > maxScore {
		maxScore = res.GoScore
		res.Recommended = Go
	}
	if res.JavaScore > maxScore {
		maxScore = res.JavaScore
		res.Recommended = Java
	}
	if res.CSharpScore > maxScore {
		res.Recommended = CSharp
	}
	if res.TypeScriptScore == 0 && res.PythonScore == 0 && res.GoScore == 0 && res.JavaScore == 0 && res.CSharpScore == 0 {
		res.Recommended = ""
	}

	return res
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func countTopLevelSources(root string) (tsCount, pyCount, goCount, javaCount, csCount int) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, 0, 0, 0, 0
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
		case strings.HasSuffix(name, ".go"):
			goCount++
		case strings.HasSuffix(name, ".java"):
			javaCount++
		case strings.HasSuffix(name, ".cs"), strings.HasSuffix(name, ".csproj"),
			strings.HasSuffix(name, ".sln"):
			csCount++
		}
	}
	return tsCount, pyCount, goCount, javaCount, csCount
}
