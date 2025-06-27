package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Helper to get the module name from go.mod
func getModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "myproject" // fallback, but should error in real use
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "module ") {
		return strings.TrimSpace(strings.TrimPrefix(lines[0], "module "))
	}
	return "myproject"
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from templates (coming soon)",
}

var gCmd = &cobra.Command{
	Use:   "g",
	Short: "Alias for generate",
}

// Helper to ensure module exists (creates if not)
func ensureModuleDirs(moduleName string) error {
	subdirs := []string{"controller", "repository", "route", "service"}
	moduleDir := filepath.Join("internal", moduleName)
	for _, sub := range subdirs {
		path := filepath.Join(moduleDir, sub)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("Error creating %s: %v", path, err)
		}
	}
	return nil
}

var controllerCmd = &cobra.Command{
	Use:   "controller [name] [in_module]",
	Short: "Generate a controller in a module (creates module if needed)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		module := args[1]
		titleName := strings.Title(name)
		moduleName := getModuleName()
		if err := ensureModuleDirs(module); err != nil {
			fmt.Println(err)
			return
		}
		controllerFile := filepath.Join("internal", module, "controller", fmt.Sprintf("%sController.go", name))
		if _, err := os.Stat(controllerFile); err == nil {
			fmt.Printf("Controller already exists: %s\n", controllerFile)
			return
		}
		content := fmt.Sprintf(`package controller

import (
	"github.com/gofiber/fiber/v2"
	"%s/internal/%s/service"
)

type %sController struct {
	Service *service.%sService `+"`inject:\"type\"`"+`
}

// Create%s handles creating a new %s
func (c *%sController) Create%s(ctx *fiber.Ctx) error {
	// TODO: Implement create logic
	return nil
}

// Get%s handles retrieving a %s by ID
func (c *%sController) Get%s(ctx *fiber.Ctx) error {
	// TODO: Implement get logic
	return nil
}

// Update%s handles updating a %s by ID
func (c *%sController) Update%s(ctx *fiber.Ctx) error {
	// TODO: Implement update logic
	return nil
}

// Delete%s handles deleting a %s by ID
func (c *%sController) Delete%s(ctx *fiber.Ctx) error {
	// TODO: Implement delete logic
	return nil
}
`,
			moduleName, module, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName)
		if err := os.WriteFile(controllerFile, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", controllerFile, err)
			return
		}
		fmt.Printf("Controller '%s' created in internal/%s/controller\n", name, module)
	},
}

var serviceCmd = &cobra.Command{
	Use:   "service [name] [in_module]",
	Short: "Generate a service in a module (creates module if needed)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		module := args[1]
		titleName := strings.Title(name)
		moduleName := getModuleName()
		if err := ensureModuleDirs(module); err != nil {
			fmt.Println(err)
			return
		}
		serviceFile := filepath.Join("internal", module, "service", fmt.Sprintf("%sService.go", name))
		if _, err := os.Stat(serviceFile); err == nil {
			fmt.Printf("Service already exists: %s\n", serviceFile)
			return
		}
		content := fmt.Sprintf(`package service

import (
	"%s/internal/%s/repository"
)

type %sService struct {
	Repository *repository.%sRepository `+"`inject:\"type\"`"+`
}

// Create%s creates a new %s
func (s *%sService) Create%s(data interface{}) error {
	// TODO: Implement create logic
	return nil
}

// Get%s retrieves a %s by ID
func (s *%sService) Get%s(id string) (interface{}, error) {
	// TODO: Implement get logic
	return nil, nil
}

// Update%s updates a %s by ID
func (s *%sService) Update%s(id string, data interface{}) error {
	// TODO: Implement update logic
	return nil
}

// Delete%s deletes a %s by ID
func (s *%sService) Delete%s(id string) error {
	// TODO: Implement delete logic
	return nil
}
`,
			moduleName, module, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName)
		if err := os.WriteFile(serviceFile, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", serviceFile, err)
			return
		}
		fmt.Printf("Service '%s' created in internal/%s/service\n", name, module)
	},
}

var repositoryCmd = &cobra.Command{
	Use:   "repository [name] [in_module]",
	Short: "Generate a repository in a module (creates module if needed)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		module := args[1]
		titleName := strings.Title(name)
		if err := ensureModuleDirs(module); err != nil {
			fmt.Println(err)
			return
		}
		repositoryFile := filepath.Join("internal", module, "repository", fmt.Sprintf("%sRepository.go", name))
		if _, err := os.Stat(repositoryFile); err == nil {
			fmt.Printf("Repository already exists: %s\n", repositoryFile)
			return
		}
		content := fmt.Sprintf(`package repository

type %sRepository struct{}

// Create%s persists a new %s
func (r *%sRepository) Create%s(data interface{}) error {
	// TODO: Implement create logic
	return nil
}

// Get%s retrieves a %s by ID
func (r *%sRepository) Get%s(id string) (interface{}, error) {
	// TODO: Implement get logic
	return nil, nil
}

// Update%s updates a %s by ID
func (r *%sRepository) Update%s(id string, data interface{}) error {
	// TODO: Implement update logic
	return nil
}

// Delete%s deletes a %s by ID
func (r *%sRepository) Delete%s(id string) error {
	// TODO: Implement delete logic
	return nil
}
`,
			titleName, titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName)
		if err := os.WriteFile(repositoryFile, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", repositoryFile, err)
			return
		}
		fmt.Printf("Repository '%s' created in internal/%s/repository\n", name, module)
	},
}

var moduleCmd = &cobra.Command{
	Use:   "module [name]",
	Short: "Generate a new module in internal/",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		titleName := strings.Title(name)
		moduleName := getModuleName()
		moduleDir := filepath.Join("internal", name)
		subdirs := []string{"controller", "repository", "route", "service"}
		for _, sub := range subdirs {
			path := filepath.Join(moduleDir, sub)
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Printf("Error creating %s: %v\n", path, err)
				return
			}
		}
		// Create module.go
		moduleGo := filepath.Join(moduleDir, "module.go")
		moduleGoContent := fmt.Sprintf(`package %s

import (
	"%s/app"
	"%s/internal/%s/controller"
	"%s/internal/%s/repository"
	"%s/internal/%s/route"
	"%s/internal/%s/service"

	"github.com/gofiber/fiber/v2"
)

type %sModule struct {
	Controller *controller.%sController
}

func New%sModule() *%sModule {
	return &%sModule{}
}

func (m *%sModule) Register(container *app.Container) {
	repo := &repository.%sRepository{}
	service := &service.%sService{}
	controller := &controller.%sController{}
	app.RegisterModuleComponents(container, repo, service, controller)
	m.Controller = controller
}

func (m *%sModule) MountRoutes(router fiber.Router) {
	group := router.Group("/%ss")
	route.Register%sRoutes(group, m.Controller)
}
`,
			name,
			moduleName, moduleName, name, moduleName, name, moduleName, name, moduleName, name,
			titleName, titleName,
			titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, name, titleName)
		if err := os.WriteFile(moduleGo, []byte(moduleGoContent), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", moduleGo, err)
			return
		}
		// Controller with CRUD and inject tag
		controllerFile := filepath.Join(moduleDir, "controller", fmt.Sprintf("%sController.go", name))
		controllerContent := fmt.Sprintf(`package controller

import (
	"github.com/gofiber/fiber/v2"
	"%s/internal/%s/service"
)

type %sController struct {
	Service *service.%sService `+"`inject:\"type\"`"+`
}

// Create%s handles creating a new %s
func (c *%sController) Create%s(ctx *fiber.Ctx) error {
	// TODO: Implement create logic
	return nil
}

// Get%s handles retrieving a %s by ID
func (c *%sController) Get%s(ctx *fiber.Ctx) error {
	// TODO: Implement get logic
	return nil
}

// Update%s handles updating a %s by ID
func (c *%sController) Update%s(ctx *fiber.Ctx) error {
	// TODO: Implement update logic
	return nil
}

// Delete%s handles deleting a %s by ID
func (c *%sController) Delete%s(ctx *fiber.Ctx) error {
	// TODO: Implement delete logic
	return nil
}
`,
			moduleName, name, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName)
		if err := os.WriteFile(controllerFile, []byte(controllerContent), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", controllerFile, err)
			return
		}
		// Service with CRUD and inject tag
		serviceFile := filepath.Join(moduleDir, "service", fmt.Sprintf("%sService.go", name))
		serviceContent := fmt.Sprintf(`package service

import (
	"%s/internal/%s/repository"
)

type %sService struct {
	Repository *repository.%sRepository `+"`inject:\"type\"`"+`
}

// Create%s creates a new %s
func (s *%sService) Create%s(data interface{}) error {
	// TODO: Implement create logic
	return nil
}

// Get%s retrieves a %s by ID
func (s *%sService) Get%s(id string) (interface{}, error) {
	// TODO: Implement get logic
	return nil, nil
}

// Update%s updates a %s by ID
func (s *%sService) Update%s(id string, data interface{}) error {
	// TODO: Implement update logic
	return nil
}

// Delete%s deletes a %s by ID
func (s *%sService) Delete%s(id string) error {
	// TODO: Implement delete logic
	return nil
}
`,
			moduleName, name, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName)
		if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", serviceFile, err)
			return
		}
		// Repository with CRUD
		repositoryFile := filepath.Join(moduleDir, "repository", fmt.Sprintf("%sRepository.go", name))
		repositoryContent := fmt.Sprintf(`package repository

type %sRepository struct{}

// Create%s persists a new %s
func (r *%sRepository) Create%s(data interface{}) error {
	// TODO: Implement create logic
	return nil
}

// Get%s retrieves a %s by ID
func (r *%sRepository) Get%s(id string) (interface{}, error) {
	// TODO: Implement get logic
	return nil, nil
}

// Update%s updates a %s by ID
func (r *%sRepository) Update%s(id string, data interface{}) error {
	// TODO: Implement update logic
	return nil
}

// Delete%s deletes a %s by ID
func (r *%sRepository) Delete%s(id string) error {
	// TODO: Implement delete logic
	return nil
}
`,
			titleName, titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName,
			titleName, titleName, titleName, titleName)
		if err := os.WriteFile(repositoryFile, []byte(repositoryContent), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", repositoryFile, err)
			return
		}
		// Route
		routeFile := filepath.Join(moduleDir, "route", fmt.Sprintf("%sRoute.go", name))
		routeContent := fmt.Sprintf(`package route

import (
	"github.com/gofiber/fiber/v2"
	"%s/internal/%s/controller"
)

func Register%sRoutes(route fiber.Router, ctrl *controller.%sController) {
	// TODO: Register routes for %s
}
`, moduleName, name, titleName, titleName, titleName)
		if err := os.WriteFile(routeFile, []byte(routeContent), 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", routeFile, err)
			return
		}
		fmt.Printf("Module '%s' created in internal/%s with boilerplate files and CRUD stubs.\n", name, name)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(gCmd)
	generateCmd.AddCommand(moduleCmd)
	gCmd.AddCommand(moduleCmd)
	generateCmd.AddCommand(controllerCmd)
	gCmd.AddCommand(controllerCmd)
	generateCmd.AddCommand(serviceCmd)
	gCmd.AddCommand(serviceCmd)
	generateCmd.AddCommand(repositoryCmd)
	gCmd.AddCommand(repositoryCmd)
}
