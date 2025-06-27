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
