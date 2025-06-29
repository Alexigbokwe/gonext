# GoNext CLI

A CLI tool to scaffold and manage GoNext framework projects.

## Features

- Project scaffolding (`new` command)
- Code generation (`generate` command)
- Extensible with templates
- Project start and hot-reload support

## Installation

Install the CLI globally with:

```sh
go install github.com/Alexigbokwe/gonext@main
```

Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your PATH.

## Usage

### Create a Project

```sh
gonext new <project_name>
```

### Start the Project

Start your GoNext project:

```sh
gonext start
```

Start your project in watch mode (hot reload):

```sh
gonext start --watch
```

> **Note:** Watch mode requires [`air`](https://github.com/cosmtrek/air). Install it with:
>
> ```sh
> go install github.com/cosmtrek/air@latest
> ```

### Generate Modules and Components

- Generate a new module:

  ```sh
  gonext generate module <name>
  # or
  gonext g module <name>
  ```

- Generate a controller, service, or repository in a module (creates the module if it doesn't exist):
  ```sh
  gonext generate controller <name> <in_module>
  gonext generate service <name> <in_module>
  gonext generate repository <name> <in_module>
  # or use the 'g' alias
  gonext g controller <name> <in_module>
  gonext g service <name> <in_module>
  gonext g repository <name> <in_module>
  ```

## Code Generation

### Modules

- `gonext generate module <name>` or `gonext g module <name>`
  - Scaffolds a new module with controller, service, repository, and route boilerplate.

### Individual Components

- `gonext generate controller <name> <in_module>` or `gonext g controller <name> <in_module>`
- `gonext generate service <name> <in_module>` or `gonext g service <name> <in_module>`
- `gonext generate repository <name> <in_module>` or `gonext g repository <name> <in_module>`

### DTOs

- `gonext generate dto <name> <in_module>` or `gonext g dto <name> <in_module>`

  - Generates a DTO struct in `internal/<in_module>/dto/<name>DTO.go` with sample validation tags.
  - **Example:**

    ```sh
    gonext generate dto CreateUser account
    ```

    Output:

    ```go
    package dto

    type CreateUserDTO struct {
        Username string `json:"username" validate:"required,min=3,max=20"`
        FullName string `json:"full_name" validate:"required,min=3,max=50"`
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required,min=8"`
    }
    ```

### Middleware

- `gonext generate middleware <name> <in_module>` or `gonext g middleware <name> <in_module>`

  - Generates a sample Fiber middleware in `internal/<in_module>/middleware/<name>Middleware.go`.
  - **Example:**

    ```sh
    gonext generate middleware auth account
    ```

    Output:

    ```go
    package middleware

    import (
        "github.com/gofiber/fiber/v2"
    )

    // AuthMiddleware is a sample Fiber middleware
    func AuthMiddleware() fiber.Handler {
        return func(c *fiber.Ctx) error {
            // TODO: Add middleware logic here
            return c.Next()
        }
    }
    ```
