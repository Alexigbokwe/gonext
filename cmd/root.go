package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var watchMode bool

var rootCmd = &cobra.Command{
	Use:   "gonext",
	Short: "GoNext CLI - Scaffolding and code generation for GoNext projects",
	Long:  `A CLI tool to scaffold and manage GoNext framework projects.`,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the GoNext project",
	Run: func(cmd *cobra.Command, args []string) {
		if watchMode {
			// Try to use 'air' for hot reloading
			if _, err := exec.LookPath("air"); err != nil {
				fmt.Println("Error: 'air' is required for watch mode but not installed. Install with 'go install github.com/cosmtrek/air@latest'.")
				return
			}
			fmt.Println("Starting in watch mode (hot reload)...")
			c := exec.Command("air")
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				fmt.Printf("Error running air: %v\n", err)
			}
			return
		}
		// Default: go run main.go
		fmt.Println("Starting GoNext project...")
		c := exec.Command("go", "run", "main.go")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		if err := c.Run(); err != nil {
			fmt.Printf("Error running project: %v\n", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	startCmd.Flags().BoolVar(&watchMode, "watch", false, "Enable watch mode (hot reload)")
	rootCmd.AddCommand(startCmd)
}
