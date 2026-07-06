// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package scaffold

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var validNameRe = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

func Create(targetDir, name string) error {
	if !validNameRe.MatchString(name) {
		return fmt.Errorf("invalid package name %q: must be kebab-case (e.g. my-package)", name)
	}

	pkgDir := filepath.Join(targetDir, name)
	if _, err := os.Stat(pkgDir); err == nil {
		return fmt.Errorf("directory %q already exists", pkgDir)
	}

	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	promptsDir := filepath.Join(pkgDir, "prompts")
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		return fmt.Errorf("creating prompts directory: %w", err)
	}

	files := map[string][]byte{
		"package.json":       generatePackageJSON(name),
		"api.json":           generateAPIJSON(name),
		"prompts/general.md": generatePrompt(),
		// "model-params.json":  generateModelParams(),
		// "mock.json":          generateMock(),
	}

	for filename, content := range files {
		path := filepath.Join(pkgDir, filename)
		if err := os.WriteFile(path, content, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}

	return nil
}

func generatePackageJSON(name string) []byte {
	pkg := map[string]any{
		"type":        "standalone",
		"name":        name,
		"version":     "1.0.0",
		"description": "TODO: fill in package description",
		"engine":      "^0.1.0",
		"author": map[string]string{
			"name": "TODO: fill in author name",
		},
	}
	data, _ := json.MarshalIndent(pkg, "", "  ")
	return append(data, '\n')
}

func generateAPIJSON(name string) []byte {
	api := map[string]any{
		"version":     "1.0.0",
		"name":        name,
		"description": "TODO: fill in API description",
		"input": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"text": map[string]any{
					"type":        "string",
					"description": "input text",
					"minLength":   1,
				},
			},
			"required":             []string{"text"},
			"additionalProperties": false,
		},
		"output": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"result": map[string]any{
					"type":        "string",
					"description": "processed result",
				},
			},
			"required":             []string{"result"},
			"additionalProperties": false,
		},
	}
	data, _ := json.MarshalIndent(api, "", "  ")
	return append(data, '\n')
}

func generateModelParams() []byte {
	params := map[string]any{
		"schemaVersion": "0.1.0",
		"rules": []map[string]any{
			{
				"match": map[string]string{
					"pattern": ".*",
				},
				"params": map[string]any{
					"temperature": 0.7,
					"maxTokens":   1024,
				},
			},
		},
	}
	data, _ := json.MarshalIndent(params, "", "  ")
	return append(data, '\n')
}

func generateMock() []byte {
	mock := map[string]any{
		"version": "1.0.0",
		"delay": map[string]int{
			"minMs": 30,
			"maxMs": 100,
		},
		"random": map[string]any{
			"enabled": false,
			"seed":    "default",
		},
		"cases": []map[string]any{
			{
				"id":          "basic",
				"description": "basic example",
				"output": map[string]string{
					"result": "This is a sample output.",
				},
			},
		},
	}
	data, _ := json.MarshalIndent(mock, "", "  ")
	return append(data, '\n')
}

func generatePrompt() []byte {
	return []byte(`# System

You are a strict JSON function. Return only JSON that conforms to the output schema. Do not include Markdown or any additional explanation.

# User

Input:
{{text}}
`)
}
