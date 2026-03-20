# runner

runner is a lightweight command runner for scripts and tasks.

It provides a simple way to execute source files and reusable tasks with a single command.

```bash
runner hello.py
runner build
runner
```

## Overview

`runner` is designed to make script execution and task running simple and predictable.

It focuses on three goals:

* **Unified execution** – run scripts and tasks using the same command
* **Transparency** – always show the actual command being executed
* **Minimal design** – avoid complex DSLs or heavy configuration

Unlike many task runners, `runner` does not introduce a new scripting language. Task files simply contain normal program code.

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
dotnet run ./src/build.cs
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

### List available tasks

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

If `install.run` and `runner.env` are provided together, run:

```bash
runner --env ./runner.env install.run
```

## Specification

Full specification is available here:

**docs/runner-spec.md**

## Status

The project is currently in the final specification phase.

Implementation is planned next.

## License

TBD
