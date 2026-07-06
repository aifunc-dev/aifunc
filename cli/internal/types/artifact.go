// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package types

type CompiledArtifact struct {
	SchemaVersion   string        `json:"schemaVersion"`
	ArtifactVersion string        `json:"artifactVersion"`
	Package         PackageSpec   `json:"package"`
	API             ApiSpec       `json:"api"`
	ModelParams     *ModelParams  `json:"modelParams,omitempty"`
	ModelRouting    *ModelRouting `json:"modelRouting,omitempty"`
	Prompts         PromptMap     `json:"prompts"`
	Metadata        ArtifactMeta  `json:"metadata"`
}

type ArtifactMeta struct {
	SourcePackageVersion string `json:"sourcePackageVersion"`
	GeneratedAt          string `json:"generatedAt"`
	ContentHash          string `json:"contentHash"`
}
