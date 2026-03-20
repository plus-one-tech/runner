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

It eliminates the need for multiple tools like make, shell scripts, or task runners by providing a single, unified entry point.

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

`runner` simply executes what is written in the `.run` file:

```bash
dotnet run ./src/build.cs
```

### Run default task

```bash
runner
```

If `runfile.run` exists, it will be executed.

## `.run` Files

A `.run` file is a simple executable task file.

Example:

```text
#python
print("Hello Runner")
```

The first line specifies the runtime.

Supported header styles:

```text
#python
#program.py
#.py
```

The rest of the file is executed as normal program code.

## Configuration

`runner` can be configured using `runner.env`.

Example:

```text
runtime.python=python
runtime.bash=bash
runtime.node=node

ext.py=python
ext.js=node
ext.sh=bash
```

This maps file extensions and runtimes to actual commands.

## Specification

Full specification is available here:

**docs/runner-spec.md**

## Status

The project is currently in the **specification stage**.

Implementation will start after the specification stabilizes.

## License

TBD
