// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package source

import (
	"fmt"
	"net/url"
	"strings"
)

type Kind string

const (
	KindLocal Kind = "local"
	KindGit   Kind = "git"
)

type Provider string

const (
	ProviderGitHub Provider = "github"
	ProviderGitee  Provider = "gitee"
)

type Source struct {
	Raw string

	Kind Kind
	Path string

	Provider Provider
	Owner    string
	Repo     string
	Ref      string
	SubPath  string
}

func Parse(raw string) (Source, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Source{}, fmt.Errorf("source cannot be empty")
	}
	if strings.HasPrefix(trimmed, "file:") {
		path := strings.TrimPrefix(trimmed, "file:")
		if path == "" {
			return Source{}, fmt.Errorf("file: source requires a path")
		}
		return Source{Raw: trimmed, Kind: KindLocal, Path: path}, nil
	}
	if strings.HasPrefix(trimmed, "github:") {
		return parseShorthand(trimmed, ProviderGitHub, "github:")
	}
	if strings.HasPrefix(trimmed, "gitee:") {
		return parseShorthand(trimmed, ProviderGitee, "gitee:")
	}
	if isHTTPURL(trimmed) {
		if src, err := parseTreeURL(trimmed); err == nil || isKnownGitHost(trimmed) {
			return src, err
		}
	}
	if strings.Contains(trimmed, "@") && !strings.ContainsAny(trimmed, `/\\`) {
		return Source{}, fmt.Errorf("name@version is not supported in v0.1; use a local path, github:owner/repo[/path], gitee:owner/repo[/path], or a supported tree URL")
	}
	return Source{Raw: trimmed, Kind: KindLocal, Path: trimmed}, nil
}

func parseShorthand(raw string, provider Provider, prefix string) (Source, error) {
	value := strings.TrimPrefix(raw, prefix)
	return parseGitParts(raw, provider, value, "")
}

func parseGitParts(raw string, provider Provider, value string, ref string) (Source, error) {
	parts := splitPath(value)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return Source{}, fmt.Errorf("invalid %s source %q; expected %s:owner/repo[/path]", provider, raw, provider)
	}
	return Source{
		Raw:      raw,
		Kind:     KindGit,
		Provider: provider,
		Owner:    parts[0],
		Repo:     parts[1],
		Ref:      ref,
		SubPath:  strings.Join(parts[2:], "/"),
	}, nil
}

func parseTreeURL(raw string) (Source, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return Source{}, fmt.Errorf("invalid source URL %q: %w", raw, err)
	}
	provider, ok := providerFromHost(parsed.Hostname())
	if !ok {
		return Source{}, fmt.Errorf("unsupported git host %q; supported hosts are github.com and gitee.com", parsed.Hostname())
	}

	parts := splitPath(parsed.EscapedPath())
	if len(parts) < 4 || parts[0] == "" || parts[1] == "" || parts[2] != "tree" || parts[3] == "" {
		return Source{}, fmt.Errorf("invalid %s tree URL %q; expected https://%s/owner/repo/tree/ref[/path]", provider, raw, parsed.Hostname())
	}

	subPathParts := parts[4:]
	return Source{
		Raw:      raw,
		Kind:     KindGit,
		Provider: provider,
		Owner:    parts[0],
		Repo:     parts[1],
		Ref:      parts[3],
		SubPath:  strings.Join(subPathParts, "/"),
	}, nil
}

func splitPath(value string) []string {
	trimmed := strings.Trim(value, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func isHTTPURL(value string) bool {
	return strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://")
}

func isKnownGitHost(value string) bool {
	parsed, err := url.Parse(value)
	if err != nil {
		return false
	}
	_, ok := providerFromHost(parsed.Hostname())
	return ok
}

func providerFromHost(host string) (Provider, bool) {
	switch strings.ToLower(host) {
	case "github.com", "www.github.com":
		return ProviderGitHub, true
	case "gitee.com", "www.gitee.com":
		return ProviderGitee, true
	default:
		return "", false
	}
}
