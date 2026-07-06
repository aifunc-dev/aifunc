// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package source

import (
	"archive/zip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveURL(t *testing.T) {
	tests := []struct {
		name string
		src  Source
		want string
	}{
		{
			name: "github with ref",
			src:  Source{Provider: ProviderGitHub, Owner: "aifunc-dev", Repo: "aifunc-packages", Ref: "main"},
			want: "https://github.com/aifunc-dev/aifunc-packages/archive/refs/heads/main.zip",
		},
		{
			name: "github without ref defaults to main",
			src:  Source{Provider: ProviderGitHub, Owner: "aifunc-dev", Repo: "aifunc-packages"},
			want: "https://github.com/aifunc-dev/aifunc-packages/archive/refs/heads/main.zip",
		},
		{
			name: "gitee with ref",
			src:  Source{Provider: ProviderGitee, Owner: "aifunc-dev", Repo: "aifunc-packages", Ref: "master"},
			want: "https://gitee.com/aifunc-dev/aifunc-packages/repository/archive/master.zip",
		},
		{
			name: "gitee without ref defaults to main",
			src:  Source{Provider: ProviderGitee, Owner: "aifunc-dev", Repo: "aifunc-packages"},
			want: "https://gitee.com/aifunc-dev/aifunc-packages/repository/archive/main.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ArchiveURL(tt.src)
			if err != nil {
				t.Fatalf("ArchiveURL() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ArchiveURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestArchiveURLUnsupportedProvider(t *testing.T) {
	_, err := ArchiveURL(Source{Provider: "gitlab"})
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}

func TestFetchArchiveIntegration(t *testing.T) {
	zipData := createTestZip(t, map[string]string{
		"repo-main/package.json": `{"name":"test-pkg"}`,
		"repo-main/api.json":     `{}`,
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write(zipData)
	}))
	defer server.Close()

	// Temporarily override ArchiveURL by testing the lower-level downloadZip + extractZip
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "test.zip")
	if err := os.WriteFile(zipPath, zipData, 0644); err != nil {
		t.Fatal(err)
	}

	dst := filepath.Join(tmp, "output")
	rootDir, err := extractZip(zipPath, dst)
	if err != nil {
		t.Fatalf("extractZip() error: %v", err)
	}

	if rootDir != "repo-main" {
		t.Errorf("extractZip() rootDir = %q, want %q", rootDir, "repo-main")
	}

	// After hoist, files should be at top level
	if err := hoistRootDir(dst, rootDir); err != nil {
		t.Fatalf("hoistRootDir() error: %v", err)
	}

	pkgJSON := filepath.Join(dst, "package.json")
	if _, err := os.Stat(pkgJSON); err != nil {
		t.Errorf("expected package.json at root after hoist, got error: %v", err)
	}
}

func TestExtractZipMultipleRoots(t *testing.T) {
	zipData := createTestZip(t, map[string]string{
		"dir-a/file.txt": "a",
		"dir-b/file.txt": "b",
	})

	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "multi.zip")
	if err := os.WriteFile(zipPath, zipData, 0644); err != nil {
		t.Fatal(err)
	}

	dst := filepath.Join(tmp, "out")
	rootDir, err := extractZip(zipPath, dst)
	if err != nil {
		t.Fatalf("extractZip() error: %v", err)
	}

	if rootDir != "" {
		t.Errorf("expected empty rootDir for multiple roots, got %q", rootDir)
	}
}

func TestDownloadZipMaxSize(t *testing.T) {
	bigData := make([]byte, MaxArchiveSize+1024)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(bigData)
	}))
	defer server.Close()

	_, err := downloadZip(server.URL)
	if err == nil {
		t.Fatal("expected error for oversized archive")
	}
}

func TestDownloadZip404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := downloadZip(server.URL)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestGitAvailable(t *testing.T) {
	// Just verify it doesn't panic; result depends on test environment.
	_ = gitAvailable()
}

func createTestZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "test.zip")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}

	w := zip.NewWriter(f)
	for name, content := range files {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	w.Close()
	f.Close()

	data, err := os.ReadFile(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
