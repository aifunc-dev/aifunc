// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"aifunc/cli/internal/builder/python"
	"aifunc/cli/internal/builder/typescript"
	"aifunc/cli/internal/compiler"
	"aifunc/cli/internal/fileutil"
	"aifunc/cli/internal/lockfile"
	"aifunc/cli/internal/source"
	"aifunc/cli/internal/types"
	"aifunc/cli/internal/workspace"

	"github.com/spf13/cobra"
)

func newBuildCommand() *cobra.Command {
	var lang string
	var outputDir string
	cmd := &cobra.Command{
		Use:   "build [packages...]",
		Short: "Build output from cached packages",
		Long:  `Compile packages in .aifunc/packages/ and write language interface code to outputDir.`,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBuild(args, lang, outputDir)
		},
	}
	cmd.Flags().StringVarP(&lang, "lang", "l", "", "target language (overrides aifunc.json)")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (overrides aifunc.json)")
	return cmd
}

func runBuild(packageNames []string, langOverride string, outputOverride string) error {
	ws, err := workspace.FromCurrentDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(ws.ConfigPath()); errors.Is(err, os.ErrNotExist) {
		if langOverride != "" {
			if err := autoCreateConfig(ws, langOverride, outputOverride); err != nil {
				return err
			}
		} else {
			fmt.Fprintln(os.Stdout, "aifunc.json not found, initializing first...")
			if err := doInit(ws); err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	}

	cfg, err := readConfig(ws.ConfigPath())
	if err != nil {
		return fmt.Errorf("reading aifunc.json: %w", err)
	}

	if langOverride != "" {
		cfg.Language = langOverride
	}

	if outputOverride != "" {
		cfg.OutputDir = outputOverride
	}

	switch cfg.Language {
	case "typescript", "python":
		// valid
	default:
		return fmt.Errorf("unsupported language %q; supported languages are typescript and python", cfg.Language)
	}

	ws.SetInputDir(cfg.GetInputDir())

	lock := readOrNewLock(ws)

	var toBuild []string
	if len(packageNames) == 0 {
		for name := range cfg.Packages {
			toBuild = append(toBuild, name)
		}
	} else {
		toBuild = packageNames
	}

	if len(toBuild) == 0 {
		fmt.Fprintln(os.Stdout, "No packages to build.")
		return nil
	}

	fmt.Printf("building %d package(s)...\n", len(toBuild))

	if err := os.MkdirAll(cfg.GetOutputDir(), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for _, name := range toBuild {
		if err := buildOneWithEngine(ws, cfg, name, lock); err != nil {
			return fmt.Errorf("build %s: %w", name, err)
		}
	}

	if cfg.Language == "python" {
		ensurePythonRootMarkers(cfg.GetOutputDir())
	}

	if err := lockfile.Write(ws.LockPath(), lock); err != nil {
		return fmt.Errorf("failed to write aifunc-lock.json: %w", err)
	}

	fmt.Printf("build complete.\n")
	return nil
}

func buildOneWithEngine(ws *workspace.Workspace, cfg types.AifuncConfig, name string, lock types.LockFile) error {
	pkgLock, ok := lock.Packages[name]
	if !ok {
		return fmt.Errorf("package %s not found in lock file, run aifn install first", name)
	}

	var pkgPath string
	if strings.HasPrefix(pkgLock.Source, "file:") {
		pkgPath = strings.TrimPrefix(pkgLock.Source, "file:")
		if !filepath.IsAbs(pkgPath) {
			pkgPath = filepath.Join(ws.Root, pkgPath)
		}
		if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
			return fmt.Errorf("local package %s path does not exist: %s\ncheck the path or run aifn install again", name, pkgPath)
		}
	} else {
		pkgPath = filepath.Join(ws.PackagesPath(), name)
		if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
			return fmt.Errorf("package %s is not downloaded, run aifn install first", name)
		}
	}

	fmt.Printf("  %s ... ", name)

	pkgOutputName := name
	if cfg.Language == "python" {
		pkgOutputName = strings.ReplaceAll(name, "-", "_")
	}
	pkgOutputDir := filepath.Join(cfg.GetOutputDir(), pkgOutputName)
	if err := os.MkdirAll(pkgOutputDir, 0755); err != nil {
		fmt.Println("failed")
		return fmt.Errorf("creating package output directory: %w", err)
	}

	artifact, err := compiler.CompilePackage(pkgPath, pkgOutputDir)
	if err != nil {
		fmt.Println("failed")
		return err
	}

	engineVersion := pkgLock.EngineVersion

	switch cfg.Language {
	case "typescript":
		hasMock := packageHasMock(pkgPath)
		if err := typescript.Generate(artifact, pkgOutputDir, hasMock, engineVersion); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("generating TypeScript code: %w", err)
		}
		if err := removeJSONArtifact(pkgOutputDir, artifact.Package.Name); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("removing JSON artifact: %w", err)
		}
		if err := writeTypeScriptMockModule(pkgPath, name, pkgOutputDir); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("generating TypeScript mock: %w", err)
		}
	case "python":
		hasMock := packageHasMock(pkgPath)
		if err := python.Generate(artifact, pkgOutputDir, hasMock, engineVersion); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("generating Python code: %w", err)
		}
		if err := removeJSONArtifact(pkgOutputDir, artifact.Package.Name); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("removing JSON artifact: %w", err)
		}
		if err := writePythonMockModule(pkgPath, name, pkgOutputDir); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("generating Python mock: %w", err)
		}
	default:
		fmt.Println("failed")
		return fmt.Errorf("unsupported language: %s", cfg.Language)
	}

	if cfg.Language != "typescript" && cfg.Language != "python" {
		if err := copyMockToPackageDir(pkgPath, name, pkgOutputDir); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("copying mock.json: %w", err)
		}
	}

	engineSrc := filepath.Join(ws.EngineCachePath(), cfg.Language, "v"+engineVersion)
	engineDstVersion := "v" + engineVersion
	if cfg.Language == "python" {
		engineSrc = filepath.Join(ws.EngineCachePath(), cfg.Language, "v"+strings.ReplaceAll(engineVersion, ".", "_"))
		engineDstVersion = "v" + strings.ReplaceAll(engineVersion, ".", "_")
	}
	engineDst := filepath.Join(cfg.GetOutputDir(), "_engine", cfg.Language, engineDstVersion)
	if _, err := os.Stat(engineSrc); err == nil {
		if err := os.RemoveAll(engineDst); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("cleaning engine output directory: %w", err)
		}
		if err := source.CopyLocal(engineSrc, engineDst); err != nil {
			fmt.Println("failed")
			return fmt.Errorf("copying engine to output directory: %w", err)
		}
		if cfg.Language == "python" {
			ensurePythonInitFiles(cfg.GetOutputDir(), cfg.Language, engineDstVersion)
		}
	}

	fmt.Println("done")
	return nil
}

func ensurePythonRootMarkers(outputDir string) {
	for _, name := range []string{"__init__.py", "py.typed"} {
		p := filepath.Join(outputDir, name)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			os.WriteFile(p, []byte(""), 0644)
		}
	}
}

func ensurePythonInitFiles(outputDir, language, engineDstVersion string) {
	dirs := []string{
		filepath.Join(outputDir, "_engine"),
		filepath.Join(outputDir, "_engine", language),
	}
	for _, dir := range dirs {
		initPath := filepath.Join(dir, "__init__.py")
		if _, err := os.Stat(initPath); os.IsNotExist(err) {
			os.WriteFile(initPath, []byte(""), 0644)
		}
	}
}

func packageHasMock(pkgPath string) bool {
	_, err := os.Stat(filepath.Join(pkgPath, "mock.json"))
	return err == nil
}

func removeJSONArtifact(pkgOutputDir, name string) error {
	artifactPath := filepath.Join(pkgOutputDir, sanitizeArtifactFileName(name)+".aifunc.json")
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(artifactPath)
}

func writeTypeScriptMockModule(pkgPath, name, pkgOutputDir string) error {
	mockSrc := filepath.Join(pkgPath, "mock.json")
	if _, err := os.Stat(mockSrc); os.IsNotExist(err) {
		return nil
	}

	data, err := fileutil.ReadJSON(mockSrc)
	if err != nil {
		return err
	}

	var mockData any
	if err := json.Unmarshal(data, &mockData); err != nil {
		return fmt.Errorf("parsing mock.json: %w", err)
	}

	encoded, err := json.MarshalIndent(mockData, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing mock data: %w", err)
	}

	var content strings.Builder
	content.WriteString("const mockData = ")
	content.Write(encoded)
	content.WriteString(";\n\n")
	content.WriteString("export default mockData;\n")

	mockDst := filepath.Join(pkgOutputDir, sanitizeArtifactFileName(name)+".mock.ts")
	return os.WriteFile(mockDst, []byte(content.String()), 0644)
}

func sanitizeArtifactFileName(name string) string {
	name = strings.TrimPrefix(name, "@")
	name = strings.ReplaceAll(name, "/", "__")
	return name
}

func writePythonMockModule(pkgPath, name, pkgOutputDir string) error {
	mockSrc := filepath.Join(pkgPath, "mock.json")
	if _, err := os.Stat(mockSrc); os.IsNotExist(err) {
		return nil
	}

	data, err := fileutil.ReadJSON(mockSrc)
	if err != nil {
		return err
	}

	var mockData any
	if err := json.Unmarshal(data, &mockData); err != nil {
		return fmt.Errorf("parsing mock.json: %w", err)
	}

	var content strings.Builder
	content.WriteString("# Generated by aifn - do not edit manually.\n\n")
	content.WriteString("mock_data = ")
	content.WriteString(python.ToPython(mockData, 0))
	content.WriteString("\n")

	safeName := strings.TrimPrefix(name, "@")
	safeName = strings.ReplaceAll(safeName, "/", "__")
	safeName = strings.ReplaceAll(safeName, "-", "_")
	mockDst := filepath.Join(pkgOutputDir, safeName+"_mock.py")
	return os.WriteFile(mockDst, []byte(content.String()), 0644)
}

func copyMockToPackageDir(pkgPath, name, pkgOutputDir string) error {
	mockSrc := filepath.Join(pkgPath, "mock.json")
	if _, err := os.Stat(mockSrc); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(mockSrc)
	if err != nil {
		return err
	}

	mockDst := filepath.Join(pkgOutputDir, name+".mock.json")
	return os.WriteFile(mockDst, data, 0644)
}
