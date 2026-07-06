// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package source

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const gitCloneTimeout = 30 * time.Second

// CloneRepo clones the full repository into dstPath.
// If git is not installed or clone fails, it falls back to archive download.
func CloneRepo(src Source, dstPath string) error {
	if !gitAvailable() {
		fmt.Println("git not found, falling back to archive download...")
		return fetchViaArchive(src, dstPath)
	}

	repo, err := RepoURL(src)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), gitCloneTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", cloneArgs(src, repo, dstPath)...)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("git clone failed, falling back to archive download...\n")
		os.RemoveAll(dstPath)
		return fetchViaArchive(src, dstPath)
	} else {
		_ = output
	}
	return nil
}

func FetchGit(src Source, dstPath string) error {
	if !gitAvailable() {
		fmt.Println("git not found, falling back to archive download...")
		return fetchViaArchiveWithSubPath(src, dstPath)
	}

	tmp, err := os.MkdirTemp("", "aifn-git-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	repo, err := RepoURL(src)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), gitCloneTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", cloneArgs(src, repo, tmp)...)
	if output, err := cmd.CombinedOutput(); err != nil {
		_ = output
		fmt.Printf("git clone failed, falling back to archive download...\n")
		return fetchViaArchiveWithSubPath(src, dstPath)
	}

	packagePath := tmp
	if src.SubPath != "" {
		packagePath = filepath.Join(tmp, filepath.FromSlash(src.SubPath))
	}
	return CopyLocal(packagePath, dstPath)
}

// fetchViaArchive downloads the repo archive and extracts it to dstPath (full repo).
func fetchViaArchive(src Source, dstPath string) error {
	err := FetchArchive(src, dstPath)
	if err != nil && src.Ref == "" {
		// If ref was empty (defaulted to "main") and failed, retry with "master".
		alt := src
		alt.Ref = "master"
		if retryErr := FetchArchive(alt, dstPath); retryErr == nil {
			return nil
		}
	}
	return err
}

// fetchViaArchiveWithSubPath downloads the repo archive, extracts it, then copies SubPath to dstPath.
func fetchViaArchiveWithSubPath(src Source, dstPath string) error {
	tmp, err := os.MkdirTemp("", "aifn-archive-repo-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := fetchViaArchive(src, tmp); err != nil {
		return err
	}

	packagePath := tmp
	if src.SubPath != "" {
		packagePath = filepath.Join(tmp, filepath.FromSlash(src.SubPath))
	}
	return CopyLocal(packagePath, dstPath)
}

func gitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func RepoURL(src Source) (string, error) {
	switch src.Provider {
	case ProviderGitHub:
		return fmt.Sprintf("https://github.com/%s/%s.git", src.Owner, src.Repo), nil
	case ProviderGitee:
		return fmt.Sprintf("https://gitee.com/%s/%s.git", src.Owner, src.Repo), nil
	default:
		return "", fmt.Errorf("unsupported git provider %q", src.Provider)
	}
}

func cloneArgs(src Source, repoURL string, dstPath string) []string {
	args := []string{"clone", "--depth", "1", "--single-branch"}
	if src.Ref != "" {
		args = append(args, "--branch", src.Ref)
	}
	return append(args, repoURL, dstPath)
}
