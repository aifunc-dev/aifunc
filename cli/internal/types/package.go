// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/json"
	"fmt"
)

type PackageSpec struct {
	SchemaVersion string         `json:"schemaVersion,omitempty"`
	Type          string         `json:"type"`
	Name          string         `json:"name"`
	Version       string         `json:"version"`
	Description   string         `json:"description"`
	Author        *AuthorInfo    `json:"author,omitempty"`
	Engine        string         `json:"engine,omitempty"`
	EngineOptions *EngineOptions `json:"engineOptions,omitempty"`
}

func (p *PackageSpec) UnmarshalJSON(data []byte) error {
	type Alias PackageSpec
	raw := struct {
		Alias
		Engine json.RawMessage `json:"engine,omitempty"`
	}{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = PackageSpec(raw.Alias)

	if len(raw.Engine) == 0 {
		return nil
	}

	var str string
	if err := json.Unmarshal(raw.Engine, &str); err == nil {
		p.Engine = str
		return nil
	}

	var obj struct {
		MinVersion string `json:"minVersion"`
	}
	if err := json.Unmarshal(raw.Engine, &obj); err == nil && obj.MinVersion != "" {
		p.Engine = ">=" + obj.MinVersion
		return nil
	}

	return fmt.Errorf("engine field must be a semver range string (e.g. \"^0.1.0\") or legacy object {\"minVersion\": \"...\"}; got: %s", string(raw.Engine))
}

type AuthorInfo struct {
	Name string `json:"name,omitempty"`
}

type EngineOptions struct {
	InjectOutputSchema *bool `json:"injectOutputSchema,omitempty"`
}

type ApiSpec struct {
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	Input              map[string]any `json:"input"`
	Output             map[string]any `json:"output"`
	InjectOutputSchema *bool          `json:"injectOutputSchema,omitempty"`
}

type ModelParams struct {
	SchemaVersion string        `json:"schemaVersion"`
	Rules         []ParamPreset `json:"rules"`
}

type ParamPreset struct {
	Match          MatchRule      `json:"match"`
	Params         StandardParams `json:"params"`
	ProviderParams map[string]any `json:"providerParams,omitempty"`
}

type MatchRule struct {
	Model   string   `json:"model,omitempty"`
	Models  []string `json:"models,omitempty"`
	Pattern string   `json:"pattern,omitempty"`
}

type StandardParams struct {
	Temperature      *float64 `json:"temperature,omitempty"`
	TopP             *float64 `json:"topP,omitempty"`
	MaxTokens        *int     `json:"maxTokens,omitempty"`
	StructuredOutput *bool    `json:"structuredOutput,omitempty"`
}

type ModelRouting struct {
	SchemaVersion string   `json:"schemaVersion"`
	Default       string   `json:"default"`
	Fallback      []string `json:"fallback,omitempty"`
	Allowed       []string `json:"allowed,omitempty"`
	Denied        []string `json:"denied,omitempty"`
}

type PromptMap map[string]string
