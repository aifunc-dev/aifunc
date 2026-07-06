// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package types

type AifuncConfig struct {
	ConfigVersion int               `json:"configVersion"`
	Language      string            `json:"language"`
	InputDir      string            `json:"inputDir,omitempty"`
	OutputDir     string            `json:"outputDir,omitempty"`
	Alias         interface{}       `json:"alias,omitempty"`
	Packages      map[string]string `json:"packages"`
}

func (c AifuncConfig) GetInputDir() string {
	if c.InputDir != "" {
		return c.InputDir
	}
	return ".aifunc"
}

func (c AifuncConfig) GetOutputDir() string {
	if c.OutputDir != "" {
		return c.OutputDir
	}
	if c.Language == "python" {
		return "aifunc"
	}
	return "src/aifunc"
}

func (c AifuncConfig) GetAlias() string {
	if c.Alias == nil {
		if c.Language == "typescript" {
			return "@aifunc"
		}
		return ""
	}
	if v, ok := c.Alias.(string); ok {
		return v
	}
	if v, ok := c.Alias.(bool); ok && !v {
		return ""
	}
	return ""
}

type LockFile struct {
	LockVersion int                    `json:"lockVersion"`
	Packages    map[string]PackageLock `json:"packages"`
	Engines     map[string]EngineLock  `json:"engines"`
}

type PackageLock struct {
	Source        string `json:"source"`
	Version       string `json:"version"`
	Hash          string `json:"hash,omitempty"`
	EngineVersion string `json:"engineVersion"`
	ResolvedAt    string `json:"resolvedAt"`
}

type EngineLock struct {
	Source       string   `json:"source"`
	Language     string   `json:"language"`
	Path         string   `json:"path"`
	Integrity    string   `json:"integrity,omitempty"`
	ResolvedFrom string   `json:"resolvedFrom"`
	ResolvedAt   string   `json:"resolvedAt"`
	UsedBy       []string `json:"usedBy"`
}
