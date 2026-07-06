// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package compiler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aifunc/cli/internal/fileutil"
	"aifunc/cli/internal/types"
)

// CompilePackage reads the package source at pkgDir and writes the compiled
// artifact to outputDir/<name>.aifunc.json. Returns the compiled artifact.
func CompilePackage(pkgDir, outputDir string) (*types.CompiledArtifact, error) {
	spec, err := loadPackageSpec(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("read package.json: %w", err)
	}

	api, err := loadAPISpec(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("read api.json: %w", err)
	}

	prompts, err := loadPrompts(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("read prompts: %w", err)
	}

	modelParams, err := loadModelParams(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("read model-params.json: %w", err)
	}

	modelRouting, err := loadModelRouting(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("read model-routing.json: %w", err)
	}

	if spec.EngineOptions != nil && spec.EngineOptions.InjectOutputSchema != nil {
		api.InjectOutputSchema = spec.EngineOptions.InjectOutputSchema
	}

	artifact := &types.CompiledArtifact{
		SchemaVersion:   "0.1.0",
		ArtifactVersion: "0.1.0",
		Package:         spec,
		API:             api,
		ModelParams:     modelParams,
		ModelRouting:    modelRouting,
		Prompts:         prompts,
		Metadata: types.ArtifactMeta{
			SourcePackageVersion: spec.Version,
			GeneratedAt:          time.Now().UTC().Format(time.RFC3339),
			ContentHash:          computeContentHash(spec, api, prompts),
		},
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	safeName := sanitizeName(spec.Name)
	outPath := filepath.Join(outputDir, safeName+".aifunc.json")
	data, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal artifact: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return nil, fmt.Errorf("write artifact: %w", err)
	}

	return artifact, nil
}

func loadPackageSpec(dir string) (types.PackageSpec, error) {
	data, err := fileutil.ReadJSON(filepath.Join(dir, "package.json"))
	if err != nil {
		return types.PackageSpec{}, err
	}
	var spec types.PackageSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return types.PackageSpec{}, err
	}
	return spec, nil
}

func loadAPISpec(dir string) (types.ApiSpec, error) {
	data, err := fileutil.ReadJSON(filepath.Join(dir, "api.json"))
	if err != nil {
		return types.ApiSpec{}, err
	}
	var api types.ApiSpec
	if err := json.Unmarshal(data, &api); err != nil {
		return types.ApiSpec{}, err
	}
	if strings.TrimSpace(api.Name) == "" {
		return types.ApiSpec{}, fmt.Errorf("api.json: name is required")
	}
	return api, nil
}

func loadPrompts(dir string) (types.PromptMap, error) {
	promptsDir := filepath.Join(dir, "prompts")
	entries, err := os.ReadDir(promptsDir)
	if err != nil {
		return nil, err
	}

	prompts := make(types.PromptMap)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		content, err := fileutil.ReadText(filepath.Join(promptsDir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read prompt %s: %w", e.Name(), err)
		}
		key := strings.TrimSuffix(e.Name(), ".md")
		prompts[key] = content
	}

	if _, ok := prompts["general"]; !ok {
		return nil, fmt.Errorf("prompts/general.md is required but not found")
	}
	return prompts, nil
}

func loadModelParams(dir string) (*types.ModelParams, error) {
	path := filepath.Join(dir, "model-params.json")
	data, err := fileutil.ReadJSON(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var mp types.ModelParams
	if err := json.Unmarshal(data, &mp); err != nil {
		return nil, err
	}
	return &mp, nil
}

func loadModelRouting(dir string) (*types.ModelRouting, error) {
	path := filepath.Join(dir, "model-routing.json")
	data, err := fileutil.ReadJSON(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var mr types.ModelRouting
	if err := json.Unmarshal(data, &mr); err != nil {
		return nil, err
	}
	return &mr, nil
}

func computeContentHash(spec types.PackageSpec, api types.ApiSpec, prompts types.PromptMap) string {
	h := sha256.New()
	data, _ := json.Marshal(struct {
		Spec    types.PackageSpec
		API     types.ApiSpec
		Prompts types.PromptMap
	}{spec, api, prompts})
	h.Write(data)
	return fmt.Sprintf("sha256:%x", h.Sum(nil))
}

// sanitizeName converts a scoped package name like @scope/name into scope__name
// so it is safe to use as a filename on all OSes.
func sanitizeName(name string) string {
	name = strings.TrimPrefix(name, "@")
	name = strings.ReplaceAll(name, "/", "__")
	return name
}
