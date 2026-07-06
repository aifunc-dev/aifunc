// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package source

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MaxArchiveSize    = 50 * 1024 * 1024 // 50MB
	archiveHTTPTimout = 60 * time.Second
)

// FetchArchive downloads the repository as a zip archive and extracts it to dstPath.
// This serves as a fallback when git is unavailable or git clone fails.
func FetchArchive(src Source, dstPath string) error {
	archiveURL, err := ArchiveURL(src)
	if err != nil {
		return err
	}

	zipPath, err := downloadZip(archiveURL)
	if err != nil {
		return fmt.Errorf("archive download failed: %w", err)
	}
	defer os.Remove(zipPath)

	repoDir, err := extractZip(zipPath, dstPath)
	if err != nil {
		return fmt.Errorf("archive extraction failed: %w", err)
	}

	// GitHub/Gitee zips contain a single root directory like "repo-branch/".
	// Move its contents up to dstPath directly.
	if repoDir != "" {
		return hoistRootDir(dstPath, repoDir)
	}
	return nil
}

// ArchiveURL returns the zip download URL for the given source.
func ArchiveURL(src Source) (string, error) {
	ref := src.Ref
	if ref == "" {
		ref = "main"
	}
	switch src.Provider {
	case ProviderGitHub:
		return fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", src.Owner, src.Repo, ref), nil
	case ProviderGitee:
		return fmt.Sprintf("https://gitee.com/%s/%s/repository/archive/%s.zip", src.Owner, src.Repo, ref), nil
	default:
		return "", fmt.Errorf("unsupported provider %q for archive download", src.Provider)
	}
}

func downloadZip(url string) (string, error) {
	client := &http.Client{Timeout: archiveHTTPTimout}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// ref might be "master" instead of "main"; caller can retry with alternate ref.
		return "", fmt.Errorf("archive not found (HTTP 404): %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status %d from %s", resp.StatusCode, url)
	}

	tmp, err := os.CreateTemp("", "aifn-archive-*.zip")
	if err != nil {
		return "", err
	}

	written, err := io.Copy(tmp, io.LimitReader(resp.Body, MaxArchiveSize+1))
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", err
	}
	tmp.Close()

	if written > MaxArchiveSize {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("archive exceeds maximum size (%dMB)", MaxArchiveSize/(1024*1024))
	}

	return tmp.Name(), nil
}

// extractZip extracts a zip file into dstPath and returns the single root directory
// name if the archive contains exactly one top-level directory (common for GitHub/Gitee).
func extractZip(zipPath, dstPath string) (rootDir string, err error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return "", err
	}

	roots := make(map[string]bool)

	for _, f := range r.File {
		name := filepath.FromSlash(f.Name)
		if isZipSlipPath(name, dstPath) {
			continue
		}

		parts := strings.SplitN(filepath.ToSlash(f.Name), "/", 2)
		if len(parts) > 0 && parts[0] != "" {
			roots[parts[0]] = true
		}

		target := filepath.Join(dstPath, name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, f.Mode())
			continue
		}

		if err := extractFile(f, target); err != nil {
			return "", err
		}
	}

	if len(roots) == 1 {
		for name := range roots {
			rootDir = name
		}
	}
	return rootDir, nil
}

func extractFile(f *zip.File, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

// hoistRootDir moves contents of dstPath/rootDir/* up to dstPath/.
func hoistRootDir(dstPath, rootDir string) error {
	nested := filepath.Join(dstPath, rootDir)
	entries, err := os.ReadDir(nested)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		src := filepath.Join(nested, entry.Name())
		dst := filepath.Join(dstPath, entry.Name())
		if err := os.Rename(src, dst); err != nil {
			return err
		}
	}
	return os.Remove(nested)
}

// isZipSlipPath checks whether extracting a file with the given name relative
// to dstPath would escape the destination directory. This is the canonical
// zip-slip defence: clean the joined path and verify it still has dstPath as
// a prefix.
func isZipSlipPath(name string, dstPath string) bool {
	cleaned := filepath.Join(dstPath, filepath.Clean("/"+name))
	// On Windows filepath.Clean normalises separators, so use os-aware prefix check.
	if !strings.HasPrefix(cleaned, filepath.Clean(dstPath)+string(os.PathSeparator)) {
		return true
	}
	return false
}
