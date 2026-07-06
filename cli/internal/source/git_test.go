// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package source

import "testing"

func TestRepoURL(t *testing.T) {
	tests := []struct {
		name string
		src  Source
		want string
	}{
		{
			name: "github",
			src: Source{
				Provider: ProviderGitHub,
				Owner:    "aifunc-dev",
				Repo:     "aifunc-packages",
			},
			want: "https://github.com/aifunc-dev/aifunc-packages.git",
		},
		{
			name: "gitee",
			src: Source{
				Provider: ProviderGitee,
				Owner:    "aifunc-dev",
				Repo:     "aifunc-packages",
			},
			want: "https://gitee.com/aifunc-dev/aifunc-packages.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RepoURL(tt.src)
			if err != nil {
				t.Fatalf("RepoURL() returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("RepoURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRepoURLError(t *testing.T) {
	if _, err := RepoURL(Source{Provider: "gitlab"}); err == nil {
		t.Fatal("RepoURL() expected error")
	}
}

func TestCloneArgs(t *testing.T) {
	got := cloneArgs(Source{Ref: "master"}, "https://example.com/repo.git", "tmp")
	want := []string{"clone", "--depth", "1", "--single-branch", "--branch", "master", "https://example.com/repo.git", "tmp"}
	if !equalStringSlices(got, want) {
		t.Fatalf("cloneArgs() = %#v, want %#v", got, want)
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
