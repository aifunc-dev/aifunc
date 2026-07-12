// Copyright 2026 GildenEye
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"aifunc/cli/internal/detect"
	"aifunc/cli/internal/fileutil"
	"aifunc/cli/internal/types"
	"aifunc/cli/internal/workspace"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize current directory as AIFunc project",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}
}

func runInit() error {
	ws, err := workspace.FromCurrentDir()
	if err != nil {
		return err
	}
	return doInit(ws)
}

func doInit(ws *workspace.Workspace) error {
	if _, err := os.Stat(ws.ConfigPath()); err == nil {
		fmt.Fprintln(os.Stdout, "aifunc.json already exists, skipping initialization.")
		fmt.Fprintln(os.Stdout, "To reinitialize, delete aifunc.json and run aifunc init again.")
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	options := []string{"TypeScript", "Python", "Go", "Java"}
	defaultIdx := recommendLanguageDefault(ws.Root, options)

	displayOptions := make([]string, len(options))
	copy(displayOptions, options)
	if defaultIdx >= 0 && defaultIdx < len(displayOptions) {
		res := detect.Scan(ws.Root)
		if res.Recommended != "" {
			displayOptions[defaultIdx] = displayOptions[defaultIdx] + " (detected)"
		}
	}

	language, err := promptSelect("Select project language", displayOptions, options, defaultIdx)
	if err != nil {
		return err
	}

	var defaultOutputDir string
	if strings.EqualFold(language, "typescript") {
		defaultOutputDir = "src/aifunc"
	} else {
		defaultOutputDir = "aifunc"
	}
	outputDir, err := promptText("Output directory", defaultOutputDir)
	if err != nil {
		return err
	}

	alias := ""
	if strings.EqualFold(language, "typescript") {
		alias, err = promptText("Path alias", "@aifunc")
		if err != nil {
			return err
		}
	}

	cfg := types.AifuncConfig{
		ConfigVersion: 1,
		Language:      strings.ToLower(language),
		OutputDir:     outputDir,
		Packages:      map[string]string{},
	}
	if alias != "" {
		cfg.Alias = alias
	}

	if err := writeConfig(ws.ConfigPath(), cfg); err != nil {
		return fmt.Errorf("failed to write aifunc.json: %w", err)
	}
	fmt.Fprintln(os.Stdout, "\naifunc.json created.")

	inputDir := cfg.GetInputDir()
	if err := ensureGitIgnore(ws.GitIgnorePath(), inputDir); err != nil {
		return fmt.Errorf("failed to update .gitignore: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Added %s/ to .gitignore\n", inputDir)
	fmt.Fprintln(os.Stdout, "Add packages to aifunc.json and run aifunc install.")
	return nil
}

// recommendLanguageDefault scans the project environment and returns the index
// of the recommended option to pre-select. It never forces a choice: when no
// signal is found it falls back to the first option, and the user can always
// pick another language at the prompt.
func recommendLanguageDefault(root string, options []string) int {
	res := detect.Scan(root)
	if res.Recommended == "" {
		return 0
	}

	fmt.Fprintf(os.Stdout, "Detected project environment, recommended language: %s (based on: %s)\n",
		languageDisplayName(res.Recommended), strings.Join(res.Reasons, ", "))

	for i, opt := range options {
		if strings.EqualFold(opt, string(res.Recommended)) {
			return i
		}
	}
	return 0
}

func languageDisplayName(lang detect.Language) string {
	switch lang {
	case detect.TypeScript:
		return "TypeScript"
	case detect.Python:
		return "Python"
	case detect.Go:
		return "Go"
	case detect.Java:
		return "Java"
	default:
		return string(lang)
	}
}

func writeConfig(path string, cfg types.AifuncConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}

func ensureGitIgnore(path string, inputDir string) error {
	entry := inputDir + "/"

	dataStr, err := fileutil.ReadText(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if strings.Contains(dataStr, entry) {
		return nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(dataStr) > 0 && !strings.HasSuffix(dataStr, "\n") {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(f, entry)
	return err
}

func promptSelect(question string, options []string, cleanValues []string, defaultIdx int) (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return promptSelectFallback(question, options, cleanValues, defaultIdx)
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return promptSelectFallback(question, options, cleanValues, defaultIdx)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	selected := defaultIdx
	if selected < 0 || selected >= len(options) {
		selected = 0
	}

	renderOptions := func() {
		fmt.Printf("\r\033[K? %s (use arrow keys to select, Enter to confirm)\n", question)
		for i, opt := range options {
			if i == selected {
				fmt.Printf("\r\033[K  \033[36m> %s\033[0m\n", opt)
			} else {
				fmt.Printf("\r\033[K    %s\n", opt)
			}
		}
		fmt.Printf("\033[%dA", len(options)+1)
	}

	renderOptions()

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			// EOF or read error: accept current selection gracefully
			fmt.Printf("\033[%dB\r\033[K", len(options)+1)
			fmt.Printf("\r\033[K? %s: \033[36m%s\033[0m\n", question, cleanValues[selected])
			return cleanValues[selected], nil
		}
		if n == 0 {
			continue
		}

		if n == 1 {
			switch buf[0] {
			case 13, 10:
				fmt.Printf("\033[%dB\r\033[K", len(options)+1)
				fmt.Printf("\r\033[K? %s: \033[36m%s\033[0m\n", question, cleanValues[selected])
				return cleanValues[selected], nil
			case 3:
				fmt.Printf("\033[%dB\r\033[K\n", len(options)+1)
				return "", fmt.Errorf("interrupted")
			case 'k', 'K':
				if selected > 0 {
					selected--
					renderOptions()
				}
			case 'j', 'J':
				if selected < len(options)-1 {
					selected++
					renderOptions()
				}
			}
		} else if n == 3 && buf[0] == 27 && buf[1] == 91 {
			switch buf[2] {
			case 65:
				if selected > 0 {
					selected--
					renderOptions()
				}
			case 66:
				if selected < len(options)-1 {
					selected++
					renderOptions()
				}
			}
		}
	}
}

func promptSelectFallback(question string, options []string, cleanValues []string, defaultIdx int) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n? %s:\n", question)
	for i, opt := range options {
		marker := " "
		if i == defaultIdx {
			marker = ">"
		}
		fmt.Printf("  %s %d) %s\n", marker, i+1, opt)
	}
	defaultVal := cleanValues[defaultIdx]
	fmt.Printf("  (default: %s) > ", defaultVal)

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal, nil
	}
	for i, opt := range cleanValues {
		if strings.EqualFold(line, opt) {
			return cleanValues[i], nil
		}
	}
	for i := range options {
		if line == fmt.Sprintf("%d", i+1) {
			return cleanValues[i], nil
		}
	}
	return defaultVal, nil
}

func promptText(question, defaultVal string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("? %s: (%s) ", question, defaultVal)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal, nil
	}
	return line, nil
}
