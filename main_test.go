package main

import (
	"testing"

	"github.com/spf13/cobra"
)

// buildCommandHierarchy creates a cobra command hierarchy from path names.
// Returns the leaf command.
func buildCommandHierarchy(path []string) *cobra.Command {
	var leaf *cobra.Command
	for i, name := range path {
		cmd := &cobra.Command{Use: name}
		if i == 0 {
			leaf = cmd
		} else {
			leaf.AddCommand(cmd)
			leaf = cmd
		}
	}
	return leaf
}

func TestIsCommandOrParent(t *testing.T) {
	tests := []struct {
		name     string
		cmdPath  []string
		names    []string
		expected bool
	}{
		{
			name:     "direct completion command",
			cmdPath:  []string{"talm", "completion"},
			names:    []string{"init", "completion"},
			expected: true,
		},
		{
			name:     "completion bash subcommand",
			cmdPath:  []string{"talm", "completion", "bash"},
			names:    []string{"init", "completion"},
			expected: true,
		},
		{
			name:     "init command",
			cmdPath:  []string{"talm", "init"},
			names:    []string{"init", "completion"},
			expected: true,
		},
		{
			name:     "apply command should not match",
			cmdPath:  []string{"talm", "apply"},
			names:    []string{"init", "completion"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaf := buildCommandHierarchy(tt.cmdPath)
			result := isCommandOrParent(leaf, tt.names...)
			if result != tt.expected {
				t.Errorf("isCommandOrParent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSkipConfigCommands(t *testing.T) {
	tests := []struct {
		name     string
		cmdPath  []string
		expected bool // true = should skip config loading
	}{
		{
			name:     "completion command",
			cmdPath:  []string{"talm", "completion"},
			expected: true,
		},
		{
			name:     "completion bash",
			cmdPath:  []string{"talm", "completion", "bash"},
			expected: true,
		},
		{
			name:     "completion zsh",
			cmdPath:  []string{"talm", "completion", "zsh"},
			expected: true,
		},
		{
			name:     "__complete (cobra internal for shell completion)",
			cmdPath:  []string{"talm", "__complete"},
			expected: true,
		},
		{
			name:     "init command",
			cmdPath:  []string{"talm", "init"},
			expected: true,
		},
		{
			name:     "apply command should load config",
			cmdPath:  []string{"talm", "apply"},
			expected: false,
		},
		{
			name:     "template command should load config",
			cmdPath:  []string{"talm", "template"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaf := buildCommandHierarchy(tt.cmdPath)
			// This uses the actual skipConfigCommands from main.go
			result := isCommandOrParent(leaf, skipConfigCommands...)
			if result != tt.expected {
				t.Errorf("skipConfigCommands check = %v, want %v (skipConfigCommands = %v)",
					result, tt.expected, skipConfigCommands)
			}
		})
	}
}
