package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// EditText writes initial to a temporary file, opens it in the user's
// preferred terminal editor, and returns the edited contents.
func EditText(initial string) (string, error) {
	f, err := os.CreateTemp("", "redm-*.md")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	path := f.Name()
	defer os.Remove(path)

	if _, err := f.WriteString(initial); err != nil {
		f.Close()
		return "", fmt.Errorf("write temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close temp file: %w", err)
	}

	editor, err := resolveEditor()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("sh", "-c", editor+` "$1"`, "sh", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor failed: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read temp file: %w", err)
	}

	return strings.TrimRight(string(data), "\r\n"), nil
}

// resolveEditor returns the editor command to use, honoring $VISUAL and
// $EDITOR before falling back to common terminal editors.
func resolveEditor() (string, error) {
	for _, name := range []string{"VISUAL", "EDITOR"} {
		if v := strings.TrimSpace(os.Getenv(name)); v != "" {
			return v, nil
		}
	}
	for _, name := range []string{"vi", "nano", "vim"} {
		if _, err := exec.LookPath(name); err == nil {
			return name, nil
		}
	}
	return "", fmt.Errorf("no editor found: set $EDITOR")
}
