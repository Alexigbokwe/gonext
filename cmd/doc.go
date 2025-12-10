package cmd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed docs/documentation.md
var docFS embed.FS

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate and open the full GoNext framework documentation",
	Long:  `Generates a 'GoNext_Documentation.md' file in the current directory and attempts to open it in your default markdown viewer or editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		const fileName = "GoNext_Documentation.md"
		const docPath = "docs/documentation.md"

		// Read the embedded documentation
		content, err := docFS.ReadFile(docPath)
		if err != nil {
			fmt.Printf("Error checking internal documentation: %v\n", err)
			return
		}

		// Write to current directory
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}

		targetPath := filepath.Join(cwd, fileName)
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			fmt.Printf("Error writing documentation file: %v\n", err)
			return
		}

		fmt.Printf("✅ Documentation generated at: %s\n", targetPath)

		// detect and open in specific editors based on environment
		if openInEditor(targetPath) {
			return
		}

		// precise open command based on OS
		var openCmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			openCmd = exec.Command("open", targetPath)
		case "windows":
			openCmd = exec.Command("cmd", "/c", "start", targetPath)
		case "linux":
			openCmd = exec.Command("xdg-open", targetPath)
		default:
			fmt.Println("Could not detect OS to auto-open file. Please open it manually.")
			return
		}

		fmt.Println("Opening with default application...")
		if err := openCmd.Start(); err != nil {
			fmt.Printf("⚠️  Could not auto-open file: %v\n", err)
			fmt.Println("Please open 'GoNext_Documentation.md' in your favorite editor.")
		}
	},
}

// openInEditor attempts to open the file using the CLI tool of the current editor
func openInEditor(path string) bool {
	termProgram := os.Getenv("TERM_PROGRAM")
	terminalEmulator := os.Getenv("TERMINAL_EMULATOR")

	// Map of environment indicators to their CLI commands
	// Order matters: specific env vars should be checked first
	type editor struct {
		name  string
		cmd   []string // Changed to slice to support alternatives (e.g. goland, idea)
		check func() bool
	}

	editors := []editor{
		{
			name:  "Cursor",
			cmd:   []string{"cursor"},
			check: func() bool { return termProgram == "cursor" },
		},
		{
			name:  "VS Code",
			cmd:   []string{"code"},
			check: func() bool { return termProgram == "vscode" },
		},
		{
			name: "JetBrains IDE (GoLand/IntelliJ)",
			cmd:  []string{"goland", "idea"}, // Try goland first, then idea
			check: func() bool {
				return strings.Contains(terminalEmulator, "JetBrains") || os.Getenv("GO_NEXT_EDITOR") == "jetbrains"
			},
		},
		{
			name:  "Sublime Text",
			cmd:   []string{"subl"},
			check: func() bool { return os.Getenv("TERM_PROGRAM") == "Sublime" },
		},
		{
			name:  "Zed",
			cmd:   []string{"zed"},
			check: func() bool { return os.Getenv("TERM_PROGRAM") == "Zed" },
		},
	}

	for _, e := range editors {
		if e.check() {
			// Try all configured commands for this editor
			for _, command := range e.cmd {
				if _, err := exec.LookPath(command); err == nil {
					fmt.Printf("Detected %s terminal. Opening with '%s'...\n", e.name, command)
					cmd := exec.Command(command, path)
					if err := cmd.Run(); err == nil {
						return true
					}
					fmt.Printf("Failed to open with '%s', trying next option...\n", command)
				}
			}
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(docCmd)
}
