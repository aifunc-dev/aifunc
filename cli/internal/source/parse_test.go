// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package source

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want Source
	}{
		{
			name: "gitee shorthand with subpath",
			raw:  "gitee:aifunc-dev/aifunc-packages/short-summary",
			want: Source{
				Raw:      "gitee:aifunc-dev/aifunc-packages/short-summary",
				Kind:     KindGit,
				Provider: ProviderGitee,
				Owner:    "aifunc-dev",
				Repo:     "aifunc-packages",
				SubPath:  "short-summary",
			},
		},
		{
			name: "github shorthand with subpath",
			raw:  "github:aifunc-dev/aifunc-packages/short-summary",
			want: Source{
				Raw:      "github:aifunc-dev/aifunc-packages/short-summary",
				Kind:     KindGit,
				Provider: ProviderGitHub,
				Owner:    "aifunc-dev",
				Repo:     "aifunc-packages",
				SubPath:  "short-summary",
			},
		},
		{
			name: "github tree URL",
			raw:  "https://github.com/aifunc-dev/aifunc-packages/tree/main/short-summary",
			want: Source{
				Raw:      "https://github.com/aifunc-dev/aifunc-packages/tree/main/short-summary",
				Kind:     KindGit,
				Provider: ProviderGitHub,
				Owner:    "aifunc-dev",
				Repo:     "aifunc-packages",
				Ref:      "main",
				SubPath:  "short-summary",
			},
		},
		{
			name: "gitee tree URL",
			raw:  "https://gitee.com/aifunc-dev/aifunc-packages/tree/master/short-summary",
			want: Source{
				Raw:      "https://gitee.com/aifunc-dev/aifunc-packages/tree/master/short-summary",
				Kind:     KindGit,
				Provider: ProviderGitee,
				Owner:    "aifunc-dev",
				Repo:     "aifunc-packages",
				Ref:      "master",
				SubPath:  "short-summary",
			},
		},
		{
			name: "local path",
			raw:  "./local-package",
			want: Source{
				Raw:  "./local-package",
				Kind: KindLocal,
				Path: "./local-package",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.raw)
			if err != nil {
				t.Fatalf("Parse() returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"",
		"gitee:bad",
		"github:bad",
		"https://github.com/aifunc-dev/aifunc-packages/blob/main/short-summary",
		"https://gitee.com/aifunc-dev/aifunc-packages/tree/",
		"name@version",
	}

	for _, raw := range tests {
		t.Run(raw, func(t *testing.T) {
			if _, err := Parse(raw); err == nil {
				t.Fatalf("Parse(%q) expected error", raw)
			}
		})
	}
}

func TestParseUnsupportedHTTPURLFallsBackToLocal(t *testing.T) {
	got, err := Parse("https://example.com/a/b/tree/master/pkg")
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}
	want := Source{
		Raw:  "https://example.com/a/b/tree/master/pkg",
		Kind: KindLocal,
		Path: "https://example.com/a/b/tree/master/pkg",
	}
	if got != want {
		t.Fatalf("Parse() = %#v, want %#v", got, want)
	}
}
