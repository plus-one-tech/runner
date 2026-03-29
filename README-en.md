# runner

[日本語版はこちら](./README.md)

runner is a lightweight command runner for scripts and tasks.

It provides a simple way to execute source files and reusable tasks with a single command.

## Example

The following shows how `runner` executes commands:

```text
> runner hello.py
[runner] python hello.py
Hello Runner

> runner build
[runner] bash build.sh
Building...
Done.

> runner
[runner] bash runfile.sh
Running default task...
>
```

`runner` prints the actual command being executed, then runs it.

## Overview

`runner` is designed to make script execution and task running simple and predictable.

## Background

In everyday development, the same tasks are repeated over and over:
write code, build, test, and run.

`runner` was created to make these tasks simple to execute with a single command.

With modern tooling (for example, running C# directly from the command line),
it becomes possible to treat source files more like scripts.

For example:

```text
>run hello.cs
hello world
```

The goal is to make execution feel lightweight and immediate.

## Philosophy

`runner` is built on a few simple principles:

* **Unified execution** – run scripts and tasks using the same command
* **Transparency** – always show the actual command being executed
* **Minimal design** – avoid complex DSLs or heavy configuration

Unlike many task runners, `runner` does not introduce a new scripting language.
Task files simply contain normal program code.

## Basic Usage

### Run a script

```bash
runner hello.py
```

### Run a task

```bash
runner build
```

If `build.run` exists:

```text
#bash
dotnet run ./src/hello.cs
```

`runner` simply executes what is written in the `.run` file.

### Run default task

```bash
runner
```

If `runfile.run` exists, it will be executed.

## Common Options

### Preview execution

```bash
runner --dry-run build.run
runner --dry-run=windows install.run
runner --dry-run=all install.run
```

### Validate without execution

```bash
runner --check build.run
```

### Use a specific config file

```bash
runner --env ./runner.env install.run
```

### List available `.run` files

```bash
runner --list
```

## `.run` Files

A `.run` file is a simple executable task file.

Example:

```text
#python
print("Hello Runner")
```

Supported header styles:

```text
#python
#program.py
#.py
#script
```

The rest of the file is executed as normal program code.

## Configuration

`runner` uses `runner.env` to map runtimes and extensions to actual commands.

Example:

```text
runtime.python=python
runtime.bash=bash
runtime.node=node

ext.py=python
ext.js=node
ext.sh=bash
```

By default, `runner.env` is loaded from the user configuration directory.
You can also specify it explicitly with `--env`.

## Install Task

If `install.run` and `runner.env` are provided together:

```powershell
.\bin\runner --env ./runner.env install.run
```

### Notes (Windows)

On Windows, an executing `.exe` file cannot overwrite itself.

To install or update `runner`, always execute it from a different location than the target.

Use a separate binary (for example, `.\bin\runner`) to perform installation.

## Specification

Full specification is available here:

**docs/runner-spec.md**

## Status

This is an early release (v0.1.0).

The core functionality is implemented and usable, but the project is still evolving.

## Contributing

Pull requests are welcome, but changes may be declined if they do not align with the project's design philosophy.

## License

MIT License
