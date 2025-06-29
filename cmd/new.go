package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/spf13/cobra"
)

const starterRepo = "https://github.com/Alexigbokwe/Go_Next.git"
const oldModuleName = "goNext" // The module name used in the starter repo

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Scaffold a new GoNext project from the official starter template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		tempDir := projectName + "-tmp"

		// Check if git is installed
		if _, err := exec.LookPath("git"); err != nil {
			fmt.Println("Error: 'git' is required but not installed.")
			return
		}

		// Clone the starter repo into a temp directory
		cmdGit := exec.Command("git", "clone", starterRepo, tempDir)
		cmdGit.Stdout = os.Stdout
		cmdGit.Stderr = os.Stderr
		fmt.Printf("Cloning starter project from %s...\n", starterRepo)
		if err := cmdGit.Run(); err != nil {
			fmt.Printf("Error cloning repository: %v\n", err)
			return
		}

		// Remove .git directory from the cloned project
		gitDir := filepath.Join(tempDir, ".git")
		if err := os.RemoveAll(gitDir); err != nil {
			fmt.Printf("Warning: could not remove .git directory: %v\n", err)
		}

		// Rename the temp directory to the target project name
		if err := os.Rename(tempDir, projectName); err != nil {
			fmt.Printf("Error renaming project directory: %v\n", err)
			return
		}

		// Prompt for module path
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter module path (e.g., github.com/yourorg/%s) [default: %s]: ", projectName, projectName)
		modulePath, _ := reader.ReadString('\n')
		modulePath = strings.TrimSpace(modulePath)
		if modulePath == "" {
			modulePath = projectName
		}

		// Update go.mod in the new project directory
		goModPath := filepath.Join(projectName, "go.mod")
		if err := updateGoMod(goModPath, modulePath); err != nil {
			fmt.Printf("Error updating go.mod: %v\n", err)
		}

		// Update all import paths in .go files
		if err := updateImports(projectName, oldModuleName, modulePath); err != nil {
			fmt.Printf("Error updating import paths: %v\n", err)
		}

		fmt.Printf("New GoNext project '%s' created.\n", projectName)
		fmt.Println("Don't forget to run 'go mod tidy' in your new project!")
	},
}

// updateGoMod updates the module path in go.mod to newModule
func updateGoMod(goModPath, newModule string) error {
	input, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(input), "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "module ") {
		lines[0] = "module " + newModule
	}
	output := strings.Join(lines, "\n")
	return ioutil.WriteFile(goModPath, []byte(output), 0644)
}

func updateImports(rootDir, oldModule, newModule string) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content := strings.ReplaceAll(string(input), oldModule+"/", newModule+"/")
		// Also handle import aliasing: import goNext "goNext/app"
		content = strings.ReplaceAll(content, "\""+oldModule+"/", "\""+newModule+"/")
		if content != string(input) {
			if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
		}
		return nil
	})
}

func init() {
	rootCmd.AddCommand(newCmd)
}
